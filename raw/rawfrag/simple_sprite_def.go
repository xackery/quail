package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSimpleSpriteDef is SimpleSpriteDef in libeq, WldFragSimpleSpriteDef Bitmap Info in openzone, SIMPLESPRITEDEF in wld, BitmapInfo in lantern
type WldFragSimpleSpriteDef struct {
	nameRef      int32
	Flags        uint32
	CurrentFrame int32
	Sleep        uint32
	BitmapRefs   []uint32
}

func (e *WldFragSimpleSpriteDef) FragCode() int {
	return FragCodeSimpleSpriteDef
}

func (e *WldFragSimpleSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.BitmapRefs)))
	if e.Flags&0x04 != 0 {
		enc.Int32(e.CurrentFrame)
	}
	if e.Flags&0x08 != 0 {
		enc.Uint32(e.Sleep)
	}
	for _, textureRef := range e.BitmapRefs {
		enc.Uint32(textureRef)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSimpleSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	textureRefCount := dec.Uint32()
	if e.Flags&0x04 != 0 {
		e.CurrentFrame = dec.Int32()
	}
	if e.Flags&0x08 != 0 {
		e.Sleep = dec.Uint32()
	}
	for i := 0; i < int(textureRefCount); i++ {
		e.BitmapRefs = append(e.BitmapRefs, dec.Uint32())
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSimpleSpriteDef) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragSimpleSpriteDef) SetNameRef(nameRef int32) {
	e.nameRef = nameRef
}
