package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// ThreeDSpriteDef information
type ThreeDSpriteDef struct {
}

func LoadThreeDSpriteDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &ThreeDSpriteDef{}
	err := parseThreeDSpriteDef(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse ThreeDSpriteDef: %w", err)
	}
	return e, nil
}

func parseThreeDSpriteDef(r io.ReadSeeker, e *ThreeDSpriteDef) error {
	if e == nil {
		return fmt.Errorf("ThreeDSpriteDef is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *ThreeDSpriteDef) FragmentType() string {
	return "ThreeDSpriteDef"
}

func (e *ThreeDSpriteDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
