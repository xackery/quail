package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// ParticleSpriteDef information
type ParticleSpriteDef struct {
}

func LoadParticleSpriteDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &ParticleSpriteDef{}
	err := parseParticleSpriteDef(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse ParticleSpriteDef: %w", err)
	}
	return e, nil
}

func parseParticleSpriteDef(r io.ReadSeeker, e *ParticleSpriteDef) error {
	if e == nil {
		return fmt.Errorf("ParticleSpriteDef is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *ParticleSpriteDef) FragmentType() string {
	return "ParticleSpriteDef"
}

func (e *ParticleSpriteDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
