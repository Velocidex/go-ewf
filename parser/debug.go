package parser

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/repr"
)

var (
	DEBUG *bool
)

func SetDebug() {
	value := true
	DEBUG = &value
}

func Debug(arg interface{}) {
	if arg != nil {
		repr.Println(arg)
	} else {
		repr.Println("nil")
	}
}

func DebugPrint(fmt_str string, v ...interface{}) {
	if DEBUG == nil {
		// os.Environ() seems very expensive in Go so we cache
		// it.
		for _, x := range os.Environ() {
			if strings.HasPrefix(x, "EWF_DEBUG=1") {
				value := true
				DEBUG = &value
				break
			}
		}
	}

	if DEBUG == nil {
		value := false
		DEBUG = &value
	}

	if *DEBUG {
		fmt.Printf(fmt_str, v...)
	}
}

func getName(reader io.ReaderAt) string {
	f, ok := reader.(*os.File)
	if ok {
		return f.Name()
	}
	return ""
}

func (self *EWFFile) WriteDebug(w io.Writer) {
	w.Write([]byte(fmt.Sprintf(`
ChunkSize: %v
NumberOfChunks: %v
TotalImageSize: %v

`, self.ChunkSize, self.NumberOfChunks, self.TotalImageSize)))

	keys := self.Metadata.Keys()
	if len(keys) > 0 {
		w.Write([]byte("Metadata:\n"))
		for _, k := range keys {
			v, _ := self.Metadata.Get(k)
			w.Write([]byte(fmt.Sprintf("   %v: %v\n", k, v)))
		}
	}

	for i, descriptor := range self.Descriptors {
		descriptor_type := strings.SplitN(descriptor.Type(), "\x00", 2)[0]
		w.Write([]byte(
			fmt.Sprintf("  descriptor %v @ %#x: %v - %v\n", i, descriptor.Offset,
				descriptor_type, getName(descriptor.Reader))))
	}

	for i, chunk := range self.Tables {
		compressed := ""
		if chunk.compressed {
			compressed = " (Compr)"
		}

		w.Write([]byte(fmt.Sprintf("  chunk %v: %v (%v) from %v%v\n",
			i, chunk.offset, chunk.size, getName(chunk.reader), compressed)))
	}
}

func DlvBreak() {

}
