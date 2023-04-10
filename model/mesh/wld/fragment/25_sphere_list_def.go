package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// SphereListDef information
type SphereListDef struct {
	name      string
	Reference uint32
}

func LoadSphereListDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &SphereListDef{}
	err := parseSphereListDef(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse SphereListDef: %w", err)
	}
	return v, nil
}

func parseSphereListDef(r io.ReadSeeker, v *SphereListDef) error {
	if v == nil {
		return fmt.Errorf("SphereListDef is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value4: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value12: %w", err)
	}

	return nil
}

func (v *SphereListDef) FragmentType() string {
	return "SphereListDef"
}

func (e *SphereListDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
