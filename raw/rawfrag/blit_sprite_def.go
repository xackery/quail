package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragBlitSpriteDef is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSprite in lantern
type WldFragBlitSpriteDef struct {
	NameRef           int32  `yaml:"name_ref"`
	Flags             uint32 `yaml:"flags"`
	SpriteInstanceRef uint32 `yaml:"sprite_instance_ref"`
	Unknown           int32  `yaml:"unknown"`
}

func (e *WldFragBlitSpriteDef) FragCode() int {
	return FragCodeBlitSpriteDef
}

func (e *WldFragBlitSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SpriteInstanceRef)
	enc.Int32(e.Unknown)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragBlitSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.SpriteInstanceRef = dec.Uint32()
	e.Unknown = dec.Int32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
