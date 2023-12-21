
package parser

// Autogenerated code from ewf_profile.json. Do not edit.

import (
    "encoding/binary"
    "fmt"
    "bytes"
    "io"
    "sort"
    "strings"
    "unicode/utf16"
    "unicode/utf8"
)

var (
   // Depending on autogenerated code we may use this. Add a reference
   // to shut the compiler up.
   _ = bytes.MinRead
   _ = fmt.Sprintf
   _ = utf16.Decode
   _ = binary.LittleEndian
   _ = utf8.RuneError
   _ = sort.Strings
   _ = strings.Join
   _ = io.Copy
)

func indent(text string) string {
    result := []string{}
    lines := strings.Split(text,"\n")
    for _, line := range lines {
         result = append(result, "  " + line)
    }
    return strings.Join(result, "\n")
}


type EWFProfile struct {
    Off_EWF_file_header_v1_EVF_sig int64
    Off_EWF_file_header_v1_Fields_start int64
    Off_EWF_file_header_v1_Segment_number int64
    Off_EWF_file_header_v1_Fields_end int64
    Off_EWF_file_header_v1_Descriptor int64
    Off_EWF_file_header_v2_EVF_sig int64
    Off_EWF_file_header_v2_Major_version int64
    Off_EWF_file_header_v2_Minor_version int64
    Off_EWF_file_header_v2_Compression_method int64
    Off_EWF_file_header_v2_Segment_number int64
    Off_EWF_file_header_v2_Set_identifier int64
    Off_EWF_section_descriptor_v1_Type int64
    Off_EWF_section_descriptor_v1_Next int64
    Off_EWF_section_descriptor_v1_SectionSize int64
    Off_EWF_section_descriptor_v1_Checksum int64
    Off_EWF_table_entry_Compressed int64
    Off_EWF_table_entry_ChunkOffset int64
    Off_EWF_table_header_v1_Number_of_entries int64
    Off_EWF_table_header_v1_Base_offset int64
    Off_EWF_table_header_v1_Checksum int64
    Off_EWF_table_header_v1_Entries int64
    Off_EWF_volume_Media_type int64
    Off_EWF_volume_Number_of_chunks int64
    Off_EWF_volume_Sectors_per_chunk int64
    Off_EWF_volume_Bytes_per_sector int64
    Off_EWF_volume_Number_of_sectors int64
    Off_EWF_volume_Chs_cylinders int64
    Off_EWF_volume_Chs_heads int64
    Off_EWF_volume_Chs_sectors int64
    Off_EWF_volume_Media_flags int64
    Off_EWF_volume_Compression_level int64
    Off_EWF_volume_Checksum int64
}

func NewEWFProfile() *EWFProfile {
    // Specific offsets can be tweaked to cater for slight version mismatches.
    self := &EWFProfile{0,8,9,11,13,0,9,10,11,13,15,0,16,24,72,0,0,0,8,20,24,0,4,8,12,16,24,28,32,36,52,90}
    return self
}

func (self *EWFProfile) EWF_file_header_v1(reader io.ReaderAt, offset int64) *EWF_file_header_v1 {
    return &EWF_file_header_v1{Reader: reader, Offset: offset, Profile: self}
}

func (self *EWFProfile) EWF_file_header_v2(reader io.ReaderAt, offset int64) *EWF_file_header_v2 {
    return &EWF_file_header_v2{Reader: reader, Offset: offset, Profile: self}
}

func (self *EWFProfile) EWF_section_descriptor_v1(reader io.ReaderAt, offset int64) *EWF_section_descriptor_v1 {
    return &EWF_section_descriptor_v1{Reader: reader, Offset: offset, Profile: self}
}

func (self *EWFProfile) EWF_table_entry(reader io.ReaderAt, offset int64) *EWF_table_entry {
    return &EWF_table_entry{Reader: reader, Offset: offset, Profile: self}
}

func (self *EWFProfile) EWF_table_header_v1(reader io.ReaderAt, offset int64) *EWF_table_header_v1 {
    return &EWF_table_header_v1{Reader: reader, Offset: offset, Profile: self}
}

func (self *EWFProfile) EWF_volume(reader io.ReaderAt, offset int64) *EWF_volume {
    return &EWF_volume{Reader: reader, Offset: offset, Profile: self}
}


type EWF_file_header_v1 struct {
    Reader io.ReaderAt
    Offset int64
    Profile *EWFProfile
}

func (self *EWF_file_header_v1) Size() int {
    return 13
}


func (self *EWF_file_header_v1) EVF_sig() string {
  return ParseString(self.Reader, self.Profile.Off_EWF_file_header_v1_EVF_sig + self.Offset, 8)
}

func (self *EWF_file_header_v1) Fields_start() int8 {
   return ParseInt8(self.Reader, self.Profile.Off_EWF_file_header_v1_Fields_start + self.Offset)
}

