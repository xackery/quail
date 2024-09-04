package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

type Ter struct {
	MetaFileName string      `yaml:"file_name"`
	Version      uint32      `yaml:"version"`
	Materials    []*Material `yaml:"materials"`
	Vertices     []Vertex    `yaml:"vertices"`
	Triangles    []Triangle  `yaml:"triangles"`
	names        []*nameEntry
	nameBuf      []byte
}

// Identity returns the type of the struct
func (ter *Ter) Identity() string {
	return "ter"
}

// Read reads a TER file
func (ter *Ter) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGT" {
		return fmt.Errorf("invalid header %s, wanted EQGT", header)
	}

	ter.Version = dec.Uint32()

	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	verticesCount := dec.Uint32()
	triangleCount := dec.Uint32()
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

	ter.NameSet(names)

	nameCounter := 0
	for i := 0; i < int(materialCount); i++ {
		material := &Material{}
		material.ID = dec.Int32()
		nameCounter++

		material.Name = ter.Name(dec.Int32())
		material.ShaderName = ter.Name(dec.Int32())

		ter.Materials = append(ter.Materials, material)

		propertyCount := dec.Uint32()
		for j := 0; j < int(propertyCount); j++ {
			property := &MaterialProperty{
				Name: material.Name,
			}

			property.Name = ter.Name(dec.Int32())

			property.Category = dec.Uint32()
			if property.Category == 0 {
				property.Value = fmt.Sprintf("%0.8f", dec.Float32())
			} else {
				val := dec.Int32()
				if property.Category == 2 {
					property.Value = ter.Name(val)
				} else {
					property.Value = fmt.Sprintf("%d", val)
				}
			}

			material.Properties = append(material.Properties, property)
		}
	}

	for i := 0; i < int(verticesCount); i++ {
		v := Vertex{}
		v.Position[0] = dec.Float32()
		v.Position[1] = dec.Float32()
		v.Position[2] = dec.Float32()
		v.Normal[0] = dec.Float32()
		v.Normal[1] = dec.Float32()
		v.Normal[2] = dec.Float32()
		if ter.Version <= 2 {
			v.Tint = [4]uint8{128, 128, 128, 255}
		} else {
			v.Tint = [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()}
		}
		v.Uv[0] = dec.Float32()
		v.Uv[1] = dec.Float32()
		if ter.Version <= 2 {
			v.Uv2[0] = 0
			v.Uv2[1] = 0
		} else {
			v.Uv2[0] = dec.Float32()
			v.Uv2[1] = dec.Float32()
		}

		ter.Vertices = append(ter.Vertices, v)
	}

	for i := 0; i < int(triangleCount); i++ {
		t := Triangle{}
		t.Index[0] = dec.Uint32()
		t.Index[1] = dec.Uint32()
		t.Index[2] = dec.Uint32()

		materialID := dec.Int32()

		var material *Material
		for _, mat := range ter.Materials {
			if mat.ID == materialID {
				material = mat
				break
			}
		}
		if material == nil {
			//if materialID != -1 {
			//log.Warnf("material %d not found", materialID)
			//return fmt.Errorf("material %d not found", materialID)
			//}
			t.MaterialName = ""
		} else {
			t.MaterialName = material.Name
		}

		t.Flag = dec.Uint32()
		ter.Triangles = append(ter.Triangles, t)
	}

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	return nil
}

// SetFileName sets the name of the file
func (ter *Ter) SetFileName(name string) {
	ter.MetaFileName = name
}

// FileName returns the name of the file
func (ter *Ter) FileName() string {
	return ter.MetaFileName
}

// Name is used during reading, returns the Name of an id
func (ter *Ter) Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if ter.names == nil {
		return fmt.Sprintf("!UNK(%d)", id)
	}
	//fmt.Println("name: [", names[id], "]")

	for _, v := range ter.names {
		if int32(v.offset) == id {
			return v.name
		}
	}
	return fmt.Sprintf("!UNK(%d)", id)
}

// NameSet is used during reading, sets the names within a buffer
func (ter *Ter) NameSet(newNames map[int32]string) {
	if newNames == nil {
		ter.names = []*nameEntry{}
		return
	}
	for k, v := range newNames {
		ter.names = append(ter.names, &nameEntry{offset: int(k), name: v})
	}
	ter.nameBuf = []byte{0x00}

	for _, v := range ter.names {
		ter.nameBuf = append(ter.nameBuf, []byte(v.name)...)
		ter.nameBuf = append(ter.nameBuf, 0)
	}
}

// NameAdd is used when writing, appending new names
func (ter *Ter) NameAdd(name string) int32 {

	if ter.names == nil {
		ter.names = []*nameEntry{
			{offset: 0, name: ""},
		}
		ter.nameBuf = []byte{0x00}
	}
	if name == "" {
		return 0
	}

	/* if name[len(ter.name)-1:] != "\x00" {
		name += "\x00"
	}
	*/
	if id := ter.NameOffset(name); id != -1 {
		return -id
	}
	ter.names = append(ter.names, &nameEntry{offset: len(ter.nameBuf), name: name})
	lastRef := int32(len(ter.nameBuf))
	ter.nameBuf = append(ter.nameBuf, []byte(name)...)
	ter.nameBuf = append(ter.nameBuf, 0)
	return int32(-lastRef)
}

func (ter *Ter) NameOffset(name string) int32 {
	if ter.names == nil {
		return -1
	}
	for _, v := range ter.names {
		if v.name == name {
			return int32(v.offset)
		}
	}
	return -1
}

// NameIndex is used when reading, returns the index of a name, or -1 if not found
func (ter *Ter) NameIndex(name string) int32 {
	if ter.names == nil {
		return -1
	}
	for k, v := range ter.names {
		if v.name == name {
			return int32(k)
		}
	}
	return -1
}

// NameData is used during writing, dumps the name cache
func (ter *Ter) NameData() []byte {

	return helper.WriteStringHash(string(ter.nameBuf))
}

// NameClear purges names and namebuf, called when encode starts
func (ter *Ter) NameClear() {
	ter.names = nil
	ter.nameBuf = nil
}
