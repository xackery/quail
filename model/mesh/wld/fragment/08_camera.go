package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// Camera information
type Camera struct {
}

func LoadCamera(r io.ReadSeeker) (common.WldFragmenter, error) {
	e := &Camera{}
	err := parseCamera(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse Camera: %w", err)
	}
	return e, nil
}

func parseCamera(r io.ReadSeeker, e *Camera) error {
	if e == nil {
		return fmt.Errorf("Camera is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *Camera) FragmentType() string {
	return "Camera"
}

func (e *Camera) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
