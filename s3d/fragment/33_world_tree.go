package fragment

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// WorldTree information
type WorldTree struct {
}

func LoadWorldTree(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &WorldTree{}
	err := parseWorldTree(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse WorldTree: %w", err)
	}
	return l, nil
}

func parseWorldTree(r io.ReadSeeker, l *WorldTree) error {
	if l == nil {
		return fmt.Errorf("WorldTree is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (l *WorldTree) FragmentType() string {
	return "WorldTree"
}
