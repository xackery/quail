package fragment

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// SimpleSpriteReference information
type SimpleSpriteReference struct {
}

func LoadSimpleSpriteReference(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &SimpleSpriteReference{}
	err := parseSimpleSpriteReference(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse SimpleSpriteReference: %w", err)
	}
	return l, nil
}

func parseSimpleSpriteReference(r io.ReadSeeker, l *SimpleSpriteReference) error {
	if l == nil {
		return fmt.Errorf("SimpleSpriteReference is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (l *SimpleSpriteReference) FragmentType() string {
	return "SimpleSpriteReference"
}
