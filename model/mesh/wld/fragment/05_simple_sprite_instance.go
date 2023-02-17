package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// SimpleSpriteInstance information
type SimpleSpriteInstance struct {
}

func LoadSimpleSpriteInstance(r io.ReadSeeker) (common.WldFragmenter, error) {
	e := &SimpleSpriteInstance{}
	err := parseSimpleSpriteInstance(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse SimpleSpriteInstance: %w", err)
	}
	return e, nil
}

func parseSimpleSpriteInstance(r io.ReadSeeker, e *SimpleSpriteInstance) error {
	if e == nil {
		return fmt.Errorf("SimpleSpriteInstance is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *SimpleSpriteInstance) FragmentType() string {
	return "SimpleSpriteInstance"
}

func (e *SimpleSpriteInstance) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
