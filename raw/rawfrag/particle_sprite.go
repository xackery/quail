package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragParticleSprite is ParticleSprite in libeq, empty in openzone, PARTICLESPRITE (ref) in wld
type WldFragParticleSprite struct {
	nameRef              int32
	ParticleSpriteDefRef int32
	Flags                uint32
}

func (e *WldFragParticleSprite) FragCode() int {
	return FragCodeParticleSprite
}

func (e *WldFragParticleSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Int32(e.ParticleSpriteDefRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragParticleSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.ParticleSpriteDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragParticleSprite) NameRef() int32 {
	return e.nameRef
}
