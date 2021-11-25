package fragment

import (
	"fmt"
	"io"
)

// Camera information
type Camera struct {
}

func loadCamera(r io.ReadSeeker) (Fragment, error) {
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
