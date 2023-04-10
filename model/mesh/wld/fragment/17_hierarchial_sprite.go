package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// HierarchialSprite information
type HierarchialSprite struct {
	name      string
	Reference uint32
	FrameMs   uint32
}

func LoadHierarchialSprite(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &HierarchialSprite{}
	err := parseHierarchialSprite(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse HierarchialSprite: %w", err)
	}
	return v, nil
}

func parseHierarchialSprite(r io.ReadSeeker, v *HierarchialSprite) error {
	if v == nil {
		return fmt.Errorf("HierarchialSprite is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flag: %w", err)
	}

	//TODO: review
	// Either 4 or 5 - maybe something to look into
	// Bits are set 0, or 2. 0 has the extra field for delay.
	// 2 doesn't have any additional fields.
	if value&1 == 1 {
		err = binary.Read(r, binary.LittleEndian, &v.FrameMs)
		if err != nil {
			return fmt.Errorf("read frame ms: %w", err)
		}
	}

	return nil
}

func (v *HierarchialSprite) FragmentType() string {
	return "HierarchialSprite"
}

func (e *HierarchialSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
