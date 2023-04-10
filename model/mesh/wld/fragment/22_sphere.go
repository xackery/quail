package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// Sphere information
type Sphere struct {
	name      string
	Reference uint32
}

func LoadSphere(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &Sphere{}
	err := parseSphere(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse Sphere: %w", err)
	}
	return v, nil
}

func parseSphere(r io.ReadSeeker, v *Sphere) error {
	if v == nil {
		return fmt.Errorf("Sphere is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value4: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value12: %w", err)
	}

	return nil
}

func (v *Sphere) FragmentType() string {
	return "Sphere"
}

func (e *Sphere) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
