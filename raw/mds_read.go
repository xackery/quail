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
	Materials    []*Material
	Bones        []*Bone
	//MainNameIndex int32
	//SubNameIndex  int32
	Models  []*MdsModel
	names   []*nameEntry
	nameBuf []byte
}

func (mds *Mds) Identity() string {
	return "mds"
}

type MdsModel struct {
	MainPiece       uint32 // 0: no, 1: yes, head is a mainpiece
	Name            string
	Vertices        []*Vertex
	Faces           []*Face
	BoneAssignments [][4]*MdsBoneWeight
}

type MdsBoneWeight struct {
	BoneIndex int32
	Value     float32
}

func (mds *Mds) String() string {
	out := fmt.Sprintf("Mds: %s,", mds.MetaFileName)
	out += fmt.Sprintf(" %d names,", len(mds.names))
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

	names := make(map[int32]string)
	chunk := []byte{}
	lastOffset := 0
	//lastElement := ""
	for i, b := range nameData {
		if b == 0 {
			names[int32(lastOffset)] = string(chunk)
			//	lastElement = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	mds.NameSet(names)

	for i := 0; i < int(materialCount); i++ {
		material := &Material{}
		material.ID = dec.Int32()
		material.Name = mds.Name(dec.Int32())
		material.EffectName = mds.Name(dec.Int32())

		mds.Materials = append(mds.Materials, material)
		propertyCount := dec.Uint32()
		for j := 0; j < int(propertyCount); j++ {
			property := &MaterialParam{
				Name: material.Name,
			}

			property.Name = mds.Name(dec.Int32())

			property.Type = MaterialParamType(dec.Uint32())
			if property.Type == 0 {
				property.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Int32()
				if property.Type == 2 {
					property.Value = mds.Name(val)
				} else {
					property.Value = fmt.Sprintf("%d", val)
				}
			}

			material.Properties = append(material.Properties, property)
		}
	}

	for i := 0; i < int(boneCount); i++ {
		bone := &Bone{}
		bone.Name = mds.Name(dec.Int32())
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
		model.Name = mds.Name(dec.Int32())
		verticesCount := dec.Uint32()
		faceCount := dec.Uint32()
		boneAssignmentCount := dec.Uint32()
		for i := 0; i < int(verticesCount); i++ {
			v := &Vertex{}
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
			f := &Face{}
			f.Index[0] = dec.Uint32()
			f.Index[1] = dec.Uint32()
			f.Index[2] = dec.Uint32()

			materialID := dec.Int32()

			var material *Material
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

// Name is used during reading, returns the Name of an id
func (mds *Mds) Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if mds.names == nil {
		return fmt.Sprintf("!UNK(%d)", id)
	}
	//fmt.Println("name: [", names[id], "]")

	for _, v := range mds.names {
		if int32(v.offset) == id {
			return v.name
		}
	}
	return fmt.Sprintf("!UNK(%d)", id)
}

// NameSet is used during reading, sets the names within a buffer
func (mds *Mds) NameSet(newNames map[int32]string) {
	if newNames == nil {
		mds.names = []*nameEntry{}
		return
	}
	for k, v := range newNames {
		mds.names = append(mds.names, &nameEntry{offset: int(k), name: v})
	}
	mds.nameBuf = []byte{0x00}

	for _, v := range mds.names {
		mds.nameBuf = append(mds.nameBuf, []byte(v.name)...)
		mds.nameBuf = append(mds.nameBuf, 0)
	}
}

// NameAdd is used when writing, appending new names
func (mds *Mds) NameAdd(name string) int32 {

	if mds.names == nil {
		mds.names = []*nameEntry{
			{offset: 0, name: ""},
		}
		mds.nameBuf = []byte{0x00}
	}
	if name == "" {
		return 0
	}

	/* if name[len(mds.name)-1:] != "\x00" {
		name += "\x00"
	}
	*/
	// if id := mds.NameOffset(name); id != -1 {
	// 	return -id
	// }
	mds.names = append(mds.names, &nameEntry{offset: len(mds.nameBuf), name: name})
	lastRef := int32(len(mds.nameBuf))
	mds.nameBuf = append(mds.nameBuf, []byte(name)...)
	mds.nameBuf = append(mds.nameBuf, 0)
	return int32(-lastRef)
}

func (mds *Mds) NameOffset(name string) int32 {
	if mds.names == nil {
		return -1
	}
	for _, v := range mds.names {
		if v.name == name {
			return int32(v.offset)
		}
	}
	return -1
}

// NameIndex is used when reading, returns the index of a name, or -1 if not found
func (mds *Mds) NameIndex(name string) int32 {
	if mds.names == nil {
		return -1
	}
	for k, v := range mds.names {
		if v.name != name {
			continue
		}
		return int32(k)
	}
	return -1
}

// NameData is used during writing, dumps the name cache
func (mds *Mds) NameData() []byte {
	if len(mds.nameBuf) == 0 {
		return nil
	}
	return mds.nameBuf[:len(mds.nameBuf)-1]

	//return helper.WriteStringHash(string(mds.nameBuf))
}

// NameClear purges names and namebuf, called when encode starts
func (mds *Mds) NameClear() {
	mds.names = nil
	mds.nameBuf = nil
}

func (mds *Mds) Names() []string {
	if mds.names == nil {
		return nil
	}
	names := []string{}
	for _, v := range mds.names {
		names = append(names, v.name)
	}
	return names
}
