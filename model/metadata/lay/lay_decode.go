package lay

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/tag"
)

// Decode loads a lay file
func Decode(model *common.Model, r io.ReadSeeker) error {
	var ok bool

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGL" {
		return fmt.Errorf("invalid header %s, wanted EQGL", header)
	}

	tag.New()

	model.Header.Version = int(dec.Uint32())
	versionOffset := 0
	switch model.Header.Version {
	case 2:
		versionOffset = 52 //32
	case 3:
		versionOffset = 16 //14
	case 4:
		versionOffset = 20
	default:
		return fmt.Errorf("unknown lay version: %d", model.Header.Version)
	}

	nameLength := int(dec.Uint32())
	layerCount := dec.Uint32()
	tag.Add(0, dec.Pos(), "red", "header")
	nameData := dec.Bytes(int(nameLength))
	tag.Add(tag.LastPos(), dec.Pos(), "green", "names")

	names := make(map[uint32]string)
	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[uint32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	for i := 0; i < int(layerCount); i++ {
		entryID := dec.Uint32()
		layer := &common.Layer{}

		if entryID != 0xffffffff {
			layer.Material, ok = names[entryID]
			if !ok {
				return fmt.Errorf("%d material 0x%x not found", i, entryID)
			}
		}

		entryID = dec.Uint32()
		if entryID != 0xffffffff {
			layer.Diffuse, ok = names[entryID]
			if !ok {
				return fmt.Errorf("%d diffuse 0x%x not found", i, entryID)
			}
		}

		entryID = dec.Uint32()
		if entryID != 0xffffffff {
			layer.Normal, ok = names[entryID]
			if !ok {
				return fmt.Errorf("%d normal 0x%x not found", i, entryID)
			}
		}
		dec.Bytes(versionOffset)
		//fmt.Println(hex.Dump())
		model.Layers = append(model.Layers, layer)
		tag.AddRand(tag.LastPos(), dec.Pos(), fmt.Sprintf("%d|%s|%s|%s", i, layer.Material, layer.Diffuse, layer.Normal))
	}

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	return nil
}
