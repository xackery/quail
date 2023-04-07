package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// CameraReference information
type CameraReference struct {
}

func LoadCameraReference(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &CameraReference{}
	err := parseCameraReference(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse CameraReference: %w", err)
	}
	return e, nil
}

func parseCameraReference(r io.ReadSeeker, e *CameraReference) error {
	if e == nil {
		return fmt.Errorf("CameraReference is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *CameraReference) FragmentType() string {
	return "CameraReference"
}

func (e *CameraReference) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
