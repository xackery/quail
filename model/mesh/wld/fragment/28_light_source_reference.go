package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// LightSourceReference information
type LightSourceReference struct {
	name      string
	Reference uint32
}

func LoadLightSourceReference(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &LightSourceReference{}
	err := parseLightSourceReference(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse LightSourceReference: %w", err)
	}
	return e, nil
}

func parseLightSourceReference(r io.ReadSeeker, v *LightSourceReference) error {
	if v == nil {
		return fmt.Errorf("lightsourceReference is nil")
	}
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHasIndex: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	return nil
}

func (e *LightSourceReference) FragmentType() string {
	return "Light Source Reference"
}

func (e *LightSourceReference) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
