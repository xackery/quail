package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// MaterialList information
type MaterialList struct {
	name               string
	MaterialReferences []uint32
}

func LoadMaterialList(r io.ReadSeeker) (archive.WldFragmenter, error) {
	m := &MaterialList{}
	err := parseMaterialList(r, m)
	if err != nil {
		return nil, fmt.Errorf("parse MaterialList: %w", err)
	}
	return m, nil
}

func parseMaterialList(r io.ReadSeeker, v *MaterialList) error {
	if v == nil {
		return fmt.Errorf("MaterialList is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	//TODO: flags support

	var materialCount uint32
	err = binary.Read(r, binary.LittleEndian, &materialCount)
	if err != nil {
		return fmt.Errorf("read materialCount: %w", err)
	}

	for i := uint32(0); i < materialCount; i++ {
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read %d materialReference: %w", i, err)
		}
		v.MaterialReferences = append(v.MaterialReferences, value)
	}
	return nil
}

func (m *MaterialList) FragmentType() string {
	return "Material List"
}

func (e *MaterialList) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
