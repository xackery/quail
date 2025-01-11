package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragPolyhedron is Polyhedron in libeq, Polygon Animation Reference in openzone, POLYHEDRON (ref) in wld, Fragment18 in lantern
type WldFragPolyhedron struct {
	nameRef     int32   `yaml:"name_ref"`
	FragmentRef int32   `yaml:"fragment_ref"`
	Flags       uint32  `yaml:"flags"`
	Scale       float32 `yaml:"scale"`
}

func (e *WldFragPolyhedron) FragCode() int {
	return FragCodePolyhedron
}

func (e *WldFragPolyhedron) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Int32(e.FragmentRef)
	enc.Uint32(e.Flags)
	enc.Float32(e.Scale)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPolyhedron) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.FragmentRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.Scale = dec.Float32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragPolyhedron) NameRef() int32 {
	return e.nameRef
}
