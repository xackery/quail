package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// Actor information
type Actor struct {
}

func LoadActor(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &Actor{}
	err := parseActor(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse Actor: %w", err)
	}
	return e, nil
}

func parseActor(r io.ReadSeeker, e *Actor) error {
	if e == nil {
		return fmt.Errorf("Actor is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *Actor) FragmentType() string {
	return "Actor"
}

func (e *Actor) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
