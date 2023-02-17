package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
)

// PolygonAnimation information
type PolygonAnimation struct {
	name string
}

func LoadPolygonAnimation(r io.ReadSeeker) (common.WldFragmenter, error) {
	v := &PolygonAnimation{}
	err := parsePolygonAnimation(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse PolygonAnimation: %w", err)
	}
	return v, nil
}

func parsePolygonAnimation(r io.ReadSeeker, v *PolygonAnimation) error {
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

func (v *PolygonAnimation) FragmentType() string {
	return "PolygonAnimation"
}

func (e *PolygonAnimation) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
