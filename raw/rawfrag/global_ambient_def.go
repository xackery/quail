package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragGlobalAmbientLightDef is GlobalAmbientLightDef in libeq, WldFragGlobalAmbientLightDef Fragment in openzone, empty in wld, GlobalAmbientLight in lantern
type WldFragGlobalAmbientLightDef struct {
	NameRef int32
}

func (e *WldFragGlobalAmbientLightDef) FragCode() int {
	return FragCodeGlobalAmbientLightDef
}

// Read writes the fragment to the writer
func (e *WldFragGlobalAmbientLightDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragGlobalAmbientLightDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
