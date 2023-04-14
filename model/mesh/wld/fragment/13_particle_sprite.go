package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// ParticleSprite information
type ParticleSprite struct {
	name      string
	Reference uint32
}

func LoadParticleSprite(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &ParticleSprite{}
	err := parseParticleSprite(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse particlesprite: %w", err)
	}
	return v, nil
}

func parseParticleSprite(r io.ReadSeeker, v *ParticleSprite) error {
	if v == nil {
		return fmt.Errorf("particlesprite is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
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
	return "ParticleSprite"
}

func (e *ParticleSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}