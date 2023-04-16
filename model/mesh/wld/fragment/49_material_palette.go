package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// MaterialPalette information
type MaterialPalette struct {
	name               string
	MaterialReferences []uint32
}

func LoadMaterialPalette(r io.ReadSeeker) (archive.WldFragmenter, error) {
	m := &MaterialPalette{}
	err := parseMaterialPalette(r, m)
	if err != nil {
		return nil, fmt.Errorf("parse MaterialPalette: %w", err)
	}
	return m, nil
}

func parseMaterialPalette(r io.ReadSeeker, v *MaterialPalette) error {
	if v == nil {
		return fmt.Errorf("MaterialPalette is nil")
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

func (m *MaterialPalette) FragmentType() string {
	return "MaterialPalette"
}

func (e *MaterialPalette) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
