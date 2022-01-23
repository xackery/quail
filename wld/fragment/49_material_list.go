package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// MaterialList information
type MaterialList struct {
	hashIndex          uint32
	MaterialReferences []uint32
}

func LoadMaterialList(r io.ReadSeeker) (common.WldFragmenter, error) {
	m := &MaterialList{}
	err := parseMaterialList(r, m)
	if err != nil {
		return nil, fmt.Errorf("parse MaterialList: %w", err)
	}
	return m, nil
}

func parseMaterialList(r io.ReadSeeker, m *MaterialList) error {
	if m == nil {
		return fmt.Errorf("MaterialList is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &m.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
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
		m.MaterialReferences = append(m.MaterialReferences, value)
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
