package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// DmSpriteDef2 information
type DmSpriteDef2 struct {
	name               string
	flags              uint32 // 0x00018003 = Zone, 0x00014003 = Object, 0x3 = NPC, 0x0 = ?
	MaterialReference  uint32 // ref to material
	AnimationReference uint32
	Center             *geo.Vector3 // zone meshes x coordinate of center of mesh, for placeholder objects, local origin
	MaxDistance        float32      // based on center max distance between any vertex an that position, a radius
	MinPosition        *geo.Vector3 // min x,y,z coords in absolute coords of any vertex in mesh
	MaxPosition        *geo.Vector3 // max x,y,z coords in absolute coords of any vertex in mesh
	verticies          []*geo.Vertex
	triangles          []*geo.Triangle
}

func NewDmSpriteDef2() *DmSpriteDef2 {
	return &DmSpriteDef2{
		Center:      &geo.Vector3{},
		MinPosition: &geo.Vector3{},
		MaxPosition: &geo.Vector3{},
	}
}

func LoadDmSpriteDef2(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := NewDmSpriteDef2()
	err := parseDmSpriteDef2(r, v, false)
	if err != nil {
		return nil, fmt.Errorf("parse DmSpriteDef2: %w", err)
	}
	return v, nil
}

func parseDmSpriteDef2(r io.ReadSeeker, v *DmSpriteDef2, isNewWorldFormat bool) error {
	var err error
	var val8 int8
	var val16 int16
	var val32 int32

	if v == nil {
		return fmt.Errorf("DmSpriteDef2 is nil")
	}

	var value uint32
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.flags)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	if value != 0x00018003 && // Zone
		value != 0x00014003 && // Object
		value != 0x3 && // NPC
		value != 0x0 {
		return fmt.Errorf("unknown DmSpriteDef2 type, got 0x%x", value)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaterialReference)
	if err != nil {
		return fmt.Errorf("read material reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.AnimationReference)
	if err != nil {
		return fmt.Errorf("read animation reference: %w", err)
	}

	// fragment3, usually empty
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown: %w", err)
	}

	// fragment4, usually first ref in textureimagefragments
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown2: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, v.Center)
	if err != nil {
		return fmt.Errorf("read center: %w", err)
	}

	// typically 0
	// dword1-3 Seems to be related to lighting models? (torches, etc.)
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknowndword1: %w", err)
	}

	// typically 0
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknowndword2: %w", err)
	}

	// typically 0
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknowndword3: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxDistance)
	if err != nil {
		return fmt.Errorf("read max distance: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, v.MinPosition)
	if err != nil {
		return fmt.Errorf("read min position: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, v.MaxPosition)
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

	skinAssignmentGroupCount := uint16(0)
	err = binary.Read(r, binary.LittleEndian, &skinAssignmentGroupCount)
	if err != nil {
		return fmt.Errorf("read skin assignment group count: %w", err)
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

	var meshopCount uint16
	err = binary.Read(r, binary.LittleEndian, &meshopCount)
	if err != nil {
		return fmt.Errorf("read meshop count: %w", err)
	}

	/// This allows vertex coordinates to be stored as integral values instead of
	/// floating-point values, without losing precision based on mesh size. Vertex
	/// values are multiplied by (1 shl `scale`) and stored in the vertex entries.
	/// FPSCALE is the internal name.
	var scaleRaw int16
	err = binary.Read(r, binary.LittleEndian, &scaleRaw)
	if err != nil {
		return fmt.Errorf("read scaleRaw: %w", err)
	}
	scale := float32(1 / float32(int(scaleRaw)<<value))

	// TODO: hacky scale fix for mesh
	//scale /= 100

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

			//TODO: fix scale
			//uv.X *= scale
			//uv.Y *= scale

			vUvs = append(vUvs, uv)
			continue
		}

		// old world format

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

		//TODO: fix scale
		//uv.X *= scale
		//uv.Y *= scale

		vUvs = append(vUvs, uv)

	}

	vNormals := []*geo.Vector3{}
	for i := 0; i < int(normalsCount); i++ {

		pos := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &val8)
		if err != nil {
			return fmt.Errorf("read normals x %d: %w", i, err)
		}
		fmt.Println(float32(val8))
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

		err = binary.Read(r, binary.LittleEndian, tint)
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
		vert := geo.NewVertex()
		vert.Position = vPositions[i]
		vert.Normal = vNormals[i]
		vert.Tint = vTints[i]
		vert.Uv = vUvs[i]
		vert.Uv2 = vUvs[i]

		v.verticies = append(v.verticies, vert)
	}

	for i := 0; i < int(triangleCount); i++ {
		var notSolidFlag int16
		err = binary.Read(r, binary.LittleEndian, &notSolidFlag)
		if err != nil {
			return fmt.Errorf("read notSolidFlag %d: %w", i, err)
		}

		triangle := geo.NewTriangle()
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

	for i := 0; i < int(skinAssignmentGroupCount); i++ {
		/// The first element of the tuple is the number of vertices in a skeleton piece.
		///
		/// The second element of the tuple is the index of the piece according to the
		/// [SkeletonTrackSet] fragment. The very first piece (index 0) is usually not referenced here
		/// as it is usually jsut a "stem" starting point for the skeleton. Only those pieces
		/// referenced here in the mesh should actually be rendered. Any other pieces in the skeleton
		/// contain no vertices or faces And have other purposes.

		// number of verts in a skelton piece
		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read count %d: %w", i, err)
		}
		// index of the piece according to skeletontrackset fragment
		err = binary.Read(r, binary.LittleEndian, &val16)
		if err != nil {
			return fmt.Errorf("read index1 %d: %w", i, err)
		}

		//v.triangles[triangleIndex].MaterialName = fmt.Sprintf("%d", val16)
	}

	for i := 0; i < int(triangleTextureCount); i++ {
		//TODO: resort
		/// The first element of the tuple is the number of faces that use the same material. All
		/// polygon entries are sorted by material index so that faces use the same material are
		/// grouped together.
		///
		/// The second element of the tuple is the index of the material that the faces use according
		/// to the [MaterialListFragment] that this fragment references.

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

		v.triangles[triangleIndex-1].MaterialName = fmt.Sprintf("%d", val16)
	}

	return nil
}

func (v *DmSpriteDef2) FragmentType() string {
	return "DmSpriteDef2"
}
func (e *DmSpriteDef2) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}

func (e *DmSpriteDef2) Vertices() []*geo.Vertex {
	return e.verticies
}

func (e *DmSpriteDef2) Triangles() []*geo.Triangle {
	return e.triangles
}

// Name returns the name of the fragment as specified in the header.
func (e *DmSpriteDef2) Name() string {
	return e.name
}
