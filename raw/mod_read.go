package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/tag"
)

type Mod struct {
	Version   uint32      `yaml:"version"`
	Materials []*Material `yaml:"materials"`
	Bones     []*Bone     `yaml:"bones"`
	Vertices  []*Vertex   `yaml:"vertices"`
	Triangles []Triangle  `yaml:"triangles"`
}

// Decode reads a MOD file
func (mod *Mod) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGM" {
		return fmt.Errorf("invalid header %s, wanted EQGM", header)
	}

	tag.New()
	mod.Version = dec.Uint32()

	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	verticesCount := dec.Uint32()
	triangleCount := dec.Uint32()
	bonesCount := dec.Uint32()
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

	NamesSet(names)

	//log.Debugf("names: %+v", names)

	for i := 0; i < int(materialCount); i++ {
		material := &Material{}
		material.ID = dec.Int32()
		material.Name = Name(dec.Int32())
		material.ShaderName = Name(dec.Int32())
		mod.Materials = append(mod.Materials, material)

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
		v := &Vertex{}
		v.Position.X = dec.Float32()
		v.Position.Y = dec.Float32()
		v.Position.Z = dec.Float32()
		v.Normal.X = dec.Float32()
		v.Normal.Y = dec.Float32()
		v.Normal.Z = dec.Float32()
		if mod.Version <= 2 {
			v.Tint = RGBA{R: 128, G: 128, B: 128, A: 255}
		} else {
			v.Tint = RGBA{R: dec.Uint8(), G: dec.Uint8(), B: dec.Uint8(), A: dec.Uint8()}
		}
		v.Uv.X = dec.Float32()
		v.Uv.Y = dec.Float32()

		if mod.Version <= 2 {
			v.Uv2.X = 0
			v.Uv2.Y = 0
		} else {
			v.Uv2.X = dec.Float32()
			v.Uv2.Y = dec.Float32()
		}

		mod.Vertices = append(mod.Vertices, v)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "yellow", "vertices")

	for i := 0; i < int(triangleCount); i++ {
		t := Triangle{}
		t.Index.X = dec.Uint32()
		t.Index.Y = dec.Uint32()
		t.Index.Z = dec.Uint32()

		materialID := dec.Int32()

		var material *Material
		for _, mat := range mod.Materials {
			if mat.ID == materialID {
				material = mat
				break
			}
		}
		if material == nil {
			if materialID != -1 {
				log.Debugf("material %d not found", materialID)
				//return fmt.Errorf("material %d not found", materialID)
			}
			t.MaterialName = ""
		} else {
			t.MaterialName = material.Name
		}

		t.Flag = dec.Uint32()
		mod.Triangles = append(mod.Triangles, t)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "purple", "triangles")

	for i := 0; i < int(bonesCount); i++ {
		bone := &Bone{}
		bone.Name = Name(dec.Int32())
		bone.Next = dec.Int32()
		bone.ChildrenCount = dec.Uint32()
		bone.ChildIndex = dec.Int32()
		bone.Pivot.X = dec.Float32()
		bone.Pivot.Y = dec.Float32()
		bone.Pivot.Z = dec.Float32()
		bone.Rotation.X = dec.Float32()
		bone.Rotation.Y = dec.Float32()
		bone.Rotation.Z = dec.Float32()
		bone.Scale.X = dec.Float32()
		bone.Scale.Y = dec.Float32()
		bone.Scale.Z = dec.Float32()
		bone.Scale2 = dec.Float32()

		mod.Bones = append(mod.Bones, bone)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "orange", "bones")

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	return nil
}
