package fragment

import (
	"encoding/binary"
	"fmt"
	"io"
)

// VertexColorReference, Referenced by an ObjectInstance fragment.
type VertexColorReference struct {
	VertexColor *VertexColor
	Reference   uint32
	hashIndex   uint32
}

func loadVertexColorReference(r io.ReadSeeker) (Fragment, error) {
	v := &VertexColorReference{}
	err := parseVertexColorReference(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse VertexColorReference: %w", err)
	}
	return v, nil
}

func parseVertexColorReference(r io.ReadSeeker, v *VertexColorReference) error {
	if v == nil {
		return fmt.Errorf("VertexColorReference is nil")
	}
	err := binary.Read(r, binary.LittleEndian, &v.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	return nil
}

func (v *VertexColorReference) FragmentType() string {
	return "Vertex Color Reference"
}