func (self *EWF_file_header_v1) Segment_number() uint16 {
   return ParseUint16(self.Reader, self.Profile.Off_EWF_file_header_v1_Segment_number + self.Offset)
}

func (self *EWF_file_header_v1) Fields_end() uint16 {
   return ParseUint16(self.Reader, self.Profile.Off_EWF_file_header_v1_Fields_end + self.Offset)
}

func (self *EWF_file_header_v1) Descriptor() *EWF_section_descriptor_v1 {
    return self.Profile.EWF_section_descriptor_v1(self.Reader, self.Profile.Off_EWF_file_header_v1_Descriptor + self.Offset)
}
func (self *EWF_file_header_v1) DebugString() string {
    result := fmt.Sprintf("struct EWF_file_header_v1 @ %#x:\n", self.Offset)
    result += fmt.Sprintf("  EVF_sig: %v\n", string(self.EVF_sig()))
    result += fmt.Sprintf("  Fields_start: %#0x\n", self.Fields_start())
    result += fmt.Sprintf("  Segment_number: %#0x\n", self.Segment_number())
    result += fmt.Sprintf("  Fields_end: %#0x\n", self.Fields_end())
    result += fmt.Sprintf("  Descriptor: {\n%v}\n", indent(self.Descriptor().DebugString()))
    return result
}

type EWF_file_header_v2 struct {
    Reader io.ReaderAt
    Offset int64
    Profile *EWFProfile
}

func (self *EWF_file_header_v2) Size() int {
    return 0
}


func (self *EWF_file_header_v2) EVF_sig() string {
  return ParseString(self.Reader, self.Profile.Off_EWF_file_header_v2_EVF_sig + self.Offset, 8)
}

func (self *EWF_file_header_v2) Major_version() int8 {
   return ParseInt8(self.Reader, self.Profile.Off_EWF_file_header_v2_Major_version + self.Offset)
}

func (self *EWF_file_header_v2) Minor_version() int8 {
   return ParseInt8(self.Reader, self.Profile.Off_EWF_file_header_v2_Minor_version + self.Offset)
}

func (self *EWF_file_header_v2) Compression_method() *Enumeration {
   value := ParseUint16(self.Reader, self.Profile.Off_EWF_file_header_v2_Compression_method + self.Offset)
   name := "Unknown"
   switch value {

      case 0:
         name = "NONE"

      case 1:
         name = "DEFLATE"

      case 2:
         name = "BZIP2"
}
   return &Enumeration{Value: uint64(value), Name: name}
}


func (self *EWF_file_header_v2) Segment_number() uint16 {
   return ParseUint16(self.Reader, self.Profile.Off_EWF_file_header_v2_Segment_number + self.Offset)
}


func (self *EWF_file_header_v2) Set_identifier() string {
  return ParseString(self.Reader, self.Profile.Off_EWF_file_header_v2_Set_identifier + self.Offset, 16)
}
func (self *EWF_file_header_v2) DebugString() string {
    result := fmt.Sprintf("struct EWF_file_header_v2 @ %#x:\n", self.Offset)
    result += fmt.Sprintf("  EVF_sig: %v\n", string(self.EVF_sig()))
    result += fmt.Sprintf("  Major_version: %#0x\n", self.Major_version())
    result += fmt.Sprintf("  Minor_version: %#0x\n", self.Minor_version())
    result += fmt.Sprintf("  Compression_method: %v\n", self.Compression_method().DebugString())
    result += fmt.Sprintf("  Segment_number: %#0x\n", self.Segment_number())
    result += fmt.Sprintf("  Set_identifier: %v\n", string(self.Set_identifier()))
    return result
}

type EWF_section_descriptor_v1 struct {
    Reader io.ReaderAt
    Offset int64
    Profile *EWFProfile
}

func (self *EWF_section_descriptor_v1) Size() int {
    return 76
}


func (self *EWF_section_descriptor_v1) Type() string {
  return ParseString(self.Reader, self.Profile.Off_EWF_section_descriptor_v1_Type + self.Offset, 16)
}

func (self *EWF_section_descriptor_v1) Next() *EWF_section_descriptor_v1 {
   deref := ParseUint64(self.Reader, self.Profile.Off_EWF_section_descriptor_v1_Next + self.Offset)
   return self.Profile.EWF_section_descriptor_v1(self.Reader, int64(deref))
}

func (self *EWF_section_descriptor_v1) SectionSize() uint64 {
    return ParseUint64(self.Reader, self.Profile.Off_EWF_section_descriptor_v1_SectionSize + self.Offset)
}

