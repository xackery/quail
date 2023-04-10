package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// ActiveGeoRegion information
type ActiveGeoRegion struct {
	name      string
	Reference uint32
}

func LoadActiveGeoRegion(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &ActiveGeoRegion{}
	err := parseActiveGeoRegion(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse ActiveGeoRegion: %w", err)
	}
	return e, nil
}

func parseActiveGeoRegion(r io.ReadSeeker, v *ActiveGeoRegion) error {
	if v == nil {
		return fmt.Errorf("ActiveGeoRegion is nil")
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

func (e *ActiveGeoRegion) FragmentType() string {
	return "ActiveGeoRegion"
}

func (e *ActiveGeoRegion) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
