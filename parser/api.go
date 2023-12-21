package parser

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

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

		n, _ = r.Read(page_buf)
		page_buf = page_buf[:n]

		self.lru.Add(chunk_id, page_buf)
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

		// Are we done?
		if to_read == 0 {
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
			fmt.Printf("PageCache hit %v miss %v (%v)\n", self.Hits, self.Miss,
				float64(self.Hits)/float64(self.Miss))
		}
	}
}

func parseVolume(ewf *EWFFile, descriptor *EWF_section_descriptor_v1) error {
	volume := descriptor.Profile.EWF_volume(descriptor.Reader,
		descriptor.Offset+int64(descriptor.Size()))
	ewf.ChunkSize = int64(volume.Sectors_per_chunk() * volume.Bytes_per_sector())
	ewf.NumberOfChunks = int64(volume.Number_of_chunks())

	// Prepare some space for chunks
	ewf.Tables = make([]chunk, 0, ewf.NumberOfChunks)

	return nil
}

func parseTable(ewf *EWFFile, descriptor *EWF_section_descriptor_v1) error {
	table := descriptor.Profile.EWF_table_header_v1(descriptor.Reader,
		descriptor.Offset+int64(descriptor.Size()))

	base_offset := table.Base_offset()
	start := table.Entries().Offset

	for i := 0; i < int(table.Number_of_entries()); i++ {
		e := table.Profile.EWF_table_entry(table.Reader, start+4*int64(i))

		current_offset := e.ChunkOffset() + base_offset
		ewf.Tables = append(ewf.Tables, chunk{
			reader:     e.Reader,
			compressed: e.Compressed() > 0,
			offset:     current_offset,
			size:       ewf.ChunkSize,
		})

		// Update the size of the last chunk based on the current
		// offset.
		if i > 0 {
			ewf.Tables[i-1].size = int64(current_offset - ewf.Tables[i-1].offset)
		}
	}

	return nil
}

func parseDescriptor(ewf *EWFFile, descriptor *EWF_section_descriptor_v1) error {

	ewf.Descriptors = append(ewf.Descriptors, descriptor)

	descriptor_type := strings.SplitN(descriptor.Type(), "\x00", 2)[0]
	switch descriptor_type {
	case "header", "header2":
		// These are useless metadata fields we dont care about right now.
		return nil

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
		lru: cache,
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
			descriptor = descriptor.Next()
		}
	}

	return ewf, nil
}
