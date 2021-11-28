package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// AmbientLight information
type AmbientLight struct {
	Unk1 uint32
}

func LoadAmbientLight(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &AmbientLight{}
	err := parseAmbientLight(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse AmbientLight: %w", err)
	}
	return l, nil
}

func parseAmbientLight(r io.ReadSeeker, l *AmbientLight) error {
	if l == nil {
		return fmt.Errorf("AmbientLight is nil")
	}

	err := binary.Read(r, binary.LittleEndian, l)
	if err != nil {
		return fmt.Errorf("read light source : %w", err)
	}
	return nil
}

func (l *AmbientLight) FragmentType() string {
	return "Ambient Light"
}
