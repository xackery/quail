package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/pfs/archive"
)

// Polyhedron information
type Polyhedron struct {
	name string
}

func LoadPolyhedron(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &Polyhedron{}
	err := parsePolyhedron(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse Polyhedron: %w", err)
	}
	return v, nil
}

func parsePolyhedron(r io.ReadSeeker, v *Polyhedron) error {
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	// Always 2 when used in main zone, and object files.
	// This means, it has a bounding radius
	// Some differences in character + model archives
	// Confirmed

	flags := int32(0)
	err = binary.Read(r, binary.LittleEndian, &flags)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}
	dump.Hex(flags, "flags=%d", flags)

	boneCount := int32(0)
	err = binary.Read(r, binary.LittleEndian, &boneCount)
	if err != nil {
		return fmt.Errorf("read boneCount: %w", err)
	}
	dump.Hex(flags, "boneCount=%d", flags)

	return nil
}

func (v *Polyhedron) FragmentType() string {
	return "Polyhedron"
}

func (e *Polyhedron) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
