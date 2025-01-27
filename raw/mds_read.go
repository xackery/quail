package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type Mds struct {
	MetaFileName string
	Version      uint32
	Materials    []*ModMaterial
	Bones        []*ModBone
	//MainNameIndex int32
	//SubNameIndex  int32
	Models []*MdsModel
	name   *eqgName
}

func (mds *Mds) Identity() string {
	return "mds"
}

type MdsModel struct {
	MainPiece       uint32 // 0: no, 1: yes, head is a mainpiece
	Name            string
	Vertices        []*ModVertex
	Faces           []*ModFace
	BoneAssignments [][4]*MdsBoneWeight
}

type MdsBoneWeight struct {
	BoneIndex int32
	Value     float32
}

func (mds *Mds) String() string {
	if mds.name == nil {
		mds.name = &eqgName{}
	}
	out := fmt.Sprintf("Mds: %s,", mds.MetaFileName)
	out += fmt.Sprintf(" %d names,", mds.name.len())
	out += fmt.Sprintf(" %d materials", len(mds.Materials))
	if len(mds.Materials) > 0 {
		out += " ["

		for i, material := range mds.Materials {
			out += material.Name
			if i < len(mds.Materials)-1 {
				out += ", "
			}
		}
		out += "]"

	}

	out += fmt.Sprintf(", %d bones,", len(mds.Bones))
	out += fmt.Sprintf(" %d models", len(mds.Models))
	if len(mds.Models) > 0 {
		out += " ["

		for i, model := range mds.Models {
			out += model.Name
			if i < len(mds.Models)-1 {
				out += ", "
			}
		}
		out += "]"

	}
	return out
}

// Read reads a mds file
func (mds *Mds) Read(r io.ReadSeeker) error {
	if mds.name == nil {
		mds.name = &eqgName{}
	}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGS" {
		return fmt.Errorf("invalid header %s, wanted EQGS", header)
	}

	mds.Version = dec.Uint32()
	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	boneCount := dec.Uint32()
	modelCount := dec.Uint32()

	nameData := dec.Bytes(int(nameLength))

	err := mds.name.parse(nameData)
	if err != nil {
		return fmt.Errorf("nameDataParse: %w", err)
	}

	for i := 0; i < int(materialCount); i++ {
		material := &ModMaterial{}
		material.ID = dec.Int32()
		material.Name = mds.name.byOffset(dec.Int32())
		material.EffectName = mds.name.byOffset(dec.Int32())

		mds.Materials = append(mds.Materials, material)
		propertyCount := dec.Uint32()
		for j := 0; j < int(propertyCount); j++ {
			property := &ModMaterialParam{
				Name: material.Name,
			}

			property.Name = mds.name.byOffset(dec.Int32())

			property.Type = MaterialParamType(dec.Uint32())
			if property.Type == 0 {
				property.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Int32()
				if property.Type == 2 {
					property.Value = mds.name.byOffset(val)
				} else {
					property.Value = fmt.Sprintf("%d", val)
				}
			}

			material.Properties = append(material.Properties, property)
		}
	}

	for i := 0; i < int(boneCount); i++ {
		bone := &ModBone{}
		bone.Name = mds.name.byOffset(dec.Int32())
		bone.Next = dec.Int32()
		bone.ChildrenCount = dec.Uint32()
		bone.ChildIndex = dec.Int32()
		bone.Pivot[0] = dec.Float32()
		bone.Pivot[1] = dec.Float32()
		bone.Pivot[2] = dec.Float32()
		bone.Quaternion[0] = dec.Float32()
		bone.Quaternion[1] = dec.Float32()
		bone.Quaternion[2] = dec.Float32()
		bone.Quaternion[3] = dec.Float32()
		bone.Scale[0] = dec.Float32()
		bone.Scale[1] = dec.Float32()
		bone.Scale[2] = dec.Float32()
		mds.Bones = append(mds.Bones, bone)
	}

	//mds.MainNameIndex = dec.Int32()
	//mds.SubNameIndex = dec.Int32()

	//verticesCount := dec.Uint32()
	//triangleCount := dec.Uint32()

	//boneAssignmentCount := dec.Uint32()

	for i := 0; i < int(modelCount); i++ {
		model := &MdsModel{}
		model.MainPiece = dec.Uint32()
		model.Name = mds.name.byOffset(dec.Int32())
		verticesCount := dec.Uint32()
		faceCount := dec.Uint32()
		boneAssignmentCount := dec.Uint32()
		for i := 0; i < int(verticesCount); i++ {
			v := &ModVertex{}
			v.Position[0] = dec.Float32()
			v.Position[1] = dec.Float32()
			v.Position[2] = dec.Float32()
			v.Normal[0] = dec.Float32()
			v.Normal[1] = dec.Float32()
			v.Normal[2] = dec.Float32()
			if mds.Version <= 2 {
				v.Tint = [4]uint8{128, 128, 128, 255}
			} else {
				v.Tint = [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()}
			}
			v.Uv[0] = dec.Float32()
			v.Uv[1] = dec.Float32()
			if mds.Version <= 2 {
				v.Uv2[0] = 0
				v.Uv2[1] = 0
			} else {
				v.Uv2[0] = dec.Float32()
				v.Uv2[1] = dec.Float32()
			}

			model.Vertices = append(model.Vertices, v)
		}

		for i := 0; i < int(faceCount); i++ {
			f := &ModFace{}
			f.Index[0] = dec.Uint32()
			f.Index[1] = dec.Uint32()
			f.Index[2] = dec.Uint32()

			materialID := dec.Int32()

			var material *ModMaterial
			for _, mat := range mds.Materials {
				if mat.ID == materialID {
					material = mat
					break
				}
			}
			if material == nil {
				if materialID != -1 && materialID != 65536 {
					fmt.Printf("Material %d not found", materialID)
					//return fmt.Errorf("material %d not found", materialID)
				}
				f.MaterialName = ""
			} else {
				f.MaterialName = material.Name
			}

			f.Flags = dec.Uint32()
			model.Faces = append(model.Faces, f)
		}

		if boneAssignmentCount > 99999 {
			return fmt.Errorf("bone assignment count too high: %d", boneAssignmentCount)
		}

		for i := 0; i < int(boneAssignmentCount); i++ {
			_ = dec.Uint32() //weightCount
			weights := [4]*MdsBoneWeight{}
			for j := 0; j < int(4); j++ {
				weight := &MdsBoneWeight{}
				weight.BoneIndex = dec.Int32()
				weight.Value = dec.Float32()
				weights[j] = weight
			}
			model.BoneAssignments = append(model.BoneAssignments, weights)
		}

		mds.Models = append(mds.Models, model)

	}

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	return nil
}

// SetFileName sets the name of the file
func (mds *Mds) SetFileName(name string) {
	mds.MetaFileName = name
}

// FileName returns the name of the file
func (mds *Mds) FileName() string {
	return mds.MetaFileName
}
