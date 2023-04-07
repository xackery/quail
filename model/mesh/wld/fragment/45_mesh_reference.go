package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// MeshReference information
type MeshReference struct {
	name      string
	Reference uint32
	Name      string
	Position  [3]float32
	Rotation  [3]float32
	Scale     [3]float32
}

func LoadMeshReference(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &MeshReference{}
	err := parseMeshReference(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse mesh reference: %w", err)
	}
	return v, nil
}

func parseMeshReference(r io.ReadSeeker, v *MeshReference) error {
	if v == nil {
		return fmt.Errorf("mesh reference is nil")
	}

	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHasIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	return nil
}

func (v *MeshReference) FragmentType() string {
	return "Mesh Reference"
}

func (e *MeshReference) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
