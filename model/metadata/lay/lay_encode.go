package lay

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// Encode writes a lay file
func Encode(model *common.Model, w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.String("EQGL")

	enc.Uint32(uint32(model.Header.Version))
	var versionPadding []byte
	switch model.Header.Version {
	case 2:
		versionPadding = make([]byte, 52) //32
	case 3:
		versionPadding = make([]byte, 16) //14
	case 4:
		versionPadding = make([]byte, 20)
	default:
		return fmt.Errorf("unknown lay version: %d", model.Header.Version)
	}

	tmpNames := []string{}
	for _, lay := range model.Layers {
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

	enc.Uint32(uint32(len(nameData)))     // nameLength
	enc.Uint32(uint32(len(model.Layers))) //layerCount
	enc.Bytes(nameData)                   // nameData

	for _, layer := range model.Layers {
		offset := uint32(0xffffffff)
		if layer.Material != "" {
			offset = uint32(names[layer.Material])
		}
		enc.Uint32(offset)

		offset = uint32(0xffffffff)
		if layer.Diffuse != "" {
			offset = uint32(names[layer.Diffuse])
		}
		enc.Uint32(offset)

		offset = uint32(0xffffffff)
		if layer.Normal != "" {
			offset = uint32(names[layer.Normal])
		}
		enc.Uint32(offset)

		enc.Bytes(versionPadding)
	}

	if enc.Error() != nil {
		return fmt.Errorf("encode: %w", enc.Error())
	}

	return nil

}
