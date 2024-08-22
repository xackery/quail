package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model"
)

type Mds struct {
	MetaFileName    string               `yaml:"file_name"`
	Version         uint32               `yaml:"version"`
	Materials       []*Material          `yaml:"materials"`
	Bones           []*Bone              `yaml:"bones"`
	MainNameIndex   int32                `yaml:"main_name_index"`
	SubNameIndex    int32                `yaml:"sub_name_index"`
	Vertices        []*Vertex            `yaml:"vertices"`
	Triangles       []Triangle           `yaml:"triangles"`
	Subs            []*MdsSub            `yaml:"subs"`
	BoneAssignments []*MdsBoneAssignment `yaml:"bone_assignments"`
	names           []*nameEntry
	nameBuf         []byte
}

func (mds *Mds) Identity() string {
	return "mds"
}

type MdsSub struct {
}

type MdsBoneAssignment struct {
}

func (mds *Mds) String() string {
	return "mds"
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
	subCount := dec.Uint32()

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

	//model.Header.Name = lastElement

	for i := 0; i < int(materialCount); i++ {
		material := &Material{}
		material.ID = dec.Int32()
		material.Name = mds.Name(dec.Int32())
		material.ShaderName = mds.Name(dec.Int32())

		mds.Materials = append(mds.Materials, material)
		propertyCount := dec.Uint32()
		for j := 0; j < int(propertyCount); j++ {
			property := &MaterialProperty{
				Name: material.Name,
			}

			property.Name = mds.Name(dec.Int32())

			property.Category = dec.Uint32()
			if property.Category == 0 {
				property.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Int32()
				if property.Category == 2 {
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
		bone.Pivot.X = dec.Float32()
		bone.Pivot.Y = dec.Float32()
		bone.Pivot.Z = dec.Float32()
		bone.Rotation.X = dec.Float32()
		bone.Rotation.Y = dec.Float32()
		bone.Rotation.Z = dec.Float32()
		bone.Rotation.W = dec.Float32()
		bone.Scale.X = dec.Float32()
		bone.Scale.Y = dec.Float32()
		bone.Scale.Z = dec.Float32()
		mds.Bones = append(mds.Bones, bone)
	}

	mds.MainNameIndex = dec.Int32()
	mds.SubNameIndex = dec.Int32()

	verticesCount := dec.Uint32()
	triangleCount := dec.Uint32()

	boneAssignmentCount := dec.Uint32()

	for i := 0; i < int(verticesCount); i++ {
		v := &Vertex{}
		v.Position.X = dec.Float32()
		v.Position.Y = dec.Float32()
		v.Position.Z = dec.Float32()
		v.Normal.X = dec.Float32()
		v.Normal.Y = dec.Float32()
		v.Normal.Z = dec.Float32()
		if mds.Version <= 2 {
			v.Tint = model.RGBA{R: 128, G: 128, B: 128, A: 255}
		} else {
			v.Tint = model.RGBA{R: dec.Uint8(), G: dec.Uint8(), B: dec.Uint8(), A: dec.Uint8()}
		}
		v.Uv.X = dec.Float32()
		v.Uv.Y = dec.Float32()
		if mds.Version <= 2 {
			v.Uv2.X = 0
			v.Uv2.Y = 0
		} else {
			v.Uv2.X = dec.Float32()
			v.Uv2.Y = dec.Float32()
		}

		mds.Vertices = append(mds.Vertices, v)
	}

	for i := 0; i < int(triangleCount); i++ {
		t := Triangle{}
		t.Index.X = dec.Uint32()
		t.Index.Y = dec.Uint32()
		t.Index.Z = dec.Uint32()

		materialID := dec.Int32()

		var material *Material
		for _, mat := range mds.Materials {
			if mat.ID == materialID {
				material = mat
				break
			}
		}
		if material == nil {
			if materialID != -1 {
				fmt.Printf("Material %d not found", materialID)
				//return fmt.Errorf("material %d not found", materialID)
			}
			t.MaterialName = ""
		} else {
			t.MaterialName = material.Name
		}

		t.Flag = dec.Uint32()
		mds.Triangles = append(mds.Triangles, t)
	}

	for i := 0; i < int(subCount); i++ {
		// TODO: sub count
	}

	for i := 0; i < int(boneAssignmentCount); i++ {
		// TODO: bone assignment count
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
	if id := mds.NameOffset(name); id != -1 {
		return -id
	}
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
		if v.name == name {
			return int32(k)
		}
	}
	return -1
}

// NameData is used during writing, dumps the name cache
func (mds *Mds) NameData() []byte {

	return helper.WriteStringHash(string(mds.nameBuf))
}

// NameClear purges names and namebuf, called when encode starts
func (mds *Mds) NameClear() {
	mds.names = nil
	mds.nameBuf = nil
}
