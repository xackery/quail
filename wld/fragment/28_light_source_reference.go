package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// LightSourceReference information
type LightSourceReference struct {
	hashIndex uint32
	Reference uint32
}

func LoadLightSourceReference(r io.ReadSeeker) (common.WldFragmenter, error) {
	e := &LightSourceReference{}
	err := parseLightSourceReference(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse LightSourceReference: %w", err)
	}
	return e, nil
}

func parseLightSourceReference(r io.ReadSeeker, e *LightSourceReference) error {
	if e == nil {
		return fmt.Errorf("lightsourceReference is nil")
	}
	err := binary.Read(r, binary.LittleEndian, &e.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &e.Reference)
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
