package def

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/xackery/encdec"
)

// MaterialBuild prepares an EQG-styled material buffer list
func (mesh *Mesh) materialBuild(names map[string]int32) ([]byte, error) {
	var err error

	dataBuf := bytes.NewBuffer(nil)
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)
	nameOffset := int32(-1)
	for materialID, o := range mesh.Materials {
		enc.Uint32(uint32(materialID))

		nameOffset = -1
		for key, offset := range names {
			if key == o.Name {
				nameOffset = offset
				break
			}
		}
		if nameOffset == -1 {
			//log.Debugf("material %s not found ignoring", o.Name)
		}

		enc.Uint32(uint32(nameOffset))

		nameOffset = -1
		for key, offset := range names {
			if key == o.ShaderName {
				nameOffset = offset
				break
			}
		}
		if nameOffset == -1 {
			return nil, fmt.Errorf("shaderName %s not found", o.Name)
		}

		enc.Uint32(uint32(nameOffset))

		enc.Uint32(uint32(len(o.Properties)))

		for _, p := range o.Properties {
			nameOffset = -1
			for key, offset := range names {
				if key == p.Name {
					nameOffset = offset
					break
				}
			}
			if nameOffset == -1 {
				return nil, fmt.Errorf("%s prop %s not found", o.Name, p.Name)
			}

			enc.Uint32(uint32(nameOffset))
			enc.Uint32(p.Category)

			nameOffset = -1

			err = materialPropertyWrite(dataBuf, p.Value, names)
			if err != nil {
				return nil, fmt.Errorf("writePropertyValue: %w", err)
			}
		}
	}
	return dataBuf.Bytes(), nil
}

func materialPropertyWrite(buf *bytes.Buffer, value string, names map[string]int32) error {
	enc := encdec.NewEncoder(buf, binary.LittleEndian)
	val, err := strconv.Atoi(value)
	if err == nil {
		enc.Uint32(uint32(val))
		return nil
	}

	fVal, err := strconv.ParseFloat(value, 64)
	if err == nil {
		enc.Float32(float32(fVal))
		return nil
	}
	nameOffset := int32(-1)
	for key, offset := range names {
		if key == value {
			nameOffset = offset
			break
		}
	}
	if nameOffset == -1 {
		return fmt.Errorf("value %s: %w", value, err)
	}
	enc.Int32(nameOffset)
	return nil
}
