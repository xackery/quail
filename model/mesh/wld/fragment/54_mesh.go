package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// Mesh information
type Mesh struct {
	name               string
	flags              uint32
	MaterialReference  uint32
	AnimationReference uint32
	Center             *geo.Vector3
	MaxDistance        float32
	MinPosition        *geo.Vector3
	MaxPosition        *geo.Vector3
	verticies          []*geo.Vertex
	triangles          []*geo.Triangle
}

func LoadMesh(r io.ReadSeeker) (archive.WldFragmenter, error) {
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

	err = binary.Read(r, binary.LittleEndian, &v.flags)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	if value != 0x00018003 && // Zone
		value != 0x00014003 && // Object
		value != 0x3 && // NPC
		value != 0x0 {
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

	var scaleRaw int16
	err = binary.Read(r, binary.LittleEndian, &scaleRaw)
	if err != nil {
		return fmt.Errorf("read scaleRaw: %w", err)
	}
	scale := float32(1 / float32(int(scaleRaw)<<value))

	vPositions := []*geo.Vector3{}
	for i := 0; i < int(vertexCount); i++ {

		pos := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex x %d: %w", i, err)
		}

		pos.X = float32(v.Center.X) + float32(val16)*scale

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex y %d: %w", i, err)
		}

		pos.Y = float32(v.Center.Y) + float32(val16)*scale

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex z %d: %w", i, err)
		}
		pos.Z = float32(v.Center.Z) + float32(val16)*scale

		vPositions = append(vPositions, pos)
	}

	vUvs := []*geo.Vector2{}
	for i := 0; i < int(textureCoordinateCount); i++ {

		uv := &geo.Vector2{}
		if isNewWorldFormat {
			err = binary.Read(r, binary.LittleEndian, &val32)
			if err != nil {
				return fmt.Errorf("read texture coordinate 32 %d: %w", i, err)
			}
			uv.X = float32(val32 / 256)
			err = binary.Read(r, binary.LittleEndian, &val32)
			if err != nil {
				return fmt.Errorf("read texture coordinate 32 %d: %w", i, err)
			}
			uv.Y = float32(val32 / 256)

			vUvs = append(vUvs, uv)
			continue
		}

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read texture coordinate 32 %d: %w", i, err)
		}
		uv.X = float32(val16 / 256)
		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read texture coordinate 32 %d: %w", i, err)
		}
		uv.Y = float32(val16 / 256)

		vUvs = append(vUvs, uv)

	}

	vNormals := []*geo.Vector3{}
	for i := 0; i < int(normalsCount); i++ {

		pos := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &val8)
		if err != nil {
			return fmt.Errorf("read normals x %d: %w", i, err)
		}
		pos.X = float32(val8) / float32(128)

		err = binary.Read(r, binary.LittleEndian, &val8)
		if err != nil {
			return fmt.Errorf("read normals y %d: %w", i, err)
		}
		pos.Y = float32(val8) / float32(128)

		err = binary.Read(r, binary.LittleEndian, &val8)
		if err != nil {
			return fmt.Errorf("read normals z %d: %w", i, err)
		}
		pos.Z = float32(val8) / float32(128)

		vNormals = append(vNormals, pos)
	}

	vTints := []*geo.RGBA{}
	for i := 0; i < int(colorsCount); i++ {

		tint := &geo.RGBA{}

		err = binary.Read(r, binary.LittleEndian, &tint)
		if err != nil {
			return fmt.Errorf("read color %d: %w", i, err)
		}
		vTints = append(vTints, tint)
	}
	for len(vPositions) > len(vTints) {
		vTints = append(vTints, &geo.RGBA{})
	}

	if len(vPositions) != len(vNormals) ||
		len(vNormals) != len(vTints) ||
		len(vTints) != len(vUvs) {
		return fmt.Errorf("mismatch on length of verticies")
	}

	for i := range vPositions {
		v.verticies = append(v.verticies, &geo.Vertex{
			Position: vPositions[i],
			Normal:   vNormals[i],
			Tint:     vTints[i],
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

		triangle := &geo.Triangle{}
		if notSolidFlag == 0 {
			//TODO: export separate collision flag
			//p.IsSolid = true
			triangle.Flag = 1
		}
		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex1 %d: %w", i, err)
		}
		triangle.Index.X = uint32(val16)

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex2 %d: %w", i, err)
		}
		triangle.Index.Y = uint32(val16)

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read vertex3 %d: %w", i, err)
		}
		triangle.Index.Z = uint32(val16)

		v.triangles = append(v.triangles, triangle)
	}

	for i := 0; i < int(vertexPieceCount); i++ {
		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read count %d: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read index1 %d: %w", i, err)
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

func (e *Mesh) Vertices() []*geo.Vertex {
	return e.verticies
}

func (e *Mesh) Triangles() []*geo.Triangle {
	return e.triangles
}
