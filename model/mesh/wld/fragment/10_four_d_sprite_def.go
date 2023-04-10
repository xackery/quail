package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// FourDSpriteDef information
type FourDSpriteDef struct {
}

func LoadFourDSpriteDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &FourDSpriteDef{}
	err := parseFourDSpriteDef(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse FourDSpriteDef: %w", err)
	}
	return e, nil
}

func parseFourDSpriteDef(r io.ReadSeeker, e *FourDSpriteDef) error {
	if e == nil {
		return fmt.Errorf("FourDSpriteDef is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *FourDSpriteDef) FragmentType() string {
	return "FourDSpriteDef"
}

func (e *FourDSpriteDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
