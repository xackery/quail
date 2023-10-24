package common

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/xackery/encdec"
)

// Animation is an animation
type Animation struct {
	Header   *Header          `yaml:"header,omitempty"`
	Bones    []*BoneAnimation `yaml:"bones,omitempty"`
	IsStrict bool             `yaml:"is_strict,omitempty"`
}

func NewAnimation(name string) *Animation {
	return &Animation{
		Header: &Header{
			Name: name,
		},
	}
}

// NameBuild prepares an EQG-styled name buffer list
func (anim *Animation) NameBuild(miscNames []string) (map[string]int32, []byte, error) {
	var err error

	names := make(map[string]int32)
	nameBuf := bytes.NewBuffer(nil)
	tmpNames := []string{}
	// append materials to tmpNames
	for _, o := range anim.Bones {
		tmpNames = append(tmpNames, o.Name)
	}

	for _, name := range miscNames {
		isNew := true
		for key := range names {
			if key == name {
				isNew = false
				break
			}
		}
		if !isNew {
			continue
		}

		tmpNames = append(tmpNames, name)
	}

	for _, name := range tmpNames {
		isNew := true
		for key := range names {
			if key == name {
				isNew = false
				break
			}
		}
		if !isNew {
			continue
		}

		names[name] = int32(nameBuf.Len())

		_, err = nameBuf.Write([]byte(name))
		if err != nil {
			return nil, nil, fmt.Errorf("write name: %w", err)
		}
		_, err = nameBuf.Write([]byte{0})
		if err != nil {
			return nil, nil, fmt.Errorf("write 0: %w", err)
		}
	}

	return names, nameBuf.Bytes(), nil
}

// BoneBuild prepares an EQG-styled bone buffer list
func (anim *Animation) BoneBuild(version uint32, isMod bool, names map[string]int32) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)

	// bones
	for _, o := range anim.Bones {
		nameOffset := int32(-1)
		for key, val := range names {
			if key == o.Name {
				nameOffset = val
				break
			}
		}
		if nameOffset == -1 {
			return nil, fmt.Errorf("bone %s not found", o.Name)
		}

		enc.Uint32(o.FrameCount)
		enc.Int32(nameOffset)
		for _, frame := range o.Frames {
			enc.Uint32(frame.Milliseconds)
			enc.Float32(frame.Translation.X)
			enc.Float32(frame.Translation.Y)
			enc.Float32(frame.Translation.Z)

			enc.Float32(frame.Rotation.X)
			enc.Float32(frame.Rotation.Y)
			enc.Float32(frame.Rotation.Z)
			enc.Float32(frame.Rotation.W)

			enc.Float32(frame.Scale.X)
			enc.Float32(frame.Scale.Y)
			enc.Float32(frame.Scale.Z)
		}
	}
	return dataBuf.Bytes(), nil
}
