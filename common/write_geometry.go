package common

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/xackery/quail/helper"
)

func WriteGeometry(materials []*Material, vertices []*Vertex, triangles []*Triangle) ([]byte, []byte, error) {
	var err error

	names := []string{}
	nameBuf := bytes.NewBuffer(nil)
	dataBuf := bytes.NewBuffer(nil)
	nameID := -1
	// materials

	for materialID, o := range materials {
		err = binary.Write(dataBuf, binary.LittleEndian, uint32(materialID))
		if err != nil {
			return nil, nil, fmt.Errorf("write material id %s: %w", o.Name, err)
		}

		nameID = -1
		for i, name := range names {
			if name == o.Name {
				nameID = i
				break
			}
		}
		if nameID == -1 {
			names = append(names, o.Name)
			nameID = len(names) - 1
			err = helper.WriteString(nameBuf, o.Name)
			if err != nil {
				return nil, nil, fmt.Errorf("writestring material %s: %w", o.Name, err)
			}
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameID))
		if err != nil {
			return nil, nil, fmt.Errorf("write material name id %s: %w", o.Name, err)
		}

		nameID = -1
		for i, name := range names {
			if name == o.ShaderName {
				nameID = i
				break
			}
		}
		if nameID == -1 {
			names = append(names, o.ShaderName)
			nameID = len(names) - 1
			err = helper.WriteString(nameBuf, o.ShaderName)
			if err != nil {
				return nil, nil, fmt.Errorf("writestring material %s shader: %w", o.Name, err)
			}
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameID))
		if err != nil {
			return nil, nil, fmt.Errorf("write material shader id %s: %w", o.Name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(len(o.Properties)))
		if err != nil {
			return nil, nil, fmt.Errorf("write material property count %s: %w", o.Name, err)
		}

		for propertyID, p := range o.Properties {
			nameID = -1
			for i, name := range names {
				if name == p.Name {
					nameID = i
					break
				}
			}
			if nameID == -1 {
				names = append(names, p.Name)
				nameID = len(names) - 1
				err = helper.WriteString(nameBuf, p.Name)
				if err != nil {
					return nil, nil, fmt.Errorf("writestring material %s property %s: %w", o.Name, p.Name, err)
				}
			}

			err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameID))
			if err != nil {
				return nil, nil, fmt.Errorf("write material %s property %s id %d: %w", o.Name, p.Name, propertyID, err)
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
		nameID = -1
		for i, name := range names {
			if name == o.MaterialName {
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