func (self *EWF_section_descriptor_v1) Checksum() uint32 {
   return ParseUint32(self.Reader, self.Profile.Off_EWF_section_descriptor_v1_Checksum + self.Offset)
}
func (self *EWF_section_descriptor_v1) DebugString() string {
    result := fmt.Sprintf("struct EWF_section_descriptor_v1 @ %#x:\n", self.Offset)
    result += fmt.Sprintf("  Type: %v\n", string(self.Type()))
    result += fmt.Sprintf("  SectionSize: %#0x\n", self.SectionSize())
    result += fmt.Sprintf("  Checksum: %#0x\n", self.Checksum())
    return result
}

type EWF_table_entry struct {
    Reader io.ReaderAt
    Offset int64
    Profile *EWFProfile
}

func (self *EWF_table_entry) Size() int {
    return 4
}

func (self *EWF_table_entry) Compressed() uint64 {
   value := ParseUint64(self.Reader, self.Profile.Off_EWF_table_entry_Compressed + self.Offset)
   return (uint64(value) & 0xffffffff) >> 0x1f
}

func (self *EWF_table_entry) ChunkOffset() uint64 {
   value := ParseUint64(self.Reader, self.Profile.Off_EWF_table_entry_ChunkOffset + self.Offset)
   return (uint64(value) & 0x7fffffff) >> 0x0
}
func (self *EWF_table_entry) DebugString() string {
    result := fmt.Sprintf("struct EWF_table_entry @ %#x:\n", self.Offset)
    result += fmt.Sprintf("  Compressed: %#0x\n", self.Compressed())
    result += fmt.Sprintf("  ChunkOffset: %#0x\n", self.ChunkOffset())
    return result
}

type EWF_table_header_v1 struct {
    Reader io.ReaderAt
    Offset int64
    Profile *EWFProfile
}

func (self *EWF_table_header_v1) Size() int {
    return 0
}

func (self *EWF_table_header_v1) Number_of_entries() uint64 {
    return ParseUint64(self.Reader, self.Profile.Off_EWF_table_header_v1_Number_of_entries + self.Offset)
}

func (self *EWF_table_header_v1) Base_offset() uint64 {
    return ParseUint64(self.Reader, self.Profile.Off_EWF_table_header_v1_Base_offset + self.Offset)
}

func (self *EWF_table_header_v1) Checksum() uint32 {
   return ParseUint32(self.Reader, self.Profile.Off_EWF_table_header_v1_Checksum + self.Offset)
}

func (self *EWF_table_header_v1) Entries() *EWF_table_entry {
    return self.Profile.EWF_table_entry(self.Reader, self.Profile.Off_EWF_table_header_v1_Entries + self.Offset)
}
func (self *EWF_table_header_v1) DebugString() string {
    result := fmt.Sprintf("struct EWF_table_header_v1 @ %#x:\n", self.Offset)
    result += fmt.Sprintf("  Number_of_entries: %#0x\n", self.Number_of_entries())
    result += fmt.Sprintf("  Base_offset: %#0x\n", self.Base_offset())
    result += fmt.Sprintf("  Checksum: %#0x\n", self.Checksum())
    result += fmt.Sprintf("  Entries: {\n%v}\n", indent(self.Entries().DebugString()))
    return result
}

type EWF_volume struct {
    Reader io.ReaderAt
    Offset int64
    Profile *EWFProfile
}

func (self *EWF_volume) Size() int {
    return 94
}

func (self *EWF_volume) Media_type() *Enumeration {
   value := ParseUint64(self.Reader, self.Profile.Off_EWF_volume_Media_type + self.Offset)
   name := "Unknown"
   switch value {

      case 0:
         name = "remobable_disk"

      case 1:
         name = "fixed_disk"

      case 2:
         name = "optical_disk"

      case 3:
         name = "LVF"

      case 4:
         name = "memory"
}
   return &Enumeration{Value: uint64(value), Name: name}
}


func (self *EWF_volume) Number_of_chunks() uint32 {
   return ParseUint32(self.Reader, self.Profile.Off_EWF_volume_Number_of_chunks + self.Offset)
}

func (self *EWF_volume) Sectors_per_chunk() uint32 {
   return ParseUint32(self.Reader, self.Profile.Off_EWF_volume_Sectors_per_chunk + self.Offset)
}

func (self *EWF_volume) Bytes_per_sector() uint32 {
   return ParseUint32(self.Reader, self.Profile.Off_EWF_volume_Bytes_per_sector + self.Offset)
}

func (self *EWF_volume) Number_of_sectors() uint64 {
    return ParseUint64(self.Reader, self.Profile.Off_EWF_volume_Number_of_sectors + self.Offset)
}

func (self *EWF_volume) Chs_cylinders() uint32 {
   return ParseUint32(self.Reader, self.Profile.Off_EWF_volume_Chs_cylinders + self.Offset)
}

