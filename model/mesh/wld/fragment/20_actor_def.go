package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// ActorDef information
type ActorDef struct {
}

func LoadActorDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &ActorDef{}
	err := parseActorDef(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse ActorDef: %w", err)
	}
	return e, nil
}

func parseActorDef(r io.ReadSeeker, e *ActorDef) error {
	if e == nil {
		return fmt.Errorf("ActorDef is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *ActorDef) FragmentType() string {
	return "ActorDef"
}

func (e *ActorDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
