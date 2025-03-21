package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragBlitSpriteDef is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSprite in lantern
type WldFragBlitSpriteDef struct {
	nameRef           int32
	Flags             uint32
	SpriteInstanceRef uint32
	RenderMethod      uint32
}

func (e *WldFragBlitSpriteDef) FragCode() int {
	return FragCodeBlitSpriteDef
}

func (e *WldFragBlitSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SpriteInstanceRef)
	enc.Uint32(e.RenderMethod)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragBlitSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.SpriteInstanceRef = dec.Uint32()
	e.RenderMethod = dec.Uint32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragBlitSpriteDef) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragBlitSpriteDef) SetNameRef(id int32) {
	e.nameRef = id
}
