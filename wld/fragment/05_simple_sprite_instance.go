package fragment

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// SimpleSpriteInstance information
type SimpleSpriteInstance struct {
}

func LoadSimpleSpriteInstance(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &SimpleSpriteInstance{}
	err := parseSimpleSpriteInstance(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse SimpleSpriteInstance: %w", err)
	}
	return l, nil
}

func parseSimpleSpriteInstance(r io.ReadSeeker, l *SimpleSpriteInstance) error {
	if l == nil {
		return fmt.Errorf("SimpleSpriteInstance is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (l *SimpleSpriteInstance) FragmentType() string {
	return "SimpleSpriteInstance"
}
