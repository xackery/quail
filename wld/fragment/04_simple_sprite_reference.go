package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// SimpleSpriteReference information
type SimpleSpriteReference struct {
	flags        uint32
	textureCount uint32
}

func LoadSimpleSpriteReference(r io.ReadSeeker) (common.WldFragmenter, error) {
	e := &SimpleSpriteReference{}
	err := parseSimpleSpriteReference(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse SimpleSpriteReference: %w", err)
	}
	return e, nil
}

func parseSimpleSpriteReference(r io.ReadSeeker, e *SimpleSpriteReference) error {
	if e == nil {
		return fmt.Errorf("SimpleSpriteReference is nil")
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

func (e *SimpleSpriteReference) FragmentType() string {
	return "SimpleSpriteReference"
}

func (e *SimpleSpriteReference) Data() []byte {
	buf := bytes.NewBuffer(nil)

	binary.Write(buf, binary.LittleEndian, e.flags)
	binary.Write(buf, binary.LittleEndian, e.textureCount)
	return buf.Bytes()
}
