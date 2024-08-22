package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model"
)

type Ani struct {
	MetaFileName string     `yaml:"file_name"`
	Version      uint32     `yaml:"version"`
	Bones        []*AniBone `yaml:"bones,omitempty"`
	IsStrict     bool       `yaml:"is_strict,omitempty"`
	names        []*nameEntry
	nameBuf      []byte
}

func (ani *Ani) Identity() string {
	return "ani"
}

type AniBone struct {
	Name   string          `yaml:"name"`
	Frames []*AniBoneFrame `yaml:"frames,omitempty"`
}

// AniBoneFrame is a bone animation frame
type AniBoneFrame struct {
	Milliseconds uint32        `yaml:"milliseconds"`
	Translation  model.Vector3 `yaml:"translation"`
	Rotation     model.Quad4   `yaml:"rotation"`
	Scale        model.Vector3 `yaml:"scale"`
}

// Read an ANI file
func (ani *Ani) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGA" {
		return fmt.Errorf("invalid header %s, wanted EQGA", header)
	}

	ani.Version = dec.Uint32()
	nameLength := int(dec.Uint32())
	boneCount := dec.Uint32()

	ani.IsStrict = false
	if ani.Version > 1 {
		isStrict := dec.Uint32()
		if isStrict > 0 {
			ani.IsStrict = true
		}
	}

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

	ani.NameSet(names)

	for i := 0; i < int(boneCount); i++ {
		bone := &AniBone{}
		frameCount := dec.Uint32()

		bone.Name = ani.Name(dec.Int32())

		for j := 0; j < int(frameCount); j++ {
			frame := &AniBoneFrame{}

			frame.Milliseconds = dec.Uint32()
			frame.Translation = model.Vector3{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
			}

			frame.Rotation = model.Quad4{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
				W: dec.Float32(),
			}

			frame.Scale = model.Vector3{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
			}
			bone.Frames = append(bone.Frames, frame)
		}

		ani.Bones = append(ani.Bones, bone)
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// SetFileName sets the name of the file
func (ani *Ani) SetFileName(name string) {
	ani.MetaFileName = name
}

func (ani *Ani) FileName() string {
	return ani.MetaFileName
}

// Name is used during reading, returns the Name of an id
func (ani *Ani) Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if ani.names == nil {
		return fmt.Sprintf("!UNK(%d)", id)
	}
	//fmt.Println("name: [", names[id], "]")

	for _, v := range ani.names {
		if int32(v.offset) == id {
			return v.name
		}
	}
	return fmt.Sprintf("!UNK(%d)", id)
}

// NameSet is used during reading, sets the names within a buffer
func (ani *Ani) NameSet(newNames map[int32]string) {
	if newNames == nil {
		ani.names = []*nameEntry{}
		return
	}
	for k, v := range newNames {
		ani.names = append(ani.names, &nameEntry{offset: int(k), name: v})
	}
	ani.nameBuf = []byte{0x00}

	for _, v := range ani.names {
		ani.nameBuf = append(ani.nameBuf, []byte(v.name)...)
		ani.nameBuf = append(ani.nameBuf, 0)
	}
}

// NameAdd is used when writing, appending new names
func (ani *Ani) NameAdd(name string) int32 {

	if ani.names == nil {
		ani.names = []*nameEntry{
			{offset: 0, name: ""},
		}
		ani.nameBuf = []byte{0x00}
	}
	if name == "" {
		return 0
	}

	/* if name[len(ani.name)-1:] != "\x00" {
		name += "\x00"
	}
	*/
	if id := ani.NameOffset(name); id != -1 {
		return -id
	}
	ani.names = append(ani.names, &nameEntry{offset: len(ani.nameBuf), name: name})
	lastRef := int32(len(ani.nameBuf))
	ani.nameBuf = append(ani.nameBuf, []byte(name)...)
	ani.nameBuf = append(ani.nameBuf, 0)
	return int32(-lastRef)
}

func (ani *Ani) NameOffset(name string) int32 {
	if ani.names == nil {
		return -1
	}
	for _, v := range ani.names {
		if v.name == name {
			return int32(v.offset)
		}
	}
	return -1
}

// NameIndex is used when reading, returns the index of a name, or -1 if not found
func (ani *Ani) NameIndex(name string) int32 {
	if ani.names == nil {
		return -1
	}
	for k, v := range ani.names {
		if v.name == name {
			return int32(k)
		}
	}
	return -1
}

// NameData is used during writing, dumps the name cache
func (ani *Ani) NameData() []byte {

	return helper.WriteStringHash(string(ani.nameBuf))
}

// NameClear purges names and namebuf, called when encode starts
func (ani *Ani) NameClear() {
	ani.names = nil
	ani.nameBuf = nil
}
