package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type Ter struct {
	MetaFileName string
	Version      uint32
	Materials    []*ModMaterial
	Vertices     []*TerVertex
	Faces        []ModFace
	name         *eqgName
}

// TerVertex is a vertex
type TerVertex struct {
	Position [3]float32
	Normal   [3]float32
	Tint     [4]uint8
	Uv       [2]float32
	Uv2      [2]float32
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
	ter.Version = dec.Uint32()
	if header != "EQGT" {
		return fmt.Errorf("invalid header %s on version %d, wanted EQGT", header, ter.Version)
	}

	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	verticesCount := dec.Uint32()
	faceCount := dec.Uint32()
	nameData := dec.Bytes(int(nameLength))
	ter.name.parse(nameData)

	for i := 0; i < int(materialCount); i++ {
		material := &ModMaterial{}
		material.ID = dec.Int32()

		material.Name = ter.name.byOffset(dec.Int32())
		material.EffectName = ter.name.byOffset(dec.Int32())

		ter.Materials = append(ter.Materials, material)

		paramCount := dec.Uint32()
		for j := 0; j < int(paramCount); j++ {
			param := &ModMaterialParam{
				Name: material.Name,
			}

			param.Name = ter.name.byOffset(dec.Int32())

			param.Type = MaterialParamType(dec.Uint32())
			if param.Type == 0 {
				param.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Int32()
				if param.Type == 2 {
					param.Value = ter.name.byOffset(val)

				} else {
					param.Value = fmt.Sprintf("%d", val)
				}
			}
			material.Properties = append(material.Properties, param)
		}
	}

	for i := 0; i < int(verticesCount); i++ {
		v := &TerVertex{}
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

	for i := 0; i < int(faceCount); i++ {
		t := ModFace{}
		t.Index[0] = dec.Uint32()
		t.Index[1] = dec.Uint32()
		t.Index[2] = dec.Uint32()

		materialID := dec.Int32()

		var material *ModMaterial
		if materialID != -1 && len(ter.Materials) > int(materialID) {
			material = ter.Materials[materialID]
			t.MaterialName = material.Name
		}

		t.Flags = dec.Uint32()
		ter.Faces = append(ter.Faces, t)
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
