package lay

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/tag"
)

// Decode loads a lay file
func Decode(mesh *common.Model, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGL" {
		return fmt.Errorf("invalid header %s, wanted EQGL", header)
	}

	tag.New()

	version := dec.Uint32()

	fmt.Println("version is", version)
	versionOffset := 0
	switch version {
	case 2:
		versionOffset = 40 //32
	case 3:
		versionOffset = 18 //14
	case 4:
		versionOffset = 20
	default:
		return fmt.Errorf("unknown lay version: %d", version)
	}

	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	tag.Add(0, int(dec.Pos()-1), "red", "header")
	nameData := dec.Bytes(int(nameLength))
	tag.Add(tag.LastPos(), int(dec.Pos()), "green", "names")

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

	for i := 0; i < int(materialCount); i++ {
		entryID := dec.Uint32()

		name, ok := names[entryID]
		if !ok {
			return fmt.Errorf("%d names materialID 0x%x not found", i, entryID)
		}
		entryID = dec.Uint32()
		colorTexture, ok := names[entryID]
		if !ok {
			return fmt.Errorf("%d names colorTexture 0x%x not found", i, entryID)
		}

		normalTexture := ""
		entryID = dec.Uint32()
		if entryID != 0xffffffff {
			normalTexture, ok = names[entryID]
			if !ok {
				return fmt.Errorf("%d names normalTexture 0x%x not found", i, entryID)
			}
		}

		print("name: ", name, " colorTexture: ", colorTexture, " normalTexture: ", normalTexture, "\n")

		fmt.Println(hex.Dump(dec.Bytes(versionOffset)))
	}

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	return nil
}
