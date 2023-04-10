package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// TwoDSpriteDef information
type TwoDSpriteDef struct {
}

// LoadTwoDSpriteDef loads a TwoDSpriteDef fragment
func LoadTwoDSpriteDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &TwoDSpriteDef{}
	err := parseTwoDSpriteDef(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse TwoDSpriteDef: %w", err)
	}
	return e, nil
}

// parseTwoDSpriteDef parses a TwoDSpriteDef fragment
func parseTwoDSpriteDef(r io.ReadSeeker, e *TwoDSpriteDef) error {
	if e == nil {
		return fmt.Errorf("TwoDSpriteDef is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

// FragmentType returns the fragment type
func (e *TwoDSpriteDef) FragmentType() string {
	return "TwoDSpriteDef"
}

// Data returns the fragment data
func (e *TwoDSpriteDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
