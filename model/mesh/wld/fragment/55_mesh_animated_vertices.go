package fragment

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// MeshAnimatedVertices information
type MeshAnimatedVertices struct {
}

func LoadMeshAnimatedVertices(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &MeshAnimatedVertices{}
	err := parseMeshAnimatedVertices(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse MeshAnimatedVertices: %w", err)
	}
	return e, nil
}

func parseMeshAnimatedVertices(r io.ReadSeeker, e *MeshAnimatedVertices) error {
	if e == nil {
		return fmt.Errorf("MeshAnimatedVertices is nil")
	}
	/*
		err := binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			return fmt.Errorf("read light source : %w", err)
		}*/
	return nil
}

func (e *MeshAnimatedVertices) FragmentType() string {
	return "MeshAnimatedVertices"
}

func (e *MeshAnimatedVertices) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
