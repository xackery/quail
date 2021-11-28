package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// GlobalAmbientLight information
type GlobalAmbientLight struct {
	Unk1 uint32
}

func LoadGlobalAmbientLight(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &GlobalAmbientLight{}
	err := parseGlobalAmbientLight(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse GlobalAmbientLight: %w", err)
	}
	return l, nil
}

func parseGlobalAmbientLight(r io.ReadSeeker, l *GlobalAmbientLight) error {
	if l == nil {
		return fmt.Errorf("globalAmbientLight is nil")
	}

	err := binary.Read(r, binary.LittleEndian, l)
	if err != nil {
		return fmt.Errorf("read light source : %w", err)
	}
	return nil
}

func (l *GlobalAmbientLight) FragmentType() string {
	return "Global Ambient Light"
}