func (self *EWF_volume) Chs_heads() uint32 {
   return ParseUint32(self.Reader, self.Profile.Off_EWF_volume_Chs_heads + self.Offset)
}

func (self *EWF_volume) Chs_sectors() uint32 {
   return ParseUint32(self.Reader, self.Profile.Off_EWF_volume_Chs_sectors + self.Offset)
}

func (self *EWF_volume) Media_flags() *Flags {
   value := ParseUint64(self.Reader, self.Profile.Off_EWF_volume_Media_flags + self.Offset)
   names := make(map[string]bool)


   if value & 1 != 0 {
      names["image"] = true
   }

   if value & 2 != 0 {
      names["physical"] = true
   }

   if value & 4 != 0 {
      names["Fastblock Tableau write blocker"] = true
   }

   if value & 8 != 0 {
      names["Tableau write blocker"] = true
   }

   return &Flags{Value: uint64(value), Names: names}
}


func (self *EWF_volume) Compression_level() *Enumeration {
   value := ParseUint64(self.Reader, self.Profile.Off_EWF_volume_Compression_level + self.Offset)
   name := "Unknown"
   switch value {

      case 0:
         name = "no compression"

      case 1:
         name = "fast/good compression"

      case 2:
         name = "best compression"
}
   return &Enumeration{Value: uint64(value), Name: name}
}


func (self *EWF_volume) Checksum() int32 {
   return ParseInt32(self.Reader, self.Profile.Off_EWF_volume_Checksum + self.Offset)
}
func (self *EWF_volume) DebugString() string {
    result := fmt.Sprintf("struct EWF_volume @ %#x:\n", self.Offset)
    result += fmt.Sprintf("  Media_type: %v\n", self.Media_type().DebugString())
    result += fmt.Sprintf("  Number_of_chunks: %#0x\n", self.Number_of_chunks())
    result += fmt.Sprintf("  Sectors_per_chunk: %#0x\n", self.Sectors_per_chunk())
    result += fmt.Sprintf("  Bytes_per_sector: %#0x\n", self.Bytes_per_sector())
    result += fmt.Sprintf("  Number_of_sectors: %#0x\n", self.Number_of_sectors())
    result += fmt.Sprintf("  Chs_cylinders: %#0x\n", self.Chs_cylinders())
    result += fmt.Sprintf("  Chs_heads: %#0x\n", self.Chs_heads())
    result += fmt.Sprintf("  Chs_sectors: %#0x\n", self.Chs_sectors())
    result += fmt.Sprintf("  Media_flags: %v\n", self.Media_flags().DebugString())
    result += fmt.Sprintf("  Compression_level: %v\n", self.Compression_level().DebugString())
    result += fmt.Sprintf("  Checksum: %#0x\n", self.Checksum())
    return result
}

type Enumeration struct {
    Value uint64
    Name  string
}

func (self Enumeration) DebugString() string {
    return fmt.Sprintf("%s (%d)", self.Name, self.Value)
}


type Flags struct {
    Value uint64
    Names  map[string]bool
}

func (self Flags) DebugString() string {
    names := []string{}
    for k, _ := range self.Names {
      names = append(names, k)
    }

    sort.Strings(names)

    return fmt.Sprintf("%d (%s)", self.Value, strings.Join(names, ","))
}

func (self Flags) IsSet(flag string) bool {
    result, _ := self.Names[flag]
    return result
}

func (self Flags) Values() []string {
    result := make([]string, 0, len(self.Names))
    for k, _ := range self.Names {
       result = append(result, k)
    }
    return result
}


func ParseInt32(reader io.ReaderAt, offset int64) int32 {
    data := make([]byte, 4)
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return int32(binary.LittleEndian.Uint32(data))
}

func ParseInt8(reader io.ReaderAt, offset int64) int8 {
    result := make([]byte, 1)
    _, err := reader.ReadAt(result, offset)
    if err != nil {
       return 0
    }
    return int8(result[0])
}

func ParseUint16(reader io.ReaderAt, offset int64) uint16 {
    data := make([]byte, 2)
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return binary.LittleEndian.Uint16(data)
}

func ParseUint32(reader io.ReaderAt, offset int64) uint32 {
    data := make([]byte, 4)
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return binary.LittleEndian.Uint32(data)
}

func ParseUint64(reader io.ReaderAt, offset int64) uint64 {
    data := make([]byte, 8)
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return binary.LittleEndian.Uint64(data)
}

func ParseTerminatedString(reader io.ReaderAt, offset int64) string {
   data := make([]byte, 1024)
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
     return ""
   }
   idx := bytes.Index(data[:n], []byte{0})
   if idx < 0 {
      idx = n
   }
   return string(data[0:idx])
}

func ParseString(reader io.ReaderAt, offset int64, length int64) string {
   data := make([]byte, length)
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
      return ""
   }
   return string(data[:n])
}


