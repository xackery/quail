package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// SimpleSprite information
type SimpleSprite struct {
	TextureCount uint32
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

	err := binary.Read(r, binary.LittleEndian, &e.TextureCount)
	if err != nil {
		return fmt.Errorf("read texture count: %w", err)
	}
	/*
		var nameLength uint16
		for i := 0; i < int(l.TextureCount); i++ {
			//log.Infof("%d/%d\n", i, l.TextureCount)
			err = binary.Read(r, binary.LittleEndian, &nameLength)
			if err != nil {
				return fmt.Errorf("read name length: %w", err)
			}
			helper.ParseFixedString(r, uint32(nameLength))
		}*/
	return nil
}

func (e *SimpleSprite) FragmentType() string {
	return "SimpleSprite"
}

func (e *SimpleSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, e.TextureCount)
	return buf.Bytes()
}
