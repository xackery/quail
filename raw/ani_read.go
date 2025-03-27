package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type Ani struct {
	MetaFileName string
	Version      uint32
	Bones        []*AniBone
	IsStrict     bool
	name         *eqgName
}

func (ani *Ani) Identity() string {
	return "ani"
}

type AniBone struct {
	Name   string
	Frames []*AniBoneFrame
}

// AniBoneFrame is a bone animation frame
type AniBoneFrame struct {
	Milliseconds uint32
	Translation  [3]float32
	Rotation     [4]float32
	Scale        [3]float32
}

func (ani *Ani) String() string {
	out := ""
	out += fmt.Sprintf("metafilename %s\n", ani.MetaFileName)
	out += fmt.Sprintf("version %d\n", ani.Version)
	out += fmt.Sprintf("strict %t\n", ani.IsStrict)
	out += fmt.Sprintf("bones %d\n", len(ani.Bones))
	if len(ani.Bones) == 0 {
		return out
	}
	out += fmt.Sprintf("frames per bone %d", len(ani.Bones[0].Frames))

	return out

}

// Read an ANI file
func (ani *Ani) Read(r io.ReadSeeker) error {
	if ani.name == nil {
		ani.name = &eqgName{}
	}
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
	ani.name.parse(nameData)

	for i := 0; i < int(boneCount); i++ {
		bone := &AniBone{}
		frameCount := dec.Uint32()

		bone.Name = ani.name.byOffset(dec.Int32())

		for j := 0; j < int(frameCount); j++ {
			frame := &AniBoneFrame{}

			frame.Milliseconds = dec.Uint32()
			frame.Translation = [3]float32{
				dec.Float32(),
				dec.Float32(),
				dec.Float32(),
			}

			frame.Rotation = [4]float32{
				dec.Float32(),
				dec.Float32(),
				dec.Float32(),
				dec.Float32(),
			}

			frame.Scale = [3]float32{
				dec.Float32(),
				dec.Float32(),
				dec.Float32(),
			}
			bone.Frames = append(bone.Frames, frame)
		}

		ani.Bones = append(ani.Bones, bone)
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

	err = dec.Error()
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
