package raw

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/xackery/encdec"
)

// Encode writes a mod file
func (mod *Mod) Write(w io.Writer) error {
	var err error
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGM")
	enc.Uint32(mod.Version)

	NameClear()

	for _, material := range mod.Materials {
		NameAdd(material.Name)
		NameAdd(material.ShaderName)
		for _, prop := range material.Properties {
			NameAdd(prop.Name)
			switch prop.Category {
			case 2:
				NameAdd(prop.Value)
			default:
			}
		}
	}

	for _, bone := range mod.Bones {
		NameAdd(bone.Name)
	}

	nameData := NameData()
	enc.Uint32(uint32(len(nameData))) // nameLength
	enc.Uint32(uint32(len(mod.Materials)))
	enc.Uint32(uint32(len(mod.Vertices)))
	enc.Uint32(uint32(len(mod.Triangles)))
	enc.Uint32(uint32(len(mod.Bones)))
	enc.Bytes(nameData)

	for _, material := range mod.Materials {
		enc.Int32(material.ID)
		enc.Uint32(uint32(NameIndex(material.Name)))
		enc.Uint32(uint32(NameIndex(material.ShaderName)))
		enc.Uint32(uint32(len(material.Properties)))
		for _, prop := range material.Properties {
			enc.Uint32(uint32(NameIndex(prop.Name)))
			enc.Uint32(uint32(prop.Category))
			switch prop.Category {
			case 0:
				fval, err := strconv.ParseFloat(prop.Value, 32)
				if err != nil {
					return fmt.Errorf("parse float: %w", err)
				}
				enc.Float32(float32(fval))
			case 2:
				enc.Int32(NameIndex(prop.Value))
			default:
				val, err := strconv.Atoi(prop.Value)
				if err != nil {
					return fmt.Errorf("parse int: %w", err)
				}
				enc.Int32(int32(val))
			}
		}
	}

	for _, vert := range mod.Vertices {
		enc.Float32(vert.Position.X)
		enc.Float32(vert.Position.Y)
		enc.Float32(vert.Position.Z)
		enc.Float32(vert.Normal.X)
		enc.Float32(vert.Normal.Y)
		enc.Float32(vert.Normal.Z)
		if mod.Version > 2 {
			enc.Uint8(vert.Tint.R)
			enc.Uint8(vert.Tint.G)
			enc.Uint8(vert.Tint.B)
			enc.Uint8(vert.Tint.A)
		}
		enc.Float32(vert.Uv.X)
		enc.Float32(vert.Uv.Y)
		if mod.Version > 2 {
			enc.Float32(vert.Uv2.X)
			enc.Float32(vert.Uv2.Y)
		}
	}

	for _, tri := range mod.Triangles {
		enc.Uint32(tri.Index.X)
		enc.Uint32(tri.Index.Y)
		enc.Uint32(tri.Index.Z)
		matID := int32(0)
		for _, mat := range mod.Materials {
			if mat.Name == tri.MaterialName {
				matID = mat.ID
				break
			}
		}
		enc.Int32(matID)
		enc.Uint32(tri.Flag)
	}

	for _, bone := range mod.Bones {
		enc.Int32(NameIndex(bone.Name))
		enc.Int32(bone.Next)
		enc.Uint32(bone.ChildrenCount)
		enc.Int32(bone.ChildIndex)
		enc.Float32(bone.Pivot.X)
		enc.Float32(bone.Pivot.Y)
		enc.Float32(bone.Pivot.Z)
		enc.Float32(bone.Rotation.X)
		enc.Float32(bone.Rotation.Y)
		enc.Float32(bone.Rotation.Z)
		enc.Float32(bone.Scale.X)
		enc.Float32(bone.Scale.Y)
		enc.Float32(bone.Scale.Z)
		enc.Float32(bone.Scale2)
	}

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
