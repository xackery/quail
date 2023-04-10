package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// TwoDSprite information
type TwoDSprite struct {
}

// LoadTwoDSprite loads a TwoDSprite fragment
func LoadTwoDSprite(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &TwoDSprite{}
	err := parseTwoDSprite(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse TwoDSprite: %w", err)
	}
	return e, nil
}

func parseTwoDSprite(r io.ReadSeeker, e *TwoDSprite) error {
	if e == nil {
		return fmt.Errorf("TwoDSprite is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

// FragmentType returns the fragment type
func (e *TwoDSprite) FragmentType() string {
	return "TwoDSprite"
}

// Data returns the fragment data
func (e *TwoDSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
