package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragCompositeSprite is empty in libeq, empty in openzone, COMPOSITESPRITE (ref) in wld
type WldFragCompositeSprite struct {
	nameRef               int32
	CompositeSpriteDefRef int32
	Flags                 uint32
}

func (e *WldFragCompositeSprite) FragCode() int {
	return FragCodeCompositeSprite
}

func (e *WldFragCompositeSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Int32(e.CompositeSpriteDefRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragCompositeSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.CompositeSpriteDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragCompositeSprite) NameRef() int32 {
	return e.nameRef
}
