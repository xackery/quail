package fragment

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// Actor information
type Actor struct {
}

func LoadActor(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &Actor{}
	err := parseActor(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse Actor: %w", err)
	}
	return l, nil
}

func parseActor(r io.ReadSeeker, l *Actor) error {
	if l == nil {
		return fmt.Errorf("Actor is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (l *Actor) FragmentType() string {
	return "Actor"
}
