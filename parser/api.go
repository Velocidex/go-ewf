package parser

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/Velocidex/ordereddict"
	"www.velocidex.com/golang/go-ntfs/parser"
)

type chunk struct {
	reader     io.ReaderAt
	compressed bool
	offset     uint64
	size       int64
}

type EWFFile struct {
	mu sync.Mutex

	ChunkSize      int64
	NumberOfChunks int64
	TotalImageSize int64

	Metadata *ordereddict.Dict

	Tables []chunk

	lru  *parser.LRU
	Hits int64
	Miss int64

	Files []*EWF_file_header_v1

	// All the descriptors in the correct order.
	Descriptors []*EWF_section_descriptor_v1
}

func (self *EWFFile) reachChunk(chunk_id int) []byte {
	cached_page_buf, pres := self.lru.Get(chunk_id)
	if !pres {
		self.Miss += 1
		DebugPrint("Cache miss for %x (%x) (%d)\n",
			chunk_id, self.ChunkSize, self.lru.Len())

		// Read this page into memory.
		page_buf := make([]byte, self.ChunkSize)
		if chunk_id > len(self.Tables)-1 || chunk_id < 0 {
			return page_buf
		}

		// Chunk is not compressed, just return it as is
		chunk := self.Tables[chunk_id]
		if !chunk.compressed {
			n, _ := chunk.reader.ReadAt(page_buf, int64(chunk.offset))
			chunk_size := chunk.size
			if chunk_size > int64(n) {
				chunk_size = int64(n)
			}

			return page_buf[:chunk_size]
		}

		// Chunk is compressed
		compressed_chunk := make([]byte, chunk.size)
		n, err := chunk.reader.ReadAt(compressed_chunk, int64(chunk.offset))
		if err != nil {
			return page_buf
		}

		// Now attempt to decompress the chunk.
		c_reader := bytes.NewReader(compressed_chunk)
		r, err := zlib.NewReader(c_reader)
		if err != nil {
			return page_buf
		}

		// Sometimes it seems this returns a short read so we need to
		// resume it until we have an EOF or a full chunk.
		for offset := 0; n > 0 && err == nil; offset += n {
			n, err = r.Read(page_buf[offset:])
			DebugPrint("Decompressing error %v could only read %v bytes from %v buffer\n",
				err, n, len(page_buf))
		}

		self.lru.Add(chunk_id, page_buf)
		DebugPrint("Decompressed chunk %v into %v bytes\n", chunk_id, n)
		return page_buf
	}

	self.Hits += 1
	return cached_page_buf.([]byte)
}

func (self *EWFFile) ReadAt(buf []byte, offset int64) (int, error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	// Write on the output buffer
	buf_idx := 0
	for {
		// How much is left in this page to read?
		to_read := int(self.ChunkSize - offset%self.ChunkSize)

		// How much do we need to read into the buffer?
		if to_read > len(buf)-buf_idx {
			to_read = len(buf) - buf_idx
		}

		// Do not exceed the image size
		if to_read > int(self.TotalImageSize-offset) {
			to_read = int(self.TotalImageSize - offset)
		}

		// Are we done?
		if to_read <= 0 {
			return buf_idx, nil
		}

		// Fetch the chunk from the LRU
		chunkd_id := offset / self.ChunkSize
		page_buf := self.reachChunk(int(chunkd_id))

		// Copy the relevant data from the page.
		page_offset := int(offset % self.ChunkSize)
		copy(buf[buf_idx:buf_idx+to_read],
			page_buf[page_offset:page_offset+to_read])

		offset += int64(to_read)
		buf_idx += to_read
		if DEBUG != nil && (self.Hits+self.Miss)%10000 == 0 {
			DebugPrint("PageCache hit %v miss %v (%v)\n", self.Hits, self.Miss,
				float64(self.Hits)/float64(self.Miss))
		}
	}
}

// These are mostly useless metadata fields we dont care about right
// now. It is not an error if we can not handle any of it.
func parseHeader(ewf *EWFFile, descriptor *EWF_section_descriptor_v1) error {
	if len(ewf.Metadata.Keys()) > 0 {
		return nil
	}

	section_size := descriptor.SectionSize()
	if section_size > 1*1024*1024 {
		return nil
	}

	c_buf := make([]byte, section_size)
	n, err := descriptor.Reader.ReadAt(c_buf, descriptor.Offset+
		int64(descriptor.Size()))
	if err != nil {
		return nil
	}

	// Try to decompress this data
	c_reader := bytes.NewReader(c_buf[:n])
	r, err := zlib.NewReader(c_reader)
	if err != nil {
		return nil
	}

	buff := &bytes.Buffer{}
	_, err = io.CopyN(buff, r, 1*1024*1024)
	if err != io.EOF && err != nil {
		return nil
	}

	data_bytes := buff.Bytes()
	if len(data_bytes) < 10 {
		return nil
	}

	// Starts with a BOM
	data := string(data_bytes)
	if data_bytes[0] == 255 || data_bytes[1] == 254 {
		data = UTF16ToUTF8(data)
	}

	lines := strings.Split(data, "\n")
	if len(lines) < 5 || lines[1] != "main" {
		return nil
	}

	keys := strings.Split(lines[2], "\t")
	values := strings.Split(lines[3], "\t")

	if len(keys) == len(values) {
		for i, k := range keys {
			v := values[i]

			switch k {
			case "a":
				k = "Unique description"
			case "c":
				k = "Case Number"
			case "n":
				k = "Evidence number"
			case "e":
				k = "Examiner name"
			case "t":
				k = "Notes"
			case "md":
				k = "Model"
			case "sn":
				k = "Serial Number"
			case "l":
				k = "Device label"
			case "av":
				k = "Version"
			case "ov":
				k = "Platform"
			case "m":
				k = "Acquisition date and time"
			case "u":
				k = "Systemdate and time"
			case "p":
				k = "Password hash"
			case "pid":
				k = "Process identifier"
			case "ext":
				k = "Extents"
			case "r":
				k = "Compression"
			}
			ewf.Metadata.Set(k, v)
		}
	}
	return nil
}

