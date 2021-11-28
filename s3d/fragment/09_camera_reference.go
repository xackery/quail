package fragment

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// CameraReference information
type CameraReference struct {
}

func LoadCameraReference(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &CameraReference{}
	err := parseCameraReference(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse CameraReference: %w", err)
	}
	return l, nil
}

func parseCameraReference(r io.ReadSeeker, l *CameraReference) error {
	if l == nil {
		return fmt.Errorf("CameraReference is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (l *CameraReference) FragmentType() string {
	return "CameraReference"
}
