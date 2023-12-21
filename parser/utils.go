package parser

import (
	"encoding/binary"
	"strings"
	"unicode/utf16"
)

func UTF16ToUTF8(in string) string {
	buff := strings.NewReader(in)
	u16 := make([]uint16, len(in)/2)
	binary.Read(buff, binary.LittleEndian, &u16)
	return string(utf16.Decode(u16))
}
