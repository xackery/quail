package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// SimpleSprite information
type SimpleSprite struct {
	TextureCount uint32
	Textures     []string
}

func LoadSimpleSprite(r io.ReadSeeker) (common.WldFragmenter, error) {
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

	err := binary.Read(r, binary.LittleEndian, &l.TextureCount)
	if err != nil {
		return fmt.Errorf("read texture count: %w", err)
	}
	var nameLength uint16
	for i := 0; i < int(l.TextureCount); i++ {
		//log.Infof("%d/%d\n", i, l.TextureCount)
		err = binary.Read(r, binary.LittleEndian, &nameLength)
		if err != nil {
			return fmt.Errorf("read name length: %w", err)
		}
		helper.ParseFixedString(r, uint32(nameLength))

	}
	return nil
}

func (l *SimpleSprite) FragmentType() string {
	return "SimpleSprite"
}
