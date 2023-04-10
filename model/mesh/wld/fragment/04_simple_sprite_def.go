package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// SimpleSpriteDef information
type SimpleSpriteDef struct {
	flags        uint32
	textureCount uint32
}

func LoadSimpleSpriteDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &SimpleSpriteDef{}
	err := parseSimpleSpriteDef(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse SimpleSpriteDef: %w", err)
	}
	return e, nil
}

func parseSimpleSpriteDef(r io.ReadSeeker, e *SimpleSpriteDef) error {
	if e == nil {
		return fmt.Errorf("SimpleSpriteDef is nil")
	}

	err := binary.Read(r, binary.LittleEndian, &e.flags)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &e.textureCount)
	if err != nil {
		return fmt.Errorf("read textureCount: %w", err)
	}

	return nil
}

func (e *SimpleSpriteDef) FragmentType() string {
	return "SimpleSpriteDef"
}

func (e *SimpleSpriteDef) Data() []byte {
	buf := bytes.NewBuffer(nil)

	binary.Write(buf, binary.LittleEndian, e.flags)
	binary.Write(buf, binary.LittleEndian, e.textureCount)
	return buf.Bytes()
}
