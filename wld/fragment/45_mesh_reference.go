package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

// MeshReference information
type MeshReference struct {
	hashIndex uint32
	Reference uint32
	Name      string
	Position  math32.Vector3
	Rotation  math32.Vector3
	Scale     math32.Vector3
}

func LoadMeshReference(r io.ReadSeeker) (common.WldFragmenter, error) {
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

	err := binary.Read(r, binary.LittleEndian, &v.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	//TODO: name from hash

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	return nil
}

func (v *MeshReference) FragmentType() string {
	return "Mesh Reference"
}
