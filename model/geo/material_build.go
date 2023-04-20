package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// MaterialBuild prepares an EQG-styled material buffer list
func MaterialBuild(version uint32, names map[string]int32, matManager *MaterialManager) ([]byte, error) {
	var err error
	dataBuf := bytes.NewBuffer(nil)
	nameOffset := int32(-1)
	for materialID, o := range matManager.materials {
		err = binary.Write(dataBuf, binary.LittleEndian, uint32(materialID))
		if err != nil {
			return nil, fmt.Errorf("write material id %s: %w", o.Name, err)
		}

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

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameOffset))
		if err != nil {
			return nil, fmt.Errorf("write material name offset %s: %w", o.Name, err)
		}

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

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameOffset))
		if err != nil {
			return nil, fmt.Errorf("write shader name offset %s: %w", o.ShaderName, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(len(o.Properties)))
		if err != nil {
			return nil, fmt.Errorf("write material property count %s: %w", o.Name, err)
		}

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

			err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameOffset))
			if err != nil {
				return nil, fmt.Errorf("write %s property %s name offset: %w", o.Name, p.Name, err)
			}

			err = binary.Write(dataBuf, binary.LittleEndian, p.Category)
			if err != nil {
				return nil, fmt.Errorf("write %s property %s type: %w", o.Name, p.Name, err)
			}

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
	val, err := strconv.Atoi(value)
	if err == nil {
		err = binary.Write(buf, binary.LittleEndian, uint32(val))
		if err != nil {
			return fmt.Errorf("write int %s: %w", value, err)
		}
		return nil
	}

	fVal, err := strconv.ParseFloat(value, 64)
	if err == nil {
		err = binary.Write(buf, binary.LittleEndian, float32(fVal))
		if err != nil {
			return fmt.Errorf("write float %s: %w", value, err)
		}
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
	err = binary.Write(buf, binary.LittleEndian, int32(nameOffset))
	if err != nil {
		return fmt.Errorf("write property offset: %w", err)
	}
	return nil
}
