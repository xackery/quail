package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// BlitSprite information
type BlitSprite struct {
	name      string
	Reference uint32
}

func LoadBlitSprite(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &BlitSprite{}
	err := parseBlitSprite(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse BlitSprite: %w", err)
	}
	return e, nil
}

func parseBlitSprite(r io.ReadSeeker, v *BlitSprite) error {
	if v == nil {
		return fmt.Errorf("BlitSprite is nil")
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

func (e *BlitSprite) FragmentType() string {
	return "BlitSprite"
}

func (e *BlitSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
