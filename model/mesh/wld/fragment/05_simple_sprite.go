package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// SimpleSprite information
type SimpleSprite struct {
}

func LoadSimpleSprite(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &SimpleSprite{}
	err := parseSimpleSprite(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse SimpleSprite: %w", err)
	}
	return e, nil
}

func parseSimpleSprite(r io.ReadSeeker, e *SimpleSprite) error {
	if e == nil {
		return fmt.Errorf("SimpleSprite is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *SimpleSprite) FragmentType() string {
	return "SimpleSprite"
}

func (e *SimpleSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
