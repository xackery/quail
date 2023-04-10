package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// DmTrack information
type DmTrack struct {
	name      string
	Reference uint32
	Name      string
	Position  [3]float32
	Rotation  [3]float32
	Scale     [3]float32
}

func LoadDmTrack(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &DmTrack{}
	err := parseDmTrack(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse DmTrack: %w", err)
	}
	return v, nil
}

func parseDmTrack(r io.ReadSeeker, v *DmTrack) error {
	if v == nil {
		return fmt.Errorf("DmTrack is nil")
	}

	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	return nil
}

func (v *DmTrack) FragmentType() string {
	return "DmTrack"
}

func (e *DmTrack) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
