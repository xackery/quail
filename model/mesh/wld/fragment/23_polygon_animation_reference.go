package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/pfs/archive"
)

// PolygonAnimationReference information
type PolygonAnimationReference struct {
	name string
}

func LoadPolygonAnimationReference(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &PolygonAnimationReference{}
	err := parsePolygonAnimationReference(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse PolygonAnimationReference: %w", err)
	}
	return v, nil
}

func parsePolygonAnimationReference(r io.ReadSeeker, v *PolygonAnimationReference) error {
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHasIndex: %w", err)
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

func (v *PolygonAnimationReference) FragmentType() string {
	return "PolygonAnimationReference"
}

func (e *PolygonAnimationReference) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
