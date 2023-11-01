package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/tag"
)

// Lay is a raw layer struct
type Lay struct {
	MetaFileName string      `yaml:"file_name"`
	Version      uint32      `yaml:"version"`
	Entries      []*LayEntry `yaml:"entries"`
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

	tag.New()

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
	tag.Add(0, dec.Pos(), "red", "header")
	nameData := dec.Bytes(int(nameLength))
	tag.Add(tag.LastPos(), dec.Pos(), "green", "names")

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

	NamesSet(names)

	for i := 0; i < int(layerCount); i++ {
		entryID := dec.Uint32()
		layEntry := &LayEntry{}

		if entryID != 0xffffffff {
			layEntry.Material = Name(int32(entryID))
		}

		entryID = dec.Uint32()
		if entryID != 0xffffffff {
			layEntry.Diffuse = Name(int32(entryID))
		}

		entryID = dec.Uint32()
		if entryID != 0xffffffff {
			layEntry.Normal = Name(int32(entryID))
		}
		dec.Bytes(versionOffset)
		//fmt.Println(hex.Dump())
		lay.Entries = append(lay.Entries, layEntry)
		tag.AddRand(tag.LastPos(), dec.Pos(), fmt.Sprintf("%d|%s|%s|%s", i, layEntry.Material, layEntry.Diffuse, layEntry.Normal))
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