func parseVolume(ewf *EWFFile, descriptor *EWF_section_descriptor_v1) error {
	volume := descriptor.Profile.EWF_volume(descriptor.Reader,
		descriptor.Offset+int64(descriptor.Size()))
	ewf.ChunkSize = int64(volume.Sectors_per_chunk() * volume.Bytes_per_sector())
	ewf.NumberOfChunks = int64(volume.Number_of_chunks())
	ewf.TotalImageSize = int64(volume.Bytes_per_sector()) *
		int64(volume.Number_of_sectors())

	// Prepare some space for chunks
	ewf.Tables = make([]chunk, 0, ewf.NumberOfChunks)

	return nil
}

func parseTable(ewf *EWFFile, descriptor *EWF_section_descriptor_v1) error {
	table := descriptor.Profile.EWF_table_header_v1(descriptor.Reader,
		descriptor.Offset+int64(descriptor.Size()))

	base_offset := table.Base_offset()
	start := table.Entries().Offset

	previous_chunk_offset := uint64(0)

	// Read all the table offsets into one buffer for speed, otherwise
	// we will be making lots of very small reads into the underlying
	// reader one per table entry.
	count := int(table.Number_of_entries())
	buff := make([]byte, (count+1)*4)
	_, err := table.Reader.ReadAt(buff, start)
	if err != nil {
		return err
	}

	mem_reader := bytes.NewReader(buff)

	for i := 0; i < count; i++ {
		e := table.Profile.EWF_table_entry(mem_reader, 4*int64(i))

		current_offset := e.ChunkOffset() + base_offset
		// Update the size of the last chunk based on the current
		// offset.
		if i > 1 {
			previous_chunk := &ewf.Tables[len(ewf.Tables)-1]
			previous_chunk_size := int64(current_offset - previous_chunk_offset)
			previous_chunk.size = previous_chunk_size

			// This can not happen!
			if previous_chunk_size < 0 {
				return fmt.Errorf("Negative chunk size for chunk %v",
					len(ewf.Tables))
			}
		}

		ewf.Tables = append(ewf.Tables, chunk{
			reader:     table.Reader,
			compressed: e.Compressed() > 0,
			offset:     current_offset,
			size:       ewf.ChunkSize,
		})

		previous_chunk_offset = current_offset
	}

	return nil
}

func parseDescriptor(ewf *EWFFile, descriptor *EWF_section_descriptor_v1) error {

	ewf.Descriptors = append(ewf.Descriptors, descriptor)

	descriptor_type := strings.SplitN(descriptor.Type(), "\x00", 2)[0]
	switch descriptor_type {
	case "header", "header2":
		return parseHeader(ewf, descriptor)

	case "volume", "disk":
		return parseVolume(ewf, descriptor)

	case "table":
		return parseTable(ewf, descriptor)

	}

	return nil
}

func OpenEWFFile(options *EWFOptions, readers ...io.ReaderAt) (
	*EWFFile, error) {
	profile := NewEWFProfile()

	cache_size := 100
	if options != nil && options.LRUSize > 0 {
		cache_size = options.LRUSize
	}

	// By default 10mb cache.
	cache, err := parser.NewLRU(cache_size, nil, "EWFPagesReader")
	if err != nil {
		return nil, err
	}

	ewf := &EWFFile{
		lru:      cache,
		Metadata: ordereddict.NewDict(),
	}

	// Read all the files and sort them in order
	for _, r := range readers {
		header := profile.EWF_file_header_v1(r, 0)
		ewf.Files = append(ewf.Files, header)
	}

	sort.Slice(ewf.Files, func(i, j int) bool {
		return ewf.Files[i].Segment_number() < ewf.Files[j].Segment_number()
	})

	current_segment := 1
	for _, header := range ewf.Files {
		segment_number := int(header.Segment_number())

		// Ignore repeated segments
		if segment_number < current_segment {
			continue
		}

		// Is the segment unexpected?
		if current_segment != segment_number {
			return nil, fmt.Errorf("Missing segment %v!", current_segment)
		}
		current_segment++

		descriptor := header.Descriptor()
		for descriptor.SectionSize() > 0 {
			err := parseDescriptor(ewf, descriptor)
			if err != nil {
				return nil, err
			}

			// Descriptors must be in order
			next_descriptor := descriptor.Next()
			if next_descriptor.Offset == 0 ||
				next_descriptor.Offset <= descriptor.Offset {
				break
			}
			descriptor = next_descriptor
		}
	}

	// We can not proceed without a valid chunk size.
	if ewf.ChunkSize == 0 {
		return nil, fmt.Errorf("Unable to parse chunk size from image.")
	}

	if ewf.TotalImageSize == 0 {
		ewf.TotalImageSize = ewf.ChunkSize * ewf.NumberOfChunks
	}

	return ewf, nil
}
