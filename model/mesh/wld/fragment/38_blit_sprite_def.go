package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// BlitSpriteDef information
type BlitSpriteDef struct {
	name      string
	Reference uint32
}

func LoadBlitSpriteDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &BlitSpriteDef{}
	err := parseBlitSpriteDef(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse BlitSpriteDef: %w", err)
	}
	return e, nil
}

func parseBlitSpriteDef(r io.ReadSeeker, v *BlitSpriteDef) error {
	if v == nil {
		return fmt.Errorf("BlitSpriteDef is nil")
	}
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	return nil
}

func (e *BlitSpriteDef) FragmentType() string {
	return "BlitSpriteDef"
}

func (e *BlitSpriteDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
