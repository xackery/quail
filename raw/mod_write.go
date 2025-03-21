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
	if mod.name == nil {
		mod.name = &eqgName{}
	}
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGM")
	enc.Uint32(mod.Version)

	mod.name.clear()

	for _, material := range mod.Materials {
		mod.name.add(material.Name)
		mod.name.add(material.EffectName)
		for _, prop := range material.Properties {
			mod.name.add(prop.Name)
			switch prop.Type {
			case 2:
				mod.name.add(prop.Value)
			default:
			}
		}
	}

	for _, bone := range mod.Bones {
		mod.name.add(bone.Name)
	}

	nameData := mod.name.data()
	enc.Uint32(uint32(len(nameData))) // nameLength
	enc.Uint32(uint32(len(mod.Materials)))
	enc.Uint32(uint32(len(mod.Vertices)))
	enc.Uint32(uint32(len(mod.Faces)))
	enc.Uint32(uint32(len(mod.Bones)))
	enc.Bytes(nameData)

	for _, material := range mod.Materials {
		enc.Int32(material.ID)
		enc.Uint32(uint32(mod.name.offsetByName(material.Name)))
		enc.Uint32(uint32(mod.name.offsetByName(material.EffectName)))
		enc.Uint32(uint32(len(material.Properties)))
		for _, prop := range material.Properties {
			enc.Uint32(uint32(mod.name.offsetByName(prop.Name)))
			enc.Uint32(uint32(prop.Type))
			switch prop.Type {
			case 0:
				fval, err := strconv.ParseFloat(prop.Value, 32)
				if err != nil {
					return fmt.Errorf("parse float: %w", err)
				}
				enc.Float32(float32(fval))
			case 2:
				enc.Int32(mod.name.offsetByName(prop.Value))
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

	for _, tri := range mod.Faces {
		enc.Uint32(tri.Index[0])
		enc.Uint32(tri.Index[1])
		enc.Uint32(tri.Index[2])
		matID := int32(0)
		for i, mat := range mod.Materials {
			if mat.Name == tri.MaterialName {
				matID = int32(i)
				break
			}
		}
		enc.Int32(matID)
		enc.Uint32(tri.Flags)
	}

	for _, bone := range mod.Bones {
		enc.Int32(mod.name.offsetByName(bone.Name))
		enc.Int32(bone.Next)
		enc.Uint32(bone.ChildrenCount)
		enc.Int32(bone.ChildIndex)
		enc.Float32(bone.Pivot[0])
		enc.Float32(bone.Pivot[1])
		enc.Float32(bone.Pivot[2])
		enc.Float32(bone.Quaternion[0])
		enc.Float32(bone.Quaternion[1])
		enc.Float32(bone.Quaternion[2])
		enc.Float32(bone.Quaternion[3])
		enc.Float32(bone.Scale[0])
		enc.Float32(bone.Scale[1])
		enc.Float32(bone.Scale[2])
	}

	if len(mod.Bones) > 0 {
		for i := 0; i < len(mod.Vertices); i++ {
			vert := mod.Vertices[i]
			enc.Int32(int32(len(vert.Weights)))
			for j := 0; j < int(4); j++ {
				if j >= len(vert.Weights) {
					enc.Int32(0)
					enc.Float32(0)
					continue
				}

				enc.Int32(vert.Weights[j].BoneIndex)
				enc.Float32(vert.Weights[j].Value)
			}
		}
	}

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
