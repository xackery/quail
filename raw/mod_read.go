package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type Mod struct {
	MetaFileName string
	Version      uint32
	Materials    []*ModMaterial
	Vertices     []*ModVertex
	Faces        []ModFace
	Bones        []*ModBone
	name         *eqgName
}

// ModVertex is a vertex
type ModVertex struct {
	Position [3]float32
	Normal   [3]float32
	Tint     [4]uint8
	Uv       [2]float32
	Uv2      [2]float32
	Weights  []*ModBoneWeight
}

// ModBone is a bone
type ModBone struct {
	Name          string
	Next          int32
	ChildrenCount uint32
	ChildIndex    int32
	Pivot         [3]float32
	Quaternion    [4]float32
	Scale         [3]float32
}

type ModBoneWeight struct {
	BoneIndex int32
	Value     float32
}

// ModFace is a triangle
type ModFace struct {
	Index        [3]uint32
	MaterialName string
	Flags        uint32
}

type ModMaterial struct {
	ID         int32
	Name       string
	ShaderName string
	Flags      uint32
	Properties []*ModMaterialParam
	Animation  ModMaterialAnimation
}

type MaterialParamType uint32

const (
	MaterialParamTypeUnused MaterialParamType = iota
	MaterialParamTypeInt
	MaterialParamTypeTexture
	MaterialParamTypeColor
)

// ModMaterialParam is a material property
type ModMaterialParam struct {
	Name  string
	Type  MaterialParamType
	Value string
	Data  []byte
}

type ModMaterialAnimation struct {
	Sleep    uint32
	Textures []string
}

type ModFaceFlag uint32

const (
	ModFaceFlagNone              ModFaceFlag = 0x00
	ModFaceFlagPassable          ModFaceFlag = 0x01
	ModFaceFlagTransparent       ModFaceFlag = 0x02
	ModFaceFlagCollisionRequired ModFaceFlag = 0x04
	ModFaceFlagCulled            ModFaceFlag = 0x08
	ModFaceFlagDegenerate        ModFaceFlag = 0x10
)

func (mod *Mod) Identity() string {
	return "mod"
}

func (mod *Mod) String() string {
	out := ""
	out += fmt.Sprintf("metafilename: %s\n", mod.MetaFileName)
	out += fmt.Sprintf("version: %d\n", mod.Version)
	out += fmt.Sprintf("materials: %d\n", len(mod.Materials))
	out += fmt.Sprintf("vertices: %d\n", len(mod.Vertices))
	out += fmt.Sprintf("faces: %d\n", len(mod.Faces))
	out += fmt.Sprintf("bones: %d", len(mod.Bones))
	return out
}

// Decode reads a MOD file
func (mod *Mod) Read(r io.ReadSeeker) error {
	if mod.name == nil {
		mod.name = &eqgName{}
	}
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	mod.Version = dec.Uint32()
	if header != "EQGM" {
		return fmt.Errorf("invalid header %s on version %d, wanted EQGM", header, mod.Version)
	}

	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	verticesCount := dec.Uint32()
	faceCount := dec.Uint32()
	bonesCount := dec.Uint32()
	nameData := dec.Bytes(int(nameLength))
	mod.name.parse(nameData)

	for i := 0; i < int(materialCount); i++ {
		material := &ModMaterial{}
		material.ID = dec.Int32()
		material.Name = mod.name.byOffset(dec.Int32())
		material.ShaderName = mod.name.byOffset(dec.Int32())
		mod.Materials = append(mod.Materials, material)

		paramCount := dec.Uint32()
		for j := 0; j < int(paramCount); j++ {
			param := &ModMaterialParam{
				Name: material.Name,
			}

			param.Name = mod.name.byOffset(dec.Int32())

			param.Type = MaterialParamType(dec.Uint32())
			if param.Type == 0 {
				param.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Int32()
				if param.Type == 2 {
					param.Value = mod.name.byOffset(val)

				} else {
					param.Value = fmt.Sprintf("%d", val)
				}
			}
			material.Properties = append(material.Properties, param)
		}
	}

	for i := 0; i < int(verticesCount); i++ {
		v := &ModVertex{}
		v.Position[0] = dec.Float32()
		v.Position[1] = dec.Float32()
		v.Position[2] = dec.Float32()
		v.Normal[0] = dec.Float32()
		v.Normal[1] = dec.Float32()
		v.Normal[2] = dec.Float32()
		if mod.Version <= 2 {
			v.Tint = [4]uint8{128, 128, 128, 255}
		} else {
			v.Tint = [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()}
		}
		v.Uv[0] = dec.Float32()
		v.Uv[1] = dec.Float32()

		if mod.Version <= 2 {
			v.Uv2[0] = 0
			v.Uv2[1] = 0
		} else {
			v.Uv2[0] = dec.Float32()
			v.Uv2[1] = dec.Float32()
		}

		mod.Vertices = append(mod.Vertices, v)
	}

	for i := 0; i < int(faceCount); i++ {
		f := ModFace{}
		f.Index[0] = dec.Uint32()
		f.Index[1] = dec.Uint32()
		f.Index[2] = dec.Uint32()

		materialID := dec.Int32()

		var material *ModMaterial

		if materialID != -1 && len(mod.Materials) > 0 {
			if len(mod.Materials) < int(materialID) {
				return fmt.Errorf("mod material %d is beyond size of materials (%d)", materialID, len(mod.Materials))
			}
			material = mod.Materials[materialID]
			f.MaterialName = material.Name
		}

		f.Flags = dec.Uint32()
		mod.Faces = append(mod.Faces, f)
	}

	for i := 0; i < int(bonesCount); i++ {
		bone := &ModBone{}
		bone.Name = mod.name.byOffset(dec.Int32())
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

		mod.Bones = append(mod.Bones, bone)
	}

	if bonesCount > 0 {
		for i := 0; i < int(verticesCount); i++ {
			count := dec.Int32()
			mod.Vertices[i].Weights = []*ModBoneWeight{}
			for j := 0; j < int(4); j++ {
				weight := &ModBoneWeight{
					BoneIndex: dec.Int32(),
					Value:     dec.Float32(),
				}
				if j >= int(count) {
					continue
				}
				mod.Vertices[i].Weights = append(mod.Vertices[i].Weights, weight)
			}
		}
	}

	pos := dec.Pos()
	endPos, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}
	if pos != endPos {
		if pos < endPos {
			return fmt.Errorf("%d bytes remaining (%d total)", endPos-pos, endPos)
		}

		return fmt.Errorf("read past end of file")
	}
	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	return nil
}

// SetFileName sets the name of the file
func (mod *Mod) SetFileName(name string) {
	mod.MetaFileName = name
}

// FileName returns the name of the file
func (mod *Mod) FileName() string {
	return mod.MetaFileName
}
