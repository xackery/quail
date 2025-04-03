package raw

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
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

type MdsModel struct {
	MainPiece uint32 // 0: no, 1: yes, head is a mainpiece
	Name      string
	Vertices  []*ModVertex
	Faces     []*ModFace
	BoneCount uint32
}

func (mds *Mds) Identity() string {
	return "mds"
}

func (mds *Mds) String() string {
	if mds.name == nil {
		mds.name = &eqgName{}
	}
	out := fmt.Sprintf("Mds: %s,", mds.MetaFileName)
	out += fmt.Sprintf(" %d version,", mds.Version)
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

	out += fmt.Sprintf(", %d bones\n", len(mds.Bones))
	if len(mds.Models) > 0 {

		for i, model := range mds.Models {
			out += fmt.Sprintf("%d model %s (%d verts, %d faces, %d bones)", i, model.Name, len(model.Vertices), len(model.Faces), model.BoneCount)
			if i < len(mds.Models)-1 {
				out += ", "
			}
			out += "\n"
		}

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
		material.ShaderName = mds.name.byOffset(dec.Int32())

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

	for i := 0; i < int(modelCount); i++ {
		model := &MdsModel{}
		model.MainPiece = dec.Uint32()
		model.Name = mds.name.byOffset(dec.Int32())
		verticesCount := dec.Uint32()
		faceCount := dec.Uint32()
		boneCount = dec.Uint32()
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
			if materialID != -1 && len(mds.Materials) > 0 {
				if len(mds.Materials) < int(materialID) {
					return fmt.Errorf("mds material %d is beyond size of materials (%d)", materialID, len(mds.Materials))
				}
				material = mds.Materials[materialID]
				f.MaterialName = material.Name
			}

			f.Flags = dec.Uint32()
			model.Faces = append(model.Faces, f)
		}

		if len(mds.Bones) > 0 {
			for i := 0; i < int(verticesCount); i++ {
				count := dec.Int32()
				model.Vertices[i].Weights = []*ModBoneWeight{}
				for j := 0; j < int(4); j++ {
					weight := &ModBoneWeight{
						BoneIndex: dec.Int32(),
						Value:     dec.Float32(),
					}
					if j < int(count) {
						continue
					}
					model.Vertices[i].Weights = append(model.Vertices[i].Weights, weight)
				}
			}
		}

		mds.Models = append(mds.Models, model)

	}

	pos := dec.Pos()
	endPos, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}
	if pos < endPos {
		remaining := dec.Bytes(int(endPos - pos))
		if !bytes.Equal(remaining, []byte{0x0, 0x0, 0x0, 0x0}) {
			fmt.Printf("remaining bytes: %s\n", hex.Dump(remaining))
			return fmt.Errorf("%d bytes remaining (%d total)", endPos-pos, endPos)
		}
	}
	if pos > endPos {
		return fmt.Errorf("read past end of file")
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
