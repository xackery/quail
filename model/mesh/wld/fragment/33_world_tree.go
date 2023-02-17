package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// WorldTree information
type WorldTree struct {
}

func LoadWorldTree(r io.ReadSeeker) (common.WldFragmenter, error) {
	e := &WorldTree{}
	err := parseWorldTree(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse WorldTree: %w", err)
	}
	return e, nil
}

func parseWorldTree(r io.ReadSeeker, e *WorldTree) error {
	if e == nil {
		return fmt.Errorf("WorldTree is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *WorldTree) FragmentType() string {
	return "WorldTree"
}

func (e *WorldTree) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
