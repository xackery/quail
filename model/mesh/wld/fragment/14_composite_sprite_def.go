package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// CompositeSpriteDef information
type CompositeSpriteDef struct {
	name      string
	Reference uint32
}

func LoadCompositeSpriteDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &CompositeSpriteDef{}
	err := parseCompositeSpriteDef(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse CompositeSpriteDef: %w", err)
	}
	return v, nil
}

func parseCompositeSpriteDef(r io.ReadSeeker, v *CompositeSpriteDef) error {
	if v == nil {
		return fmt.Errorf("CompositeSpriteDef is nil")
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

func (v *CompositeSpriteDef) FragmentType() string {
	return "CompositeSpriteDef"
}

func (e *CompositeSpriteDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
