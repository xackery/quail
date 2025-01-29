package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragCompositeSpriteDef is empty in libeq, empty in openzone, COMPOSITESPRITEDEF in wld, Actor in lantern
type WldFragCompositeSpriteDef struct {
	nameRef int32
	Flags   uint32
}

func (e *WldFragCompositeSpriteDef) FragCode() int {
	return FragCodeCompositeSpriteDef
}

func (e *WldFragCompositeSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragCompositeSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragCompositeSpriteDef) NameRef() int32 {
	return e.nameRef
}
