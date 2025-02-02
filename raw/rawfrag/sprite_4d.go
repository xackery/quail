package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSprite4D is Sprite4D in libeq, empty in openzone, 4DSPRITE (ref) in wld
type WldFragSprite4D struct {
	nameRef  int32
	FourDRef int32
	Params1  uint32
}

func (e *WldFragSprite4D) FragCode() int {
	return FragCodeSprite4D
}

func (e *WldFragSprite4D) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Int32(e.FourDRef)
	enc.Uint32(e.Params1)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite4D) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.FourDRef = dec.Int32()
	e.Params1 = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSprite4D) NameRef() int32 {
	return e.nameRef
}
