package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/pfs/archive"
)

// PolyhedronDef information
type PolyhedronDef struct {
	name string
}

func LoadPolyhedronDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &PolyhedronDef{}
	err := parsePolyhedronDef(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse PolyhedronDef: %w", err)
	}
	return v, nil
}

func parsePolyhedronDef(r io.ReadSeeker, v *PolyhedronDef) error {
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

func (v *PolyhedronDef) FragmentType() string {
	return "PolyhedronDef"
}

func (e *PolyhedronDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
