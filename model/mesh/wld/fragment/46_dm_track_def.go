package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// DmTrackDef information
type DmTrackDef struct {
	name      string
	Reference uint32
	Name      string
	Position  [3]float32
	Rotation  [3]float32
	Scale     [3]float32
}

func LoadDmTrackDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &DmTrackDef{}
	err := parseDmTrackDef(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse DmTrackDef: %w", err)
	}
	return v, nil
}

func parseDmTrackDef(r io.ReadSeeker, v *DmTrackDef) error {
	if v == nil {
		return fmt.Errorf("DmTrackDef is nil")
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

func (v *DmTrackDef) FragmentType() string {
	return "DmTrackDef"
}

func (e *DmTrackDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
