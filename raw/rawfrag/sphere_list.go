package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSphereList is SphereList in libeq, empty in openzone, SPHERELIST (ref) in wld
type WldFragSphereList struct {
	NameRef          int32  `yaml:"name_ref"`
	SphereListDefRef int32  `yaml:"sphere_list_def_ref"`
	Params1          uint32 `yaml:"params1"`
}

func (e *WldFragSphereList) FragCode() int {
	return FragCodeSphereList
}

func (e *WldFragSphereList) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.SphereListDefRef)
	enc.Uint32(e.Params1)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSphereList) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.SphereListDefRef = dec.Int32()
	e.Params1 = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
