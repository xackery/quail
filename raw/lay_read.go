package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

// Lay is a raw layer struct
type Lay struct {
	MetaFileName string      `yaml:"file_name"`
	Version      uint32      `yaml:"version"`
	Entries      []*LayEntry `yaml:"entries"`
	names        []*nameEntry
	nameBuf      []byte
}

// Identity notes this is a lay file
func (lay *Lay) Identity() string {
	return "lay"
}

// LayEntry is a raw layer entry struct
type LayEntry struct {
	Material string `yaml:"material"`
	Diffuse  string `yaml:"diffuse"`
	Normal   string `yaml:"normal"`
}

// IsRaw notes this is a raw file
func (e *Lay) IsRaw() bool {
	return true
}

// Read takes data
func (lay *Lay) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGL" {
		return fmt.Errorf("invalid header %s, wanted EQGL", header)
	}

	lay.Version = dec.Uint32()
	versionOffset := 0
	switch lay.Version {
	case 2:
		versionOffset = 52 //32
	case 3:
		versionOffset = 16 //14
	case 4:
		versionOffset = 20
	default:
		return fmt.Errorf("unknown lay version: %d", lay.Version)
	}

	nameLength := int(dec.Uint32())
	layerCount := dec.Uint32()
	nameData := dec.Bytes(int(nameLength))

	names := make(map[int32]string)
	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[int32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	lay.NameSet(names)

	for i := 0; i < int(layerCount); i++ {
		entryID := dec.Uint32()
		layEntry := &LayEntry{}

		if entryID != 0xffffffff {
			layEntry.Material = lay.Name(int32(entryID))
		}

		entryID = dec.Uint32()
		if entryID != 0xffffffff {
			layEntry.Diffuse = lay.Name(int32(entryID))
		}

		entryID = dec.Uint32()
		if entryID != 0xffffffff {
			layEntry.Normal = lay.Name(int32(entryID))
		}
		dec.Bytes(versionOffset)
		//fmt.Println(hex.Dump())
		lay.Entries = append(lay.Entries, layEntry)
	}

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}
	return nil
}

// SetFileName sets the name of the file
func (lay *Lay) SetFileName(name string) {
	lay.MetaFileName = name
}

// FileName returns the name of the file
func (lay *Lay) FileName() string {
	return lay.MetaFileName
}

// Name is used during reading, returns the Name of an id
func (lay *Lay) Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if lay.names == nil {
		return fmt.Sprintf("!UNK(%d)", id)
	}
	//fmt.Println("name: [", names[id], "]")

	for _, v := range lay.names {
		if int32(v.offset) == id {
			return v.name
		}
	}
	return fmt.Sprintf("!UNK(%d)", id)
}

// NameSet is used during reading, sets the names within a buffer
func (lay *Lay) NameSet(newNames map[int32]string) {
	if newNames == nil {
		lay.names = []*nameEntry{}
		return
	}
	for k, v := range newNames {
		lay.names = append(lay.names, &nameEntry{offset: int(k), name: v})
	}
	lay.nameBuf = []byte{0x00}

	for _, v := range lay.names {
		lay.nameBuf = append(lay.nameBuf, []byte(v.name)...)
		lay.nameBuf = append(lay.nameBuf, 0)
	}
}

// NameAdd is used when writing, appending new names
func (lay *Lay) NameAdd(name string) int32 {

	if lay.names == nil {
		lay.names = []*nameEntry{
			{offset: 0, name: ""},
		}
		lay.nameBuf = []byte{0x00}
	}
	if name == "" {
		return 0
	}

	/* if name[len(lay.name)-1:] != "\x00" {
		name += "\x00"
	}
	*/
	if id := lay.NameOffset(name); id != -1 {
		return -id
	}
	lay.names = append(lay.names, &nameEntry{offset: len(lay.nameBuf), name: name})
	lastRef := int32(len(lay.nameBuf))
	lay.nameBuf = append(lay.nameBuf, []byte(name)...)
	lay.nameBuf = append(lay.nameBuf, 0)
	return int32(-lastRef)
}

func (lay *Lay) NameOffset(name string) int32 {
	if lay.names == nil {
		return -1
	}
	for _, v := range lay.names {
		if v.name == name {
			return int32(v.offset)
		}
	}
	return -1
}

// NameIndex is used when reading, returns the index of a name, or -1 if not found
func (lay *Lay) NameIndex(name string) int32 {
	if lay.names == nil {
		return -1
	}
	for k, v := range lay.names {
		if v.name == name {
			return int32(k)
		}
	}
	return -1
}

// NameData is used during writing, dumps the name cache
func (lay *Lay) NameData() []byte {

	return helper.WriteStringHash(string(lay.nameBuf))
}

// NameClear purges names and namebuf, called when encode starts
func (lay *Lay) NameClear() {
	lay.names = nil
	lay.nameBuf = nil
}
