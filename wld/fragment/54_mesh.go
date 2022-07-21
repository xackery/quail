package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"github.com/xackery/quail/common"
)

// Mesh information
type Mesh struct {
	name               string
	MaterialReference  uint32
	AnimationReference uint32
	Center             [3]float32
	MaxDistance        float32
	MinPosition        [3]float32
	MaxPosition        [3]float32
	verticies          []*common.Vertex
	triangles          []*common.Triangle
}

func LoadMesh(r io.ReadSeeker) (common.WldFragmenter, error) {
	v := &Mesh{}
	err := parseMesh(r, v, false)
	if err != nil {
		return nil, fmt.Errorf("parse Mesh: %w", err)
	}
	return v, nil
}

func parseMesh(r io.ReadSeeker, v *Mesh, isNewWorldFormat bool) error {
	var err error
	var val8 int8
	var val16 int16
	var val32 int32

	if v == nil {
		return fmt.Errorf("mesh is nil")
	}
	var value uint32
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHasIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	if value != 0x00018003 && // Zone
		value != 0x00014003 { // Object
		return fmt.Errorf("unknown mesh type, got 0x%x", value)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaterialReference)
	if err != nil {
		return fmt.Errorf("read material reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.AnimationReference)
	if err != nil {
		return fmt.Errorf("read animation reference: %w", err)
	}
	//TODO: find fragment referred to here (MeshAnimation)

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown2: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Center)
	if err != nil {
		return fmt.Errorf("read center: %w", err)
	}

	// dword1-3 Seems to be related to lighting models? (torches, etc.)
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknowndword1: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknowndword2: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknowndword3: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxDistance)
	if err != nil {
		return fmt.Errorf("read max distance: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MinPosition)
	if err != nil {
		return fmt.Errorf("read min position: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxPosition)
	if err != nil {
		return fmt.Errorf("read max position: %w", err)
	}

	var vertexCount uint16
	err = binary.Read(r, binary.LittleEndian, &vertexCount)
	if err != nil {
		return fmt.Errorf("read vertex count: %w", err)
	}

	var textureCoordinateCount uint16
	err = binary.Read(r, binary.LittleEndian, &textureCoordinateCount)
	if err != nil {
		return fmt.Errorf("read texture coordinate count: %w", err)
	}

	var normalsCount uint16
	err = binary.Read(r, binary.LittleEndian, &normalsCount)
	if err != nil {
		return fmt.Errorf("read normals count: %w", err)
	}

	var colorsCount uint16
	err = binary.Read(r, binary.LittleEndian, &colorsCount)
	if err != nil {
		return fmt.Errorf("read colors count: %w", err)
	}

	var triangleCount uint16
	err = binary.Read(r, binary.LittleEndian, &triangleCount)
	if err != nil {
		return fmt.Errorf("read triangle count: %w", err)
	}

	var vertexPieceCount uint16
	err = binary.Read(r, binary.LittleEndian, &vertexPieceCount)
	if err != nil {
		return fmt.Errorf("read vertex piece count: %w", err)
	}

	var triangleTextureCount uint16
	err = binary.Read(r, binary.LittleEndian, &triangleTextureCount)
	if err != nil {
		return fmt.Errorf("read triangle texture count: %w", err)
	}

	var vertexTextureCount uint16
	err = binary.Read(r, binary.LittleEndian, &vertexTextureCount)
	if err != nil {
		return fmt.Errorf("read vertex texture count: %w", err)
	}

	var size9 uint16
	err = binary.Read(r, binary.LittleEndian, &size9)
	if err != nil {
		return fmt.Errorf("read size9: %w", err)
	}
	scale := float32(1 / float32(int(1)<<value))

	vPositions := [][3]float32{}
	for i := 0; i < int(vertexCount); i++ {

		pos := [3]float32{}
		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex x %d: %w", i, err)
		}

		pos[0] = float32(val16) * scale

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex y %d: %w", i, err)
		}

		pos[1] = float32(val16) * scale

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex z %d: %w", i, err)
		}
		pos[2] = float32(val16) * scale

		vPositions = append(vPositions, pos)
	}

	vUvs := [][2]float32{}
	for i := 0; i < int(textureCoordinateCount); i++ {

		uv := [2]float32{}
		if isNewWorldFormat {
			err = binary.Read(r, binary.LittleEndian, &val32)
			if err != nil {
				return fmt.Errorf("read texture coordinate 32 %d: %w", i, err)
			}
			uv[0] = float32(val32 / 256)
			err = binary.Read(r, binary.LittleEndian, &val32)
			if err != nil {
				return fmt.Errorf("read texture coordinate 32 %d: %w", i, err)
			}
			uv[1] = float32(val32 / 256)

			vUvs = append(vUvs, uv)
			continue
		}

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read texture coordinate 32 %d: %w", i, err)
		}
		uv[0] = float32(val16 / 256)
		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read texture coordinate 32 %d: %w", i, err)
		}
		uv[1] = float32(val16 / 256)

		vUvs = append(vUvs, uv)

	}

	vNormals := [][3]float32{}
	for i := 0; i < int(normalsCount); i++ {

		pos := [3]float32{}
		err = binary.Read(r, binary.LittleEndian, &val8)
		if err != nil {
			return fmt.Errorf("read normals x %d: %w", i, err)
		}
		pos[0] = float32(val8) / 128

		err = binary.Read(r, binary.LittleEndian, &val8)
		if err != nil {
			return fmt.Errorf("read normals y %d: %w", i, err)
		}
		pos[1] = float32(val8) / 128

		err = binary.Read(r, binary.LittleEndian, &val8)
		if err != nil {
			return fmt.Errorf("read normals z %d: %w", i, err)
		}
		pos[2] = float32(val8) / 128

		vNormals = append(vNormals, pos)
	}

	vTints := []color.RGBA{}
	for i := 0; i < int(colorsCount); i++ {
		c := color.RGBA{}

		err = binary.Read(r, binary.LittleEndian, &c.R)
		if err != nil {
			return fmt.Errorf("read color r %d: %w", i, err)
		}
		err = binary.Read(r, binary.LittleEndian, &c.G)
		if err != nil {
			return fmt.Errorf("read color g %d: %w", i, err)
		}
		err = binary.Read(r, binary.LittleEndian, &c.B)
		if err != nil {
			return fmt.Errorf("read color b %d: %w", i, err)
		}
		err = binary.Read(r, binary.LittleEndian, &c.A)
		if err != nil {
			return fmt.Errorf("read color a %d: %w", i, err)
		}
		vTints = append(vTints, c)
	}

	for i := range vPositions {
		v.verticies = append(v.verticies, &common.Vertex{
			Position: vPositions[i],
			Normal:   vNormals[i],
			Tint:     &common.Tint{R: vTints[i].R, G: vTints[i].G, B: vTints[i].B},
			Uv:       vUvs[i],
			Uv2:      vUvs[i],
		})
	}

	for i := 0; i < int(triangleCount); i++ {
		var notSolidFlag int16
		err = binary.Read(r, binary.LittleEndian, &notSolidFlag)
		if err != nil {
			return fmt.Errorf("read notSolidFlag %d: %w", i, err)
		}

		triangle := &common.Triangle{}
		if notSolidFlag == 0 {
			//TODO: export separate collision flag
			//p.IsSolid = true
			triangle.Flag = 1
		}
		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex1 %d: %w", i, err)
		}
		triangle.Index[0] = uint32(val16)

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex2 %d: %w", i, err)
		}
		triangle.Index[1] = uint32(val16)

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex3 %d: %w", i, err)
		}
		triangle.Index[2] = uint32(val16)

		v.triangles = append(v.triangles, triangle)
	}

	for i := 0; i < int(vertexPieceCount); i++ {
		triangleIndex := int16(0)
		err = binary.Read(r, binary.LittleEndian, &triangleIndex)
		if err != nil {
			return fmt.Errorf("read triangleIndex %d: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read materialID %d: %w", i, err)
		}

		//v.triangles[triangleIndex].MaterialName = fmt.Sprintf("%d", val16)
	}

	for i := 0; i < int(triangleTextureCount); i++ {
		triangleIndex := int16(0)
		err = binary.Read(r, binary.LittleEndian, &triangleIndex)
		if err != nil {
			return fmt.Errorf("read triangleIndex %d: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read materialID %d: %w", i, err)
		}
		//fmt.Println("TODO, fix materials?", triangleIndex)
		//v.triangles[triangleIndex-1].MaterialName = fmt.Sprintf("%d", val16)
	}

	return nil
}

func (v *Mesh) FragmentType() string {
	return "Mesh"
}
func (e *Mesh) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}

func (e *Mesh) Vertices() []*common.Vertex {
	return e.verticies
}

func (e *Mesh) Triangles() []*common.Triangle {
	return e.triangles
}
