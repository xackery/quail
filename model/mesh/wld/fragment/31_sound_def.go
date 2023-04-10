package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// SoundDef information
type SoundDef struct {
	name      string
	Reference uint32
}

func LoadSoundDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &SoundDef{}
	err := parseSoundDef(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse SoundDef: %w", err)
	}
	return e, nil
}

func parseSoundDef(r io.ReadSeeker, v *SoundDef) error {
	if v == nil {
		return fmt.Errorf("SoundDef is nil")
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

func (e *SoundDef) FragmentType() string {
	return "SoundDef"
}

func (e *SoundDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
