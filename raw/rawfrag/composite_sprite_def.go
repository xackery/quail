package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragCompositeSpriteDef is empty in libeq, empty in openzone, COMPOSITESPRITEDEF in wld, Actor in lantern
type WldFragCompositeSpriteDef struct {
	NameRef int32  `yaml:"name_ref"`
	Flags   uint32 `yaml:"flags"`
}

func (e *WldFragCompositeSpriteDef) FragCode() int {
	return FragCodeCompositeSpriteDef
}

func (e *WldFragCompositeSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragCompositeSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
