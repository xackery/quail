package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

// LegacyMesh information
type LegacyMesh struct {
	HashIndex         uint32
	Flags             uint32
	VertexCount       uint32
	TexCoordCount     uint32
	NormalCount       uint32
	ColorCount        uint32
	PolygonCount      uint32
	Size6             int16
	Fragment1Maybe    int16
	VertexPieceCount  uint32
	MaterialReference uint32
	Fragment3         uint32
	CenterPosition    math32.Vector3
	Params2           uint32
	Something2        uint32
	Something3        uint32
	verticies         []math32.Vector3
	texCoords         []math32.Vector3
	normals           []math32.Vector3
	colors            []int32
	polygons          []*LegacyPolygon
	vertexPieces      []*LegacyVertexPiece
	renderGroups      []*LegacyRenderGroup
	vertexTex         []*LegacyVertexTex
}

type LegacyPolygon struct {
	Flag int16
	Unk1 int16
	Unk2 int16
	Unk3 int16
	Unk4 int16
	I1   int16
	I2   int16
	I3   int16
}

type LegacyVertexPiece struct {
	Count  int16
	Offset int16
}

type LegacyRenderGroup struct {
	PolygonCount int16
	MaterialID   int16
}

type LegacyVertexTex struct {
	X int16
	Y int16
}

func LoadLegacyMesh(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &LegacyMesh{}
	err := parseLegacyMesh(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse LegacyMesh: %w", err)
	}
	return l, nil
}

func parseLegacyMesh(r io.ReadSeeker, l *LegacyMesh) error {
	if l == nil {
		return fmt.Errorf("LegacyMesh is nil")
	}
	err := binary.Read(r, binary.LittleEndian, l)
	if err != nil {
		return fmt.Errorf("read legacy mesh: %w", err)
	}
	for i := 0; i < int(l.VertexCount); i++ {
		err := binary.Read(r, binary.LittleEndian, &l.verticies)
		if err != nil {
			return fmt.Errorf("read vertex %d: %w", i, err)
		}

	}
	for i := 0; i < int(l.TexCoordCount); i++ {
		err := binary.Read(r, binary.LittleEndian, &l.texCoords)
		if err != nil {
			return fmt.Errorf("read tex coords %d: %w", i, err)
		}
	}
	for i := 0; i < int(l.NormalCount); i++ {
		err := binary.Read(r, binary.LittleEndian, &l.normals)
		if err != nil {
			return fmt.Errorf("read normal %d: %w", i, err)
		}
	}

	for i := 0; i < int(l.ColorCount); i++ {
		err := binary.Read(r, binary.LittleEndian, &l.colors)
		if err != nil {
			return fmt.Errorf("read color %d: %w", i, err)
		}
	}

	for i := 0; i < int(l.PolygonCount); i++ {
		p := &LegacyPolygon{}
		err := binary.Read(r, binary.LittleEndian, p)
		if err != nil {
			return fmt.Errorf("read polygon %d: %w", i, err)
		}
		l.polygons = append(l.polygons, p)
	}

	var value uint32
	for i := 0; i < int(l.Size6); i++ {
		err := binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read unk1 %d: %w", i, err)
		}
		if value != 4 {
			err = binary.Read(r, binary.LittleEndian, &value)
			if err != nil {
				return fmt.Errorf("read unk2 %d: %w", i, err)
			}
			err = binary.Read(r, binary.LittleEndian, &value)
			if err != nil {
				return fmt.Errorf("read unk3 %d: %w", i, err)
			}
		} else {
			err = binary.Read(r, binary.LittleEndian, &value)
			if err != nil {
				return fmt.Errorf("read unk4 %d: %w", i, err)
			}
			err = binary.Read(r, binary.LittleEndian, &value)
			if err != nil {
				return fmt.Errorf("read unk5 %d: %w", i, err)
			}
		}
	}

	for i := 0; uint32(i) < l.VertexPieceCount; i++ {
		vp := &LegacyVertexPiece{}
		err = binary.Read(r, binary.LittleEndian, vp)
		if err != nil {
			return fmt.Errorf("read vertex piece %d: %w", i, err)
		}
		l.vertexPieces = append(l.vertexPieces, vp)
	}

	if l.Flags&9 == 9 {
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read size8: %w", err)
		}
	}

	if l.Flags&11 == 11 {
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read polygonTexCount: %w", err)
		}
		for i := 0; uint32(i) < value; i++ {
			rg := &LegacyRenderGroup{}
			err = binary.Read(r, binary.LittleEndian, &rg)
			if err != nil {
				return fmt.Errorf("read render group %d: %w", i, err)
			}
			l.renderGroups = append(l.renderGroups, rg)
		}
	}
	if l.Flags&12 == 12 {
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read vertex count: %w", err)
		}
		for i := 0; uint32(i) < value; i++ {
			v := &LegacyVertexTex{}
			err = binary.Read(r, binary.LittleEndian, v)
			if err != nil {
				return fmt.Errorf("read vertex tex %d: %w", i, err)
			}
			l.vertexTex = append(l.vertexTex, v)
		}
	}

	if l.Flags&13 == 13 {
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read params3_1: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read params3_2: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read params3_3: %w", err)
		}
	}

	return nil
}

func (l *LegacyMesh) FragmentType() string {
	return "Legacy Mesh"
}
