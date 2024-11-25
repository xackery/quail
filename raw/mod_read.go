package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

type Mod struct {
	MetaFileName string
	Version      uint32
	Materials    []*Material
	Bones        []*Bone
	Vertices     []*Vertex
	Triangles    []Triangle
	names        []*nameEntry
	nameBuf      []byte
}

func (mod *Mod) Identity() string {
	return "mod"
}

// Decode reads a MOD file
func (mod *Mod) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)

	mod.Version = dec.Uint32()
	if header != "EQGM" {
		return fmt.Errorf("invalid header %s on version %d, wanted EQGM", header, mod.Version)
	}

	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	verticesCount := dec.Uint32()
	triangleCount := dec.Uint32()
	bonesCount := dec.Uint32()
	nameData := dec.Bytes(int(nameLength))

	names := make(map[int32]string)
	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[int32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	mod.NameSet(names)

	for i := 0; i < int(materialCount); i++ {
		material := &Material{}
		material.ID = dec.Int32()
		material.Name = mod.Name(dec.Int32())
		material.ShaderName = mod.Name(dec.Int32())
		mod.Materials = append(mod.Materials, material)

		propertyCount := dec.Uint32()
		for j := 0; j < int(propertyCount); j++ {
			property := &MaterialProperty{
				Name: material.Name,
			}

			property.Name = mod.Name(dec.Int32())

			property.Category = dec.Uint32()
			if property.Category == 0 {
				property.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Int32()
				if property.Category == 2 {
					property.Value = mod.Name(val)

				} else {
					property.Value = fmt.Sprintf("%d", val)
				}
			}
			material.Properties = append(material.Properties, property)
		}
	}

	for i := 0; i < int(verticesCount); i++ {
		v := &Vertex{}
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

	for i := 0; i < int(triangleCount); i++ {
		t := Triangle{}
		t.Index[0] = dec.Uint32()
		t.Index[1] = dec.Uint32()
		t.Index[2] = dec.Uint32()

		materialID := dec.Int32()

		var material *Material
		for _, mat := range mod.Materials {
			if mat.ID == materialID {
				material = mat
				break
			}
		}
		if material == nil {
			if materialID != -1 {
				fmt.Printf("Material mod %d not found", materialID)
				//return fmt.Errorf("material %d not found", materialID)
			}
			t.MaterialName = ""
		} else {
			t.MaterialName = material.Name
		}

		t.Flag = dec.Uint32()
		mod.Triangles = append(mod.Triangles, t)
	}

	for i := 0; i < int(bonesCount); i++ {
		bone := &Bone{}
		bone.Name = mod.Name(dec.Int32())
		bone.Next = dec.Int32()
		bone.ChildrenCount = dec.Uint32()
		bone.ChildIndex = dec.Int32()
		bone.Pivot[0] = dec.Float32()
		bone.Pivot[1] = dec.Float32()
		bone.Pivot[2] = dec.Float32()
		bone.Rotation[0] = dec.Float32()
		bone.Rotation[1] = dec.Float32()
		bone.Rotation[2] = dec.Float32()
		bone.Scale[0] = dec.Float32()
		bone.Scale[1] = dec.Float32()
		bone.Scale[2] = dec.Float32()
		bone.Scale2 = dec.Float32()

		mod.Bones = append(mod.Bones, bone)
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

// Name is used during reading, returns the Name of an id
func (mod *Mod) Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if mod.names == nil {
		return fmt.Sprintf("!UNK(%d)", id)
	}
	//fmt.Println("name: [", names[id], "]")

	for _, v := range mod.names {
		if int32(v.offset) == id {
			return v.name
		}
	}
	return fmt.Sprintf("!UNK(%d)", id)
}

// NameSet is used during reading, sets the names within a buffer
func (mod *Mod) NameSet(newNames map[int32]string) {
	if newNames == nil {
		mod.names = []*nameEntry{}
		return
	}
	for k, v := range newNames {
		mod.names = append(mod.names, &nameEntry{offset: int(k), name: v})
	}
	mod.nameBuf = []byte{0x00}

	for _, v := range mod.names {
		mod.nameBuf = append(mod.nameBuf, []byte(v.name)...)
		mod.nameBuf = append(mod.nameBuf, 0)
	}
}

// NameAdd is used when writing, appending new names
func (mod *Mod) NameAdd(name string) int32 {

	if mod.names == nil {
		mod.names = []*nameEntry{
			{offset: 0, name: ""},
		}
		mod.nameBuf = []byte{0x00}
	}
	if name == "" {
		return 0
	}

	/* if name[len(mod.name)-1:] != "\x00" {
		name += "\x00"
	}
	*/
	if id := mod.NameOffset(name); id != -1 {
		return -id
	}
	mod.names = append(mod.names, &nameEntry{offset: len(mod.nameBuf), name: name})
	lastRef := int32(len(mod.nameBuf))
	mod.nameBuf = append(mod.nameBuf, []byte(name)...)
	mod.nameBuf = append(mod.nameBuf, 0)
	return int32(-lastRef)
}

func (mod *Mod) NameOffset(name string) int32 {
	if mod.names == nil {
		return -1
	}
	for _, v := range mod.names {
		if v.name == name {
			return int32(v.offset)
		}
	}
	return -1
}

// NameIndex is used when reading, returns the index of a name, or -1 if not found
func (mod *Mod) NameIndex(name string) int32 {
	if mod.names == nil {
		return -1
	}
	for k, v := range mod.names {
		if v.name == name {
			return int32(k)
		}
	}
	return -1
}

// NameData is used during writing, dumps the name cache
func (mod *Mod) NameData() []byte {

	return helper.WriteStringHash(string(mod.nameBuf))
}

// NameClear purges names and namebuf, called when encode starts
func (mod *Mod) NameClear() {
	mod.names = nil
	mod.nameBuf = nil
}

func (mod *Mod) Names() []string {
	if mod.names == nil {
		return nil
	}
	names := []string{}
	for _, v := range mod.names {
		names = append(names, v.name)
	}
	return names
}
