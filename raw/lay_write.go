package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// Write will write a lay file
func (lay *Lay) Write(w io.Writer) error {
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

	tmpNames := []string{}
	for _, lay := range lay.Entries {
		if lay.Material != "" {
			tmpNames = append(tmpNames, lay.Material)
		}
		if lay.Diffuse != "" {
			tmpNames = append(tmpNames, lay.Diffuse)
		}
		if lay.Normal != "" {
			tmpNames = append(tmpNames, lay.Normal)
		}
	}
	names, nameData, err := common.NameBuild(tmpNames)
	if err != nil {
		return fmt.Errorf("nameBuild: %w", err)
	}

	enc.Uint32(uint32(len(nameData)))    // nameLength
	enc.Uint32(uint32(len(lay.Entries))) //layerCount
	enc.Bytes(nameData)                  // nameData

	for _, layEntry := range lay.Entries {
		offset := uint32(0xffffffff)
		if layEntry.Material != "" {
			offset = uint32(names[layEntry.Material])
		}
		enc.Uint32(offset)

		offset = uint32(0xffffffff)
		if layEntry.Diffuse != "" {
			offset = uint32(names[layEntry.Diffuse])
		}
		enc.Uint32(offset)

		offset = uint32(0xffffffff)
		if layEntry.Normal != "" {
			offset = uint32(names[layEntry.Normal])
		}
		enc.Uint32(offset)

		enc.Bytes(versionPadding)
	}

	if enc.Error() != nil {
		return fmt.Errorf("encode: %w", enc.Error())
	}

	return nil

}
