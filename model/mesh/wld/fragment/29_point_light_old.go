package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// PointLightOld information
type PointLightOld struct {
	name      string
	Reference uint32
}

func LoadPointLightOld(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &PointLightOld{}
	err := parsePointLightOld(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse PointLightOld: %w", err)
	}
	return e, nil
}

func parsePointLightOld(r io.ReadSeeker, v *PointLightOld) error {
	if v == nil {
		return fmt.Errorf("PointlightOld is nil")
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

func (e *PointLightOld) FragmentType() string {
	return "PointLightOld"
}

func (e *PointLightOld) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
