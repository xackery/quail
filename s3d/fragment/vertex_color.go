package fragment

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
)

// VertexColor information
type VertexColor struct {
	// Colors of the vertex, if applicable
	Colors    []color.RGBA
	hashIndex uint32
}

func loadVertexColor(r io.ReadSeeker) (Fragment, error) {
	v := &VertexColor{}
	err := parseVertexColor(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse VertexColor: %w", err)
	}
	return v, nil
}

func parseVertexColor(r io.ReadSeeker, v *VertexColor) error {
	if v == nil {
		return fmt.Errorf("VertexColor is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &v.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	//unknown
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown: %w", err)
	}

	var vertexColorCount int32
	err = binary.Read(r, binary.LittleEndian, &vertexColorCount)
	if err != nil {
		return fmt.Errorf("read vertex color count: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown 2: %w", err)
	}
	//if value != 1 {
	//	return fmt.Errorf("unknown 2 expected %d, got %d", 1, value)
	//}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown 3: %w", err)
	}
	//also got 70? with nexus
	//if value != 200 {
	//	return fmt.Errorf("unknown 3 expected %d, got %d", 200, value)
	//}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown 4: %w", err)
	}
	//also got 74? with nexus
	//if value != 0 {
	//	return fmt.Errorf("unknown 4 expected %d, got %d", 0, value)
	//}
	for i := 0; i < int(vertexColorCount)/4; i++ {
		rgba := color.RGBA{}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read r: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read g: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read b: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read a: %w", err)
		}
		v.Colors = append(v.Colors, rgba)
	}

	return nil
}

func (v *VertexColor) FragmentType() string {
	return "Vertex Color"
}
