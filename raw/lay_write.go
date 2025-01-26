package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Write will write a lay file
func (lay *Lay) Write(w io.Writer) error {
	if lay.name == nil {
		lay.name = &eqgName{}
	}
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.String("EQGL")

	enc.Uint32(uint32(lay.Version))
	var versionPadding []byte
	switch lay.Version {
	case 2:
		versionPadding = make([]byte, 52) //32
	case 3:
		versionPadding = make([]byte, 16) //14
	case 4:
		versionPadding = make([]byte, 20)
	default:
		return fmt.Errorf("unknown lay version: %d", lay.Version)
	}

	lay.name.clear()
	for _, layEntry := range lay.Layers {
		lay.name.add(layEntry.Material)
		lay.name.add(layEntry.Diffuse)
		lay.name.add(layEntry.Normal)
	}

	enc.Uint32(uint32(len(lay.name.data()))) // nameLength
	enc.Uint32(uint32(len(lay.Layers)))      //layerCount
	enc.Bytes(lay.name.data())               // nameData

	for _, layEntry := range lay.Layers {
		offset := uint32(0xffffffff)
		if layEntry.Material != "" {
			offset = uint32(lay.name.offsetByName(layEntry.Material))
		}
		enc.Uint32(offset)

		offset = uint32(0xffffffff)
		if layEntry.Diffuse != "" {
			offset = uint32(lay.name.offsetByName(layEntry.Diffuse))
		}
		enc.Uint32(offset)

		offset = uint32(0xffffffff)
		if layEntry.Normal != "" {
			offset = uint32(lay.name.offsetByName(layEntry.Normal))
		}
		enc.Uint32(offset)

		enc.Bytes(versionPadding)
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil

}
