package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// mesh 0x36 54
type mesh struct {
	isNewWorldFormat bool `bin:"-"`
	NameRef          int32
	Flags            uint32
	// A reference to a [MaterialListFragment] fragment. This tells the client which materials this mesh uses.
	// For zone meshes the [MaterialListFragment] contains all the materials used in the entire zone.
	// For placeable objects the [MaterialListFragment] contains all of the materials used in that object.
	MaterialListRef            int32
	AnimationRef               int32
	Fragment3Ref               int32 // unknown, usually empty
	Fragment4Ref               int32 // unknown, usually ref to first texture
	Center                     geo.Vector3
	Params2                    geo.UIndex3 // unknown, usually 0,0,0
	MaxDistance                float32     // radius from center, max distance from center
	Min                        geo.Vector3 // min x,y,z
	Max                        geo.Vector3 // max x,y,z
	VertexCount                uint16      //
	UvCount                    uint16
	NormalCount                uint16
	ColorCount                 uint16 // Vertex colors
	TriangleCount              uint16
	SkinAssignmentGroupCount   uint16 // used for bones
	TriangleMaterialGroupCount uint16
	VertexMaterialGroupCount   uint16
	MeshopCount                uint16         // used for animated mshes
	Scale                      float32        `bin:"ScaleUnmarshal,le"`
	Vertices                   []geo.Vector3  `bin:"VerticesUnmarshal,le"`
	Uvs                        []geo.Vector2  `bin:"UvsUnmarshal,le"`
	Normals                    []geo.UIndex3  `bin:"NormalsUnmarshal,le"`
	Colors                     []geo.RGBA     `bin:"len:ColorCount"`
	Triangles                  []geo.Triangle `bin:"TrianglesUnmarshal,le"`
	skinAssignmentGroups       []geo.Vector2  `bin:"-"`
	faceMaterialGroups         []geo.Vector2  `bin:"-"`
	VertexMaterialGroups       []geo.Vector2  `bin:"-"`
	meshops                    []geo.Vector2  `bin:"-"`
}

func (e *WLD) meshRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &mesh{
		isNewWorldFormat: !e.isOldWorld,
	}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("meshRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *mesh) build(e *WLD) error {
	name, ok := e.names[-v.NameRef]
	if !ok {
		return fmt.Errorf("offset 0x%x (len %d)", -v.NameRef, len(e.names))
	}

	m := &geo.Mesh{
		Name: name,
	}

	m.Triangles = v.Triangles

	e.meshManager.Add(m)

	return nil
}

/*
// UvUnmarhsal reads a uv
func (v *mesh) UvsUnmarshal(r binstruct.Reader) error {

	for i := 0; i < int(v.UvCount); i++ {
		uv := geo.Vector2{}
		if v.isNewWorldFormat {
			val32, err := r.ReadInt32()
			if err != nil {
				return fmt.Errorf("read uv x: %w", err)
			}
			uv.X = float32(val32) / 256

			val32, err = r.ReadInt32()
			if err != nil {
				return fmt.Errorf("read uv y: %w", err)
			}
			uv.Y = float32(val32) / 256
		} else {
			val16, err := r.ReadUint16()
			if err != nil {
				return fmt.Errorf("read uv x: %w", err)
			}
			uv.X = float32(val16) / 256

			val16, err = r.ReadUint16()
			if err != nil {
				return fmt.Errorf("read uv y: %w", err)
			}
			uv.Y = float32(val16) / 256
		}
		v.Uvs = append(v.Uvs, uv)
	}
	return nil
}

// NormalsUnmarshal reads a normal
func (v *mesh) NormalsUnmarshal(r binstruct.Reader) error {
	for i := 0; i < int(v.NormalCount); i++ {
		normal := geo.UIndex3{}
		val8, err := r.ReadInt8()
		if err != nil {
			return fmt.Errorf("read normal x: %w", err)
		}
		normal.X = uint32(val8) / 128

		val8, err = r.ReadInt8()
		if err != nil {
			return fmt.Errorf("read normal y: %w", err)
		}
		normal.Y = uint32(val8) / 128

		val8, err = r.ReadInt8()
		if err != nil {
			return fmt.Errorf("read normal z: %w", err)
		}
		normal.Z = uint32(val8) / 128

		v.Normals = append(v.Normals, normal)
	}
	return nil
}

// ScaleUnmarshal reads a scale
func (v *mesh) ScaleUnmarshal(r binstruct.Reader) error {
	scaleRaw, err := r.ReadUint32()
	if err != nil {
		return fmt.Errorf("read scale: %w", err)
	}
	v.Scale = float32(1 / float32(int(scaleRaw)<<16))
	return nil
}

// VerticesUnmarshal reads a vertex
func (v *mesh) VerticesUnmarshal(r binstruct.Reader) error {
	for i := 0; i < int(v.VertexCount); i++ {
		vertPos := geo.Vector3{}
		val16, err := r.ReadUint16()
		if err != nil {
			return fmt.Errorf("read vertex x: %w", err)
		}
		vertPos.X = float32(v.Center.X) + float32(val16)*v.Scale

		val16, err = r.ReadUint16()
		if err != nil {
			return fmt.Errorf("read vertex y: %w", err)
		}
		vertPos.Y = float32(v.Center.Y) + float32(val16)*v.Scale

		val16, err = r.ReadUint16()
		if err != nil {
			return fmt.Errorf("read vertex z: %w", err)
		}
		vertPos.Z = float32(v.Center.Z) + float32(val16)*v.Scale

		v.Vertices = append(v.Vertices, vertPos)
	}
	return nil
}

// TrianglesUnmarshal reads a triangle
func (v *mesh) TrianglesUnmarshal(r binstruct.Reader) error {
	for i := 0; i < int(v.TriangleCount); i++ {
		triangle := geo.Triangle{}

		flags, err := r.ReadInt16()
		if err != nil {
			return fmt.Errorf("read triangle flags: %w", err)
		}

		if flags == 0 {
			triangle.Flag = 1
		}

		val16, err := r.ReadUint16()
		if err != nil {
			return fmt.Errorf("read triangle x: %w", err)
		}
		triangle.Index.X = uint32(val16)

		val16, err = r.ReadUint16()
		if err != nil {
			return fmt.Errorf("read triangle y: %w", err)
		}
		triangle.Index.Y = uint32(val16)

		val16, err = r.ReadUint16()
		if err != nil {
			return fmt.Errorf("read triangle z: %w", err)
		}
		triangle.Index.Z = uint32(val16)

		v.Triangles = append(v.Triangles, triangle)
	}
	return nil
}
*/
