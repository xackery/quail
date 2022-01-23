package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// ParticleSprite information
type ParticleSprite struct {
	hashIndex uint32
	Reference uint32
}

func LoadParticleSprite(r io.ReadSeeker) (common.WldFragmenter, error) {
	v := &ParticleSprite{}
	err := parseParticleSprite(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse particle sprite: %w", err)
	}
	return v, nil
}

func parseParticleSprite(r io.ReadSeeker, v *ParticleSprite) error {
	if v == nil {
		return fmt.Errorf("particle sprite is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &v.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value4: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value12: %w", err)
	}

	return nil
}

func (v *ParticleSprite) FragmentType() string {
	return "Particle Sprite"
}

func (e *ParticleSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
