package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// CompositeSprite information
type CompositeSprite struct {
	name      string
	Reference uint32
}

func LoadCompositeSprite(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &CompositeSprite{}
	err := parseCompositeSprite(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse CompositeSprite: %w", err)
	}
	return v, nil
}

func parseCompositeSprite(r io.ReadSeeker, v *CompositeSprite) error {
	if v == nil {
		return fmt.Errorf("CompositeSprite is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value4: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value12: %w", err)
	}

	return nil
}

func (v *CompositeSprite) FragmentType() string {
	return "CompositeSprite"
}

func (e *CompositeSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
