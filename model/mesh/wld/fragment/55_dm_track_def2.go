package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// DmTrackDef2 information
type DmTrackDef2 struct {
}

func LoadDmTrackDef2(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &DmTrackDef2{}
	err := parseDmTrackDef2(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse DmTrackDef2: %w", err)
	}
	return e, nil
}

func parseDmTrackDef2(r io.ReadSeeker, e *DmTrackDef2) error {
	if e == nil {
		return fmt.Errorf("DmTrackDef2 is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *DmTrackDef2) FragmentType() string {
	return "DmTrackDef2"
}

func (e *DmTrackDef2) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
