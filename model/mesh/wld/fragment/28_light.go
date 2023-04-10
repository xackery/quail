package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// Light information
type Light struct {
	name      string
	Reference uint32
}

func LoadLight(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &Light{}
	err := parseLight(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse Light: %w", err)
	}
	return e, nil
}

func parseLight(r io.ReadSeeker, v *Light) error {
	if v == nil {
		return fmt.Errorf("light is nil")
	}
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	return nil
}

func (e *Light) FragmentType() string {
	return "Light"
}

func (e *Light) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
