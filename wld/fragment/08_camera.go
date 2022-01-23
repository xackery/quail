package fragment

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// Camera information
type Camera struct {
}

func LoadCamera(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &Camera{}
	err := parseCamera(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse Camera: %w", err)
	}
	return l, nil
}

func parseCamera(r io.ReadSeeker, l *Camera) error {
	if l == nil {
		return fmt.Errorf("Camera is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (l *Camera) FragmentType() string {
	return "Camera"
}
