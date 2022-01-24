package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type nameInfo struct {
	offset uint32
	name   string
}

func WriteGeometry(materials []*Material, vertices []*Vertex, triangles []*Triangle) ([]byte, []byte, error) {
	var err error

	names := []*nameInfo{}
	nameBuf := bytes.NewBuffer(nil)
	dataBuf := bytes.NewBuffer(nil)
	nameOffset := int32(-1)
	// materials

	tmpNames := []string{}
	for _, o := range materials {
		tmpNames = append(tmpNames, o.Name)
		tmpNames = append(tmpNames, o.ShaderName)
		for _, p := range o.Properties {
			tmpNames = append(tmpNames, p.Name)
		}
	}
	for _, name := range tmpNames {
		isNew := true
		for _, val := range names {
			if val.name == name {
				isNew = false
				break
			}
		}
		if !isNew {
			continue
		}

		names = append(names, &nameInfo{
			offset: uint32(nameBuf.Len()),
			name:   name,
		})
		_, err = nameBuf.Write([]byte(name))
		if err != nil {
			return nil, nil, fmt.Errorf("write name: %w", err)
		}
		_, err = nameBuf.Write([]byte{0})
		if err != nil {
			return nil, nil, fmt.Errorf("write 0: %w", err)
		}
	}

	//fmt.Println(hex.Dump(nameBuf.Bytes()))
	for materialID, o := range materials {
		err = binary.Write(dataBuf, binary.LittleEndian, uint32(materialID))
		if err != nil {
			return nil, nil, fmt.Errorf("write material id %s: %w", o.Name, err)
		}

		nameOffset = -1
		for _, val := range names {
			if val.name == o.Name {
				nameOffset = int32(val.offset)
				break
			}
		}
		if nameOffset == -1 {
			return nil, nil, fmt.Errorf("name %s not found", o.Name)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameOffset))
		if err != nil {
			return nil, nil, fmt.Errorf("write material name offset %s: %w", o.Name, err)
		}

		nameOffset = -1
		for _, val := range names {
			if val.name == o.ShaderName {
				nameOffset = int32(val.offset)
				break
			}
		}
		if nameOffset == -1 {
			return nil, nil, fmt.Errorf("shaderName %s not found", o.Name)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameOffset))
		if err != nil {
			return nil, nil, fmt.Errorf("write shader name offset %s: %w", o.ShaderName, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(len(o.Properties)))
		if err != nil {
			return nil, nil, fmt.Errorf("write material property count %s: %w", o.Name, err)
		}

		for propertyID, p := range o.Properties {
			nameOffset = -1
			for _, val := range names {
				if val.name == p.Name {
					nameOffset = int32(val.offset)
					break
				}
			}
			if nameOffset == -1 {
				return nil, nil, fmt.Errorf("%s prop %s not found", o.Name, p.Name)
			}

			err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameOffset))
			if err != nil {
				return nil, nil, fmt.Errorf("write %s property %s name offset: %w", o.Name, p.Name, err)
			}

			err = binary.Write(dataBuf, binary.LittleEndian, p.TypeValue)
			if err != nil {
				return nil, nil, fmt.Errorf("write %s property %s type: %w", o.Name, p.Name, err)
			}

			if p.TypeValue == 0 {
				err = binary.Write(dataBuf, binary.LittleEndian, p.FloatValue)
				if err != nil {
					return nil, nil, fmt.Errorf("write material %s property %s id %d value (float): %w", o.Name, p.Name, propertyID, err)
				}
			} else {
				err = binary.Write(dataBuf, binary.LittleEndian, p.IntValue)
				if err != nil {
					return nil, nil, fmt.Errorf("write material %s property %s id %d value (int): %w", o.Name, p.Name, propertyID, err)
				}
			}
		}
	}

	// verts
	for i, o := range vertices {
		err = binary.Write(dataBuf, binary.LittleEndian, o.Position)
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d position: %w", i, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.Normal)
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d normal: %w", i, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.Uv)
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d uv: %w", i, err)
		}
	}

	// triangles
	for i, o := range triangles {
		nameID := -1
		for i, val := range names {
			if val.name == o.MaterialName {
				nameID = i
				break
			}
		}
		if nameID == -1 {
			return nil, nil, fmt.Errorf("triangle %d refers to material %s, which is not declared", i, o.MaterialName)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.Index)
		if err != nil {
			return nil, nil, fmt.Errorf("write triangle %d index: %w", i, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameID))
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d position2: %w", i, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.Flag)
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d flag: %w", i, err)
		}
	}
	return nameBuf.Bytes(), dataBuf.Bytes(), nil
}
