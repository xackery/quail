package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// ParticleSpriteReference information
type ParticleSpriteReference struct {
	name      string
	Reference uint32
}

func LoadParticleSpriteReference(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &ParticleSpriteReference{}
	err := parseParticleSpriteReference(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse particle sprite reference: %w", err)
	}
	return v, nil
}

func parseParticleSpriteReference(r io.ReadSeeker, v *ParticleSpriteReference) error {
	if v == nil {
		return fmt.Errorf("particle sprite reference is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read post reference: %w", err)
	}
	if value != 8 {
		return fmt.Errorf("post reference got %d, wanted %d", value, 8)
	}

	return nil
}

func (v *ParticleSpriteReference) FragmentType() string {
	return "Particle Sprite Reference"
}
func (e *ParticleSpriteReference) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
