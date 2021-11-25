package fragment

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ParticleSprite information
type ParticleSprite struct {
	hashIndex uint32
	Reference uint32
}

func loadParticleSprite(r io.ReadSeeker) (Fragment, error) {
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
