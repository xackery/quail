package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/tag"
)

type Ani struct {
	MetaFileName string     `yaml:"file_name"`
	Version      uint32     `yaml:"version"`
	Bones        []*AniBone `yaml:"bones,omitempty"`
	IsStrict     bool       `yaml:"is_strict,omitempty"`
}

type AniBone struct {
	Name   string          `yaml:"name"`
	Frames []*AniBoneFrame `yaml:"frames,omitempty"`
}

// AniBoneFrame is a bone animation frame
type AniBoneFrame struct {
	Milliseconds uint32  `yaml:"milliseconds"`
	Translation  Vector3 `yaml:"translation"`
	Rotation     Quad4   `yaml:"rotation"`
	Scale        Vector3 `yaml:"scale"`
}

// Read an ANI file
func (ani *Ani) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGA" {
		return fmt.Errorf("invalid header %s, wanted EQGA", header)
	}
	tag.New()

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
	tag.Add(0, dec.Pos(), "red", "header")

	nameData := dec.Bytes(int(nameLength))
	tag.Add(tag.LastPos(), dec.Pos(), "green", "names")

	names := make(map[uint32]string)
	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[uint32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	for i := 0; i < int(boneCount); i++ {
		bone := &AniBone{}
		frameCount := dec.Uint32()

		nameOffset := dec.Uint32()
		bone.Name = names[nameOffset]

		for j := 0; j < int(frameCount); j++ {
			frame := &AniBoneFrame{}

			frame.Milliseconds = dec.Uint32()
			frame.Translation = Vector3{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
			}

			frame.Rotation = Quad4{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
				W: dec.Float32(),
			}

			frame.Scale = Vector3{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
			}
			bone.Frames = append(bone.Frames, frame)
		}

		ani.Bones = append(ani.Bones, bone)
	}

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
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
