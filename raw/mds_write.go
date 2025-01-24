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

	mds.NameClear()

	for _, material := range mds.Materials {
		mds.NameAdd(material.Name)
		mds.NameAdd(material.EffectName)
		for _, prop := range material.Properties {
			mds.NameAdd(prop.Name)
			switch prop.Type {
			case 2:
				mds.NameAdd(prop.Value)
			default:
			}
		}
	}

	for _, bone := range mds.Bones {
		mds.NameAdd(bone.Name)
	}

	nameData := mds.NameData()
	enc.Uint32(uint32(len(nameData))) // nameLength
	enc.Uint32(uint32(len(mds.Materials)))

	enc.Uint32(uint32(len(mds.Bones)))
	enc.Uint32(uint32(len(mds.Models)))

	enc.Bytes(nameData)

	for _, material := range mds.Materials {
		enc.Int32(material.ID)
		enc.Uint32(uint32(mds.NameIndex(material.Name)))
		enc.Uint32(uint32(mds.NameIndex(material.EffectName)))
		enc.Uint32(uint32(len(material.Properties)))
		for _, prop := range material.Properties {
			enc.Uint32(uint32(mds.NameIndex(prop.Name)))
			enc.Uint32(uint32(prop.Type))
			switch prop.Type {
			case 0:
				fval, err := strconv.ParseFloat(prop.Value, 32)
				if err != nil {
					return err
				}
				enc.Float32(float32(fval))
			case 2:
				enc.Int32(mds.NameIndex(prop.Value))
			default:
				return err
			}
		}
	}

	for _, bone := range mds.Bones {
		enc.Int32(mds.NameIndex(bone.Name))
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
		enc.Int32(mds.NameIndex(model.Name))
		enc.Uint32(uint32(len(model.Vertices)))
		enc.Uint32(uint32(len(model.Faces)))
		enc.Uint32(uint32(len(model.BoneAssignments)))

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

		for _, weights := range model.BoneAssignments {
			enc.Uint32(uint32(len(weights)))
			for i := 0; i < 4; i++ {
				if i < len(weights) {
					enc.Int32(weights[i].BoneIndex)
					enc.Float32(weights[i].Value)
				} else {
					enc.Int32(0)
					enc.Float32(0)
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
