package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// DmSprite information
type DmSprite struct {
	name      string
	Reference uint32
	Name      string
	Position  [3]float32
	Rotation  [3]float32
	Scale     [3]float32
}

func LoadDmSprite(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &DmSprite{}
	err := parseDmSprite(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse DmSprite: %w", err)
	}
	return v, nil
}

func parseDmSprite(r io.ReadSeeker, v *DmSprite) error {
	if v == nil {
		return fmt.Errorf("DmSprite is nil")
	}

	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	return nil
}

func (v *DmSprite) FragmentType() string {
	return "DmSprite"
}

func (e *DmSprite) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
