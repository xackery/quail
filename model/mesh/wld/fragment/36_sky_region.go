package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// SkyRegion information
type SkyRegion struct {
	name      string
	Reference uint32
}

func LoadSkyRegion(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &SkyRegion{}
	err := parseSkyRegion(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse SkyRegion: %w", err)
	}
	return e, nil
}

func parseSkyRegion(r io.ReadSeeker, v *SkyRegion) error {
	if v == nil {
		return fmt.Errorf("SkyRegion is nil")
	}
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	return nil
}

func (e *SkyRegion) FragmentType() string {
	return "SkyRegion"
}

func (e *SkyRegion) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
