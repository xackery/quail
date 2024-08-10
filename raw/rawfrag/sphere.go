package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSphere is Sphere in libeq, Zone Unknown in openzone, SPHERE (ref) in wld, Fragment16 in lantern
type WldFragSphere struct {
	NameRef int32   `yaml:"name_ref"`
	Radius  float32 `yaml:"radius"`
}

func (e *WldFragSphere) FragCode() int {
	return FragCodeSphere
}

func (e *WldFragSphere) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Float32(e.Radius)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSphere) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Radius = dec.Float32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
