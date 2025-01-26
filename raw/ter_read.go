package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type Ter struct {
	MetaFileName string      `yaml:"file_name"`
	Version      uint32      `yaml:"version"`
	Materials    []*Material `yaml:"materials"`
	Vertices     []Vertex    `yaml:"vertices"`
	Triangles    []Face      `yaml:"triangles"`
	name         *eqgName
}

// Identity returns the type of the struct
func (ter *Ter) Identity() string {
	return "ter"
}

// Read reads a TER file
func (ter *Ter) Read(r io.ReadSeeker) error {
	if ter.name == nil {
		ter.name = &eqgName{}
	}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGT" {
		return fmt.Errorf("invalid header %s, wanted EQGT", header)
	}

	ter.Version = dec.Uint32()

	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	verticesCount := dec.Uint32()
	triangleCount := dec.Uint32()
	nameData := dec.Bytes(int(nameLength))
	ter.name.parse(nameData)

	nameCounter := 0
	for i := 0; i < int(materialCount); i++ {
		material := &Material{}
		material.ID = dec.Int32()
		nameCounter++

		material.Name = ter.name.byOffset(dec.Int32())
		material.EffectName = ter.name.byOffset(dec.Int32())

		ter.Materials = append(ter.Materials, material)

		propertyCount := dec.Uint32()
		for j := 0; j < int(propertyCount); j++ {
			property := &MaterialParam{
				Name: material.Name,
			}

			property.Name = ter.name.byOffset(dec.Int32())

			property.Type = MaterialParamType(dec.Uint32())
			if property.Type == 0 {
				property.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Int32()
				if property.Type == 2 {
					property.Value = ter.name.byOffset(val)
				} else {
					property.Value = fmt.Sprintf("%d", val)
				}
			}

			material.Properties = append(material.Properties, property)
		}
	}

	for i := 0; i < int(verticesCount); i++ {
		v := Vertex{}
		v.Position[0] = dec.Float32()
		v.Position[1] = dec.Float32()
		v.Position[2] = dec.Float32()
		v.Normal[0] = dec.Float32()
		v.Normal[1] = dec.Float32()
		v.Normal[2] = dec.Float32()
		if ter.Version <= 2 {
			v.Tint = [4]uint8{128, 128, 128, 255}
		} else {
			v.Tint = [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()}
		}
		v.Uv[0] = dec.Float32()
		v.Uv[1] = dec.Float32()
		if ter.Version <= 2 {
			v.Uv2[0] = 0
			v.Uv2[1] = 0
		} else {
			v.Uv2[0] = dec.Float32()
			v.Uv2[1] = dec.Float32()
		}

		ter.Vertices = append(ter.Vertices, v)
	}

	for i := 0; i < int(triangleCount); i++ {
		t := Face{}
		t.Index[0] = dec.Uint32()
		t.Index[1] = dec.Uint32()
		t.Index[2] = dec.Uint32()

		materialID := dec.Int32()

		var material *Material
		for _, mat := range ter.Materials {
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

		t.Flags = dec.Uint32()
		ter.Triangles = append(ter.Triangles, t)
	}

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	return nil
}

// SetFileName sets the name of the file
func (ter *Ter) SetFileName(name string) {
	ter.MetaFileName = name
}

// FileName returns the name of the file
func (ter *Ter) FileName() string {
	return ter.MetaFileName
}
