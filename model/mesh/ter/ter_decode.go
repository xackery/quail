package ter

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/tag"
)

// Decode decodes a TER file
func Decode(model *common.Model, r io.ReadSeeker) error {
	var ok bool
	model.FileType = "ter"

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGT" {
		return fmt.Errorf("invalid header %s, wanted EQGM", header)
	}

	tag.New()
	version := dec.Uint32()
	model.Version = int(version)

	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	verticesCount := dec.Uint32()
	triangleCount := dec.Uint32()
	tag.Add(0, dec.Pos(), "red", "header")
	nameData := dec.Bytes(int(nameLength))
	tag.Add(tag.LastPos(), dec.Pos(), "green", "names")

	names := make(map[uint32]string)
	modelNames := []string{}
	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[uint32(lastOffset)] = string(chunk)
			modelNames = append(modelNames, string(chunk))
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	//log.Debugf("names: %+v", names)

	if model.Name == "" {
		model.Name = modelNames[0]
	}

	nameCounter := 0
	for i := 0; i < int(materialCount); i++ {
		material := &common.Material{}
		material.ID = dec.Int32()
		nameOffset := dec.Uint32()
		nameCounter++
		material.Name, ok = names[nameOffset]
		if !ok {
			return fmt.Errorf("material nameOffset %d not found", nameOffset)
		}
		shaderOffset := dec.Uint32()
		material.ShaderName, ok = names[shaderOffset]
		if !ok {
			return fmt.Errorf("material shader not found")
		}

		isNew := true
		for _, mat := range model.Materials {
			if mat.Name == material.Name {
				material = mat
				isNew = false
				break
			}
		}
		if isNew {
			model.Materials = append(model.Materials, material)
		}

		propertyCount := dec.Uint32()
		for j := 0; j < int(propertyCount); j++ {
			property := &common.MaterialProperty{
				Name: material.Name,
			}

			propertyNameOffset := dec.Uint32()
			nameCounter++
			property.Name, ok = names[propertyNameOffset]
			if !ok {
				return fmt.Errorf("material property name not found")
			}

			property.Category = dec.Uint32()
			if property.Category == 0 {
				property.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Uint32()
				nameCounter++
				if property.Category == 2 {
					property.Value, ok = names[val]
					if !ok {
						return fmt.Errorf("material property value %d not found", val)
					}
				} else {
					property.Value = fmt.Sprintf("%d", val)
				}
			}

			isNew := true
			for _, prop := range material.Properties {
				if prop.Name == property.Name {
					property = prop
					isNew = false
					break
				}
			}

			if isNew {
				material.Properties = append(material.Properties, property)
			}
		}
	}
	tag.Add(tag.LastPos(), dec.Pos(), "blue", "materials")

	for i := 0; i < int(verticesCount); i++ {
		v := common.Vertex{}
		v.Position.X = dec.Float32()
		v.Position.Y = dec.Float32()
		v.Position.Z = dec.Float32()
		v.Normal.X = dec.Float32()
		v.Normal.Y = dec.Float32()
		v.Normal.Z = dec.Float32()
		if version <= 2 {
			v.Tint = common.RGBA{R: 128, G: 128, B: 128, A: 255}
		} else {
			v.Tint = common.RGBA{R: dec.Uint8(), G: dec.Uint8(), B: dec.Uint8(), A: dec.Uint8()}
		}
		v.Uv.X = dec.Float32()
		v.Uv.Y = dec.Float32()
		if version <= 2 {
			v.Uv2.X = 0
			v.Uv2.Y = 0
		} else {
			v.Uv2.X = dec.Float32()
			v.Uv2.Y = dec.Float32()
		}

		model.Vertices = append(model.Vertices, v)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "yellow", "vertices")

	for i := 0; i < int(triangleCount); i++ {
		t := common.Triangle{}
		t.Index.X = dec.Uint32()
		t.Index.Y = dec.Uint32()
		t.Index.Z = dec.Uint32()

		materialID := dec.Int32()

		var material *common.Material
		for _, mat := range model.Materials {
			if mat.ID == materialID {
				material = mat
				break
			}
		}
		if material == nil {
			//if materialID != -1 {
			//log.Warnf("material %d not found", materialID)
			//return fmt.Errorf("material %d not found", materialID)
			//}
			t.MaterialName = ""
		} else {
			t.MaterialName = material.Name
		}

		t.Flag = dec.Uint32()
		model.Triangles = append(model.Triangles, t)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "purple", "triangles")

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	model.Name = strings.ToLower(model.Name)

	log.Debugf("%s (ter) decoded %d verts, %d triangles, %d bones, %d materials", model.Name, len(model.Vertices), len(model.Triangles), len(model.Bones), len(model.Materials))
	return nil
}
