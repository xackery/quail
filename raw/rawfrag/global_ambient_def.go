package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragGlobalAmbientLightDef is GlobalAmbientLightDef in libeq, WldFragGlobalAmbientLightDef Fragment in openzone, empty in wld, GlobalAmbientLight in lantern
type WldFragGlobalAmbientLightDef struct {
	Color [4]uint8
}

func (e *WldFragGlobalAmbientLightDef) FragCode() int {
	return FragCodeGlobalAmbientLightDef
}

// Read writes the fragment to the writer
func (e *WldFragGlobalAmbientLightDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Uint8(e.Color[0])
	enc.Uint8(e.Color[1])
	enc.Uint8(e.Color[2])
	enc.Uint8(e.Color[3])
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragGlobalAmbientLightDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.Color[0] = dec.Uint8()
	e.Color[1] = dec.Uint8()
	e.Color[2] = dec.Uint8()
	e.Color[3] = dec.Uint8()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragGlobalAmbientLightDef) NameRef() int32 {
	return 0
}
