package fragment

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// MeshAnimatedVertices information
type MeshAnimatedVertices struct {
}

func LoadMeshAnimatedVertices(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &MeshAnimatedVertices{}
	err := parseMeshAnimatedVertices(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse MeshAnimatedVertices: %w", err)
	}
	return l, nil
}

func parseMeshAnimatedVertices(r io.ReadSeeker, l *MeshAnimatedVertices) error {
	if l == nil {
		return fmt.Errorf("MeshAnimatedVertices is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (l *MeshAnimatedVertices) FragmentType() string {
	return "MeshAnimatedVertices"
}
