package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// FourDSprite information
type FourDSprite struct {
}

func LoadFourDSprite(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &FourDSprite{}
	err := parseFourDSprite(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse FourDSprite: %w", err)
	}
	return e, nil
}

func parseFourDSprite(r io.ReadSeeker, e *FourDSprite) error {
	if e == nil {
		return fmt.Errorf("FourDSprite is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *FourDSprite) FragmentType() string {
	return "FourDSprite"
}

func (e *FourDSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
