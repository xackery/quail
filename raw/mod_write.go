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

	mod.NameClear()

	for _, material := range mod.Materials {
		mod.NameAdd(material.Name)
		mod.NameAdd(material.ShaderName)
		for _, prop := range material.Properties {
			mod.NameAdd(prop.Name)
			switch prop.Category {
			case 2:
				mod.NameAdd(prop.Value)
			default:
			}
		}
	}

	for _, bone := range mod.Bones {
		mod.NameAdd(bone.Name)
	}

	nameData := mod.NameData()
	enc.Uint32(uint32(len(nameData))) // nameLength
	enc.Uint32(uint32(len(mod.Materials)))
	enc.Uint32(uint32(len(mod.Vertices)))
	enc.Uint32(uint32(len(mod.Triangles)))
	enc.Uint32(uint32(len(mod.Bones)))
	enc.Bytes(nameData)

	for _, material := range mod.Materials {
		enc.Int32(material.ID)
		enc.Uint32(uint32(mod.NameIndex(material.Name)))
		enc.Uint32(uint32(mod.NameIndex(material.ShaderName)))
		enc.Uint32(uint32(len(material.Properties)))
		for _, prop := range material.Properties {
			enc.Uint32(uint32(mod.NameIndex(prop.Name)))
			enc.Uint32(uint32(prop.Category))
			switch prop.Category {
			case 0:
				fval, err := strconv.ParseFloat(prop.Value, 32)
				if err != nil {
					return fmt.Errorf("parse float: %w", err)
				}
				enc.Float32(float32(fval))
			case 2:
				enc.Int32(mod.NameIndex(prop.Value))
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
		enc.Float32(vert.Position[0])
		enc.Float32(vert.Position[1])
		enc.Float32(vert.Position[2])
		enc.Float32(vert.Normal[0])
		enc.Float32(vert.Normal[1])
		enc.Float32(vert.Normal[2])
		if mod.Version > 2 {
			enc.Uint8(vert.Tint[0])
			enc.Uint8(vert.Tint[1])
			enc.Uint8(vert.Tint[2])
			enc.Uint8(vert.Tint[3])
		}
		enc.Float32(vert.Uv[0])
		enc.Float32(vert.Uv[1])
		if mod.Version > 2 {
			enc.Float32(vert.Uv2[0])
			enc.Float32(vert.Uv2[1])
		}
	}

	for _, tri := range mod.Triangles {
		enc.Uint32(tri.Index[0])
		enc.Uint32(tri.Index[1])
		enc.Uint32(tri.Index[2])
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
		enc.Int32(mod.NameIndex(bone.Name))
		enc.Int32(bone.Next)
		enc.Uint32(bone.ChildrenCount)
		enc.Int32(bone.ChildIndex)
		enc.Float32(bone.Pivot[0])
		enc.Float32(bone.Pivot[1])
		enc.Float32(bone.Pivot[2])
		enc.Float32(bone.Rotation[0])
		enc.Float32(bone.Rotation[1])
		enc.Float32(bone.Rotation[2])
		enc.Float32(bone.Scale[0])
		enc.Float32(bone.Scale[1])
		enc.Float32(bone.Scale[2])
		enc.Float32(bone.Scale2)
	}

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
