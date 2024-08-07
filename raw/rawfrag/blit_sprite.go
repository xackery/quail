package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragBlitSprite is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSpriteReference in lantern
type WldFragBlitSprite struct {
	NameRef       int32
	BlitSpriteRef int32
	Unknown       int32
}

func (e *WldFragBlitSprite) FragCode() int {
	return FragCodeBlitSprite
}

func (e *WldFragBlitSprite) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.BlitSpriteRef)
	enc.Int32(e.Unknown)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragBlitSprite) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.BlitSpriteRef = dec.Int32()
	e.Unknown = dec.Int32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
