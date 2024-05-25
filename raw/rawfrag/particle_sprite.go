package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragParticleSprite is ParticleSprite in libeq, empty in openzone, PARTICLESPRITE (ref) in wld
type WldFragParticleSprite struct {
	NameRef              int32  `yaml:"name_ref"`
	ParticleSpriteDefRef int32  `yaml:"particle_sprite_def_ref"`
	Flags                uint32 `yaml:"flags"`
}

func (e *WldFragParticleSprite) FragCode() int {
	return FragCodeParticleSprite
}

func (e *WldFragParticleSprite) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.ParticleSpriteDefRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragParticleSprite) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.ParticleSpriteDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
