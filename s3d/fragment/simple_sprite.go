package fragment

import (
	"fmt"
	"io"
)

// SimpleSprite information
type SimpleSprite struct {
}

func loadSimpleSprite(r io.ReadSeeker) (Fragment, error) {
	l := &SimpleSprite{}
	err := parseSimpleSprite(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse SimpleSprite: %w", err)
	}
	return l, nil
}

func parseSimpleSprite(r io.ReadSeeker, l *SimpleSprite) error {
	if l == nil {
		return fmt.Errorf("SimpleSprite is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (l *SimpleSprite) FragmentType() string {
	return "SimpleSprite"
}
