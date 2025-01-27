package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Lay is a raw layer struct
type Lay struct {
	MetaFileName string
	Version      uint32
	Layers       []*LayEntry
	name         *eqgName
}

// Identity notes this is a lay file
func (lay *Lay) Identity() string {
	return "lay"
}

func (lay *Lay) String() string {
	out := fmt.Sprintf("Lay: %s,", lay.MetaFileName)
	out += fmt.Sprintf(" %d names,", lay.name.len())
	out += fmt.Sprintf(" %d layers", len(lay.Layers))
	if len(lay.Layers) > 0 {
		out += " ["

		for i, layer := range lay.Layers {
			out += layer.Material
			if i < len(lay.Layers)-1 {
				out += ", "
			}
		}
		out += "]"

	}

	return out
}

// LayEntry is a raw layer entry struct
type LayEntry struct {
	Material string
	Diffuse  string
	Normal   string
}

// IsRaw notes this is a raw file
func (e *Lay) IsRaw() bool {
	return true
}

// Read takes data
func (lay *Lay) Read(r io.ReadSeeker) error {
	if lay.name == nil {
		lay.name = &eqgName{}
	}
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

	lay.name.parse(nameData)

	for i := 0; i < int(layerCount); i++ {
		entryID := dec.Uint32()
		layEntry := &LayEntry{}

		if entryID != 0xffffffff {
			layEntry.Material = lay.name.byOffset(int32(entryID))
		}

		entryID = dec.Uint32()
		if entryID != 0xffffffff {
			layEntry.Diffuse = lay.name.byOffset(int32(entryID))
		}

		entryID = dec.Uint32()
		if entryID != 0xffffffff {
			layEntry.Normal = lay.name.byOffset(int32(entryID))
		}
		dec.Bytes(versionOffset)
		//fmt.Println(hex.Dump())
		lay.Layers = append(lay.Layers, layEntry)
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
