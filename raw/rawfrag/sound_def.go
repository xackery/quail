package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSoundDef is empty in libeq, empty in openzone, SOUNDINSTANCE in wld
type WldFragSoundDef struct {
	nameRef int32  `yaml:"name_ref"`
	Flags   uint32 `yaml:"flags"`
}

func (e *WldFragSoundDef) FragCode() int {
	return FragCodeSoundDef
}

func (e *WldFragSoundDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSoundDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSoundDef) NameRef() int32 {
	return e.nameRef
}
