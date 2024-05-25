package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/tag"
)

type Ter struct {
	MetaFileName string      `yaml:"file_name"`
	Version      uint32      `yaml:"version"`
	Materials    []*Material `yaml:"materials"`
	Vertices     []Vertex    `yaml:"vertices"`
	Triangles    []Triangle  `yaml:"triangles"`
}

// Identity returns the type of the struct
func (ter *Ter) Identity() string {
	return "ter"
}

// Read reads a TER file
func (ter *Ter) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGT" {
		return fmt.Errorf("invalid header %s, wanted EQGT", header)
	}

	tag.New()
	ter.Version = dec.Uint32()

	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	verticesCount := dec.Uint32()
	triangleCount := dec.Uint32()
	tag.Add(0, dec.Pos(), "red", "header")
	nameData := dec.Bytes(int(nameLength))
	tag.Add(tag.LastPos(), dec.Pos(), "green", "names")

	names := make(map[int32]string)
	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[int32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	NameSet(names)

	//log.Debugf("names: %+v", names)

	nameCounter := 0
	for i := 0; i < int(materialCount); i++ {
		material := &Material{}
		material.ID = dec.Int32()
		nameCounter++

		material.Name = Name(dec.Int32())
		material.ShaderName = Name(dec.Int32())

		ter.Materials = append(ter.Materials, material)

		propertyCount := dec.Uint32()
		for j := 0; j < int(propertyCount); j++ {
			property := &MaterialProperty{
				Name: material.Name,
			}

			property.Name = Name(dec.Int32())

			property.Category = dec.Uint32()
			if property.Category == 0 {
				property.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Int32()
				if property.Category == 2 {
					property.Value = Name(val)
				} else {
					property.Value = fmt.Sprintf("%d", val)
				}
			}

			material.Properties = append(material.Properties, property)
		}
	}
	tag.Add(tag.LastPos(), dec.Pos(), "blue", "materials")

	for i := 0; i < int(verticesCount); i++ {
		v := Vertex{}
		v.Position.X = dec.Float32()
		v.Position.Y = dec.Float32()
		v.Position.Z = dec.Float32()
		v.Normal.X = dec.Float32()
		v.Normal.Y = dec.Float32()
		v.Normal.Z = dec.Float32()
		if ter.Version <= 2 {
			v.Tint = model.RGBA{R: 128, G: 128, B: 128, A: 255}
		} else {
			v.Tint = model.RGBA{R: dec.Uint8(), G: dec.Uint8(), B: dec.Uint8(), A: dec.Uint8()}
		}
		v.Uv.X = dec.Float32()
		v.Uv.Y = dec.Float32()
		if ter.Version <= 2 {
			v.Uv2.X = 0
			v.Uv2.Y = 0
		} else {
			v.Uv2.X = dec.Float32()
			v.Uv2.Y = dec.Float32()
		}

		ter.Vertices = append(ter.Vertices, v)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "yellow", "vertices")

	for i := 0; i < int(triangleCount); i++ {
		t := Triangle{}
		t.Index.X = dec.Uint32()
		t.Index.Y = dec.Uint32()
		t.Index.Z = dec.Uint32()

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

		t.Flag = dec.Uint32()
		ter.Triangles = append(ter.Triangles, t)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "purple", "triangles")

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
