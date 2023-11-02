package raw

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/xackery/encdec"
)

// Write writes a mds file
func (mds *Mds) Write(w io.Writer) error {
	var err error
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGS")
	enc.Uint32(mds.Version)

	NameClear()

	for _, material := range mds.Materials {
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

	for _, bone := range mds.Bones {
		NameAdd(bone.Name)
	}

	nameData := NameData()
	enc.Uint32(uint32(len(nameData))) // nameLength
	enc.Uint32(uint32(len(mds.Materials)))
	enc.Uint32(uint32(len(mds.Bones)))
	enc.Uint32(uint32(len(mds.Subs)))
	enc.Bytes(nameData)

	for _, material := range mds.Materials {
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
					return err
				}
				enc.Float32(float32(fval))
			case 2:
				enc.Int32(NameIndex(prop.Value))
			default:
				return err
			}
		}
	}

	for _, bone := range mds.Bones {
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

	enc.Int32(mds.MainNameIndex)
	enc.Int32(mds.SubNameIndex)

	enc.Uint32(uint32(len(mds.Vertices)))
	enc.Uint32(uint32(len(mds.Triangles)))
	enc.Uint32(uint32(len(mds.BoneAssignments)))

	for _, vert := range mds.Vertices {
		enc.Float32(vert.Position.X)
		enc.Float32(vert.Position.Y)
		enc.Float32(vert.Position.Z)
		enc.Float32(vert.Normal.X)
		enc.Float32(vert.Normal.Y)
		enc.Float32(vert.Normal.Z)
		if mds.Version > 2 {
			enc.Uint8(vert.Tint.R)
			enc.Uint8(vert.Tint.G)
			enc.Uint8(vert.Tint.B)
			enc.Uint8(vert.Tint.A)
		}
		enc.Float32(vert.Uv.X)
		enc.Float32(vert.Uv.Y)
		if mds.Version > 2 {
			enc.Float32(vert.Uv2.X)
			enc.Float32(vert.Uv2.Y)
		}
	}

	for _, tri := range mds.Triangles {
		enc.Uint32(tri.Index.X)
		enc.Uint32(tri.Index.Y)
		enc.Uint32(tri.Index.Z)
		matID := int32(0)
		for _, mat := range mds.Materials {
			if mat.Name == tri.MaterialName {
				matID = mat.ID
				break
			}
		}
		enc.Int32(matID)
		enc.Uint32(tri.Flag)
	}

	// TODO: sub count
	// TODO: bone assigment count

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil

}
