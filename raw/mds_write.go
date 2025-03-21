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
	if mds.name == nil {
		mds.name = &eqgName{}
	}
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGS")
	enc.Uint32(mds.Version)

	mds.name.clear()

	for _, material := range mds.Materials {
		mds.name.add(material.Name)
		mds.name.add(material.EffectName)
		for _, prop := range material.Properties {
			mds.name.add(prop.Name)
			switch prop.Type {
			case 2:
				mds.name.add(prop.Value)
			default:
			}
		}
	}

	for _, bone := range mds.Bones {
		mds.name.add(bone.Name)
	}

	for _, model := range mds.Models {
		mds.name.add(model.Name)
	}

	nameData := mds.name.data()
	enc.Uint32(uint32(len(nameData))) // nameLength
	enc.Uint32(uint32(len(mds.Materials)))

	enc.Uint32(uint32(len(mds.Bones)))
	enc.Uint32(uint32(len(mds.Models)))

	enc.Bytes(nameData)

	for _, material := range mds.Materials {
		enc.Int32(material.ID)
		enc.Uint32(uint32(mds.name.offsetByName(material.Name)))
		enc.Uint32(uint32(mds.name.offsetByName(material.EffectName)))
		enc.Uint32(uint32(len(material.Properties)))
		for _, prop := range material.Properties {
			enc.Uint32(uint32(mds.name.offsetByName(prop.Name)))
			enc.Uint32(uint32(prop.Type))
			switch prop.Type {
			case 0:
				fval, err := strconv.ParseFloat(prop.Value, 32)
				if err != nil {
					return err
				}
				enc.Float32(float32(fval))
			case 2:
				enc.Int32(mds.name.offsetByName(prop.Value))
			default:
				return err
			}
		}
	}

	for _, bone := range mds.Bones {
		enc.Int32(mds.name.offsetByName(bone.Name))
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

	//enc.Int32(mds.MainNameIndex)
	//enc.Int32(mds.SubNameIndex)

	//enc.Uint32(uint32(len(mds.Vertices)))
	//enc.Uint32(uint32(len(mds.Faces)))
	//enc.Uint32(uint32(len(mds.BoneAssignments)))

	for _, model := range mds.Models {
		enc.Uint32(model.MainPiece)
		enc.Int32(mds.name.offsetByName(model.Name))
		enc.Uint32(uint32(len(model.Vertices)))
		enc.Uint32(uint32(len(model.Faces)))
		enc.Uint32(model.BoneCount)

		for _, vert := range model.Vertices {
			enc.Float32(vert.Position[0])
			enc.Float32(vert.Position[1])
			enc.Float32(vert.Position[2])
			enc.Float32(vert.Normal[0])
			enc.Float32(vert.Normal[1])
			enc.Float32(vert.Normal[2])
			if mds.Version > 2 {
				enc.Uint8(vert.Tint[0])
				enc.Uint8(vert.Tint[1])
				enc.Uint8(vert.Tint[2])
				enc.Uint8(vert.Tint[3])
			}
			enc.Float32(vert.Uv[0])
			enc.Float32(vert.Uv[1])
			if mds.Version > 2 {
				enc.Float32(vert.Uv2[0])
				enc.Float32(vert.Uv2[1])
			}
		}

		for _, tri := range model.Faces {
			enc.Uint32(tri.Index[0])
			enc.Uint32(tri.Index[1])
			enc.Uint32(tri.Index[2])
			matID := int32(0)
			for _, mat := range mds.Materials {
				if mat.Name == tri.MaterialName {
					matID = mat.ID
					break
				}
			}
			enc.Int32(matID)
			enc.Uint32(tri.Flags)
		}

		if model.BoneCount > 0 {
			for i := 0; i < len(model.Vertices); i++ {
				vert := model.Vertices[i]
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

	}

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil

}
