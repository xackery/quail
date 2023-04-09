package fragment

import (
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// DefaultPalette is used by WLDCOM but is never known to be used by EQ
type DefaultPalette struct {
}

// LoadDefaultPalette loads a DefaultPalette
func LoadDefaultPalette(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &DefaultPalette{}
	err := parseDefaultPalette(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse DefaultPalette: %w", err)
	}
	return e, nil
}

func parseDefaultPalette(r io.ReadSeeker, e *DefaultPalette) error {
	if e == nil {
		return fmt.Errorf("DefaultPalette is nil")
	}

	return fmt.Errorf("DefaultPalette is not implemented")
}

func (e *DefaultPalette) FragmentType() string {
	return "DefaultPalette"
}

// Data returns the raw data of the fragment
func (e *DefaultPalette) Data() []byte {
	return nil
}
