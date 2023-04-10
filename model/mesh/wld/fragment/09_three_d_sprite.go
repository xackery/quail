package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// ThreeDSprite information
type ThreeDSprite struct {
}

func LoadThreeDSprite(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &ThreeDSprite{}
	err := parseThreeDSprite(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse ThreeDSprite: %w", err)
	}
	return e, nil
}

func parseThreeDSprite(r io.ReadSeeker, e *ThreeDSprite) error {
	if e == nil {
		return fmt.Errorf("ThreeDSprite is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *ThreeDSprite) FragmentType() string {
	return "ThreeDSprite"
}

func (e *ThreeDSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
