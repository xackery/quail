package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragPointLightOldDef is empty in libeq, empty in openzone, empty in wld
type WldFragPointLightOldDef struct {
	nameRef       int32 `yaml:"name_ref"`
	PointLightRef int32 `yaml:"point_light_ref"`
}

func (e *WldFragPointLightOldDef) FragCode() int {
	return FragCodePointLightOldDef
}

func (e *WldFragPointLightOldDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Int32(e.PointLightRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPointLightOldDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.PointLightRef = dec.Int32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragPointLightOldDef) NameRef() int32 {
	return e.nameRef
}
