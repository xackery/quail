package fragment

import (
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
	l := &LightSourceReference{}
	err := parseLightSourceReference(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse LightSourceReference: %w", err)
	}
	return l, nil
}

func parseLightSourceReference(r io.ReadSeeker, l *LightSourceReference) error {
	if l == nil {
		return fmt.Errorf("lightsourceReference is nil")
	}
	err := binary.Read(r, binary.LittleEndian, &l.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &l.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	return nil
}

func (l *LightSourceReference) FragmentType() string {
	return "Light Source Reference"
}
