package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// Sound information
type Sound struct {
	name      string
	Reference uint32
}

func LoadSound(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &Sound{}
	err := parseSound(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse Sound: %w", err)
	}
	return e, nil
}

func parseSound(r io.ReadSeeker, v *Sound) error {
	if v == nil {
		return fmt.Errorf("Sound is nil")
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

func (e *Sound) FragmentType() string {
	return "Sound"
}

func (e *Sound) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
