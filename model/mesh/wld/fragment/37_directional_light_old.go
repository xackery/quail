package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// DirectionalLightOld information
type DirectionalLightOld struct {
	name      string
	Reference uint32
}

func LoadDirectionalLightOld(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &DirectionalLightOld{}
	err := parseDirectionalLightOld(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse DirectionalLightOld: %w", err)
	}
	return e, nil
}

func parseDirectionalLightOld(r io.ReadSeeker, v *DirectionalLightOld) error {
	if v == nil {
		return fmt.Errorf("DirectionalLightOld is nil")
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

func (e *DirectionalLightOld) FragmentType() string {
	return "DirectionalLightOld"
}

func (e *DirectionalLightOld) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
