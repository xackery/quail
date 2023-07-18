package wld

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/quail/def"
)

// Mesh 0x36 54
type Mesh struct {
	isNewWorldFormat bool `bin:"-"`
	NameRef          int32
	Flags            uint32
	// A reference to a [MaterialListFragment] fragment. This tells the client which materials this mesh uses.
	// For zone meshes the [MaterialListFragment] contains all the materials used in the entire zone.
	// For placeable objects the [MaterialListFragment] contains all of the materials used in that object.
	MaterialListRef            uint32
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
	mesh := &def.Mesh{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	nameRef := dec.Int32()
	//nameLength := dec.Uint32()
	name, ok := e.names[-nameRef]
	if !ok {
		return fmt.Errorf("unknown name ref %d", nameRef)
	}
	mesh.Name = strings.ToLower(strings.TrimSuffix(name, "_DMSPRITEDEF"))
	flags := dec.Uint32() // flags, currently unknown, zone meshes are 0x00018003, placeable objects are 0x00014003
	dump.Hex(flags, "flags=%d", flags)

	materialListRef := dec.Int32()
	dump.Hex(materialListRef, "materialListRef=%d", materialListRef)
	animationRef := dec.Int32() //used by flags/trees only
	dump.Hex(animationRef, "animationRef=%d", animationRef)

	_ = dec.Uint32() // unknown, usually empty
	_ = dec.Uint32() // unknown, This usually seems to reference the first [TextureImagesFragment] fragment in the file.

	// for zone meshes, x coordinate of the center of the mesh
	// for placeable objects, this seems to define where the vertices will lie relative to the object’s local origin
	centerX := dec.Float32()
	centerY := dec.Float32()
	centerZ := dec.Float32()

	_ = dec.Uint32() // unknown, usually empty
	_ = dec.Uint32() // unknown, usually empty
	_ = dec.Uint32() // unknown, usually empty
	_ = dec.Uint32() // unknown, usually empty

	//maxDistance := dec.Float32 // Given the values in center, this seems to contain the maximum distance between any vertex and that position. It seems to define a radius from that position within which the mesh lies.
	minX := dec.Float32() // min x, y, and z coords in absolute coords of any vertex in the mesh.
	dump.Hex(minX, "minX=%f", minX)
	minY := dec.Float32()
	dump.Hex(minY, "minY=%f", minY)
	minZ := dec.Float32()
	dump.Hex(minZ, "minZ=%f", minZ)

	maxX := dec.Float32() // max x, y, and z coords in absolute coords of any vertex in the mesh.
	dump.Hex(maxX, "maxX=%f", maxX)
	maxY := dec.Float32()
	dump.Hex(maxY, "maxY=%f", maxY)
	maxZ := dec.Float32()
	dump.Hex(maxZ, "maxZ=%f", maxZ)

	vertexCount := dec.Uint16()   // number of vertices in the mesh (called position_count in libeq)
	uvCount := dec.Uint16()       // number of uv in the mesh (called texture_coordinate_count in libeq)
	normalCount := dec.Uint16()   // number of vertex normal entries in the mesh (called normal_count in libeq)
	colorCount := dec.Uint16()    // number of vertex color entries in the mesh (called color_count in libeq)
	triangleCount := dec.Uint16() // number of triangles in the mesh (called face_count in libeq)
	/// This seems to only be used when dealing with animated (mob) models.
	/// It contains the number of vertex piece entries. Vertices are grouped together by
	/// skeleton piece in this case and vertex piece entries tell the client how
	/// many vertices are in each piece. It’s possible that there could be more
	/// pieces in the skeleton than are in the meshes it references. Extra pieces have
	/// no faces or vertices and I suspect they are there to define attachment points for
	/// objects (e.g. weapons or shields).
	skinAssignmentGroupsCount := dec.Uint16()
	dump.Hex(skinAssignmentGroupsCount, "skinAssignmentGroupsCount=%d", skinAssignmentGroupsCount)
	triangleMaterialGroupsCount := dec.Uint16() // number of triangle texture entries. faces are grouped together by material and polygon material entries. This tells the client the number of faces using a material.
	dump.Hex(triangleMaterialGroupsCount, "triangleMaterialGroupsCount=%d", triangleMaterialGroupsCount)

	vertexMaterialGroupsCount := dec.Uint16() // number of vertex material entries. Vertices are grouped together by material and vertex material entries tell the client how many vertices there are using a material.
	dump.Hex(vertexMaterialGroupsCount, "vertexMaterialGroupsCount=%d", vertexMaterialGroupsCount)

	meshopCount := dec.Uint16() // number of entries in meshops. Seems to be used only for animated mob models.
	dump.Hex(meshopCount, "vertexCount=%d", meshopCount)

	scale := float32(1 / float32(int(1)<<dec.Uint16())) // This allows vertex coordinates to be stored as integral values instead of floating-point values, without losing precision based on mesh size. Vertex values are multiplied by (1 shl `scale`) and stored in the vertex entries. FPSCALE is the internal name.

	/// Vertices (x, y, z) belonging to this mesh. Each axis should
	/// be multiplied by (1 shl `scale`) for the final vertex position.
	for i := 0; i < int(vertexCount); i++ {
		vert := def.Vertex{}
		vert.Position.X = float32(centerX) + (float32(dec.Int16()) * scale)
		vert.Position.Y = float32(centerY) + (float32(dec.Int16()) * scale)
		vert.Position.Z = float32(centerZ) + (float32(dec.Int16()) * scale)
		mesh.Vertices = append(mesh.Vertices, vert)
	}

	for i := 0; i < int(uvCount); i++ {
		uv := def.Vector2{}
		if e.isOldWorld {
			uv.X = float32(dec.Int16()) / 256
			uv.Y = float32(dec.Int16()) / 256
		} else {
			uv.X = float32(dec.Int32()) / 256
			uv.Y = float32(dec.Int32()) / 256
		}
		mesh.Vertices[i].Uv = uv
	}

	for i := 0; i < int(normalCount); i++ {
		normal := def.Vector3{}
		normal.X = float32(dec.Int8()) / 128
		normal.Y = float32(dec.Int8()) / 128
		normal.Z = float32(dec.Int8()) / 128
		mesh.Vertices[i].Normal = normal
	}

	for i := 0; i < int(colorCount); i++ {
		color := def.RGBA{}
		color.R = dec.Uint8()
		color.G = dec.Uint8()
		color.B = dec.Uint8()
		color.A = dec.Uint8()
		mesh.Vertices[i].Tint = color
	}

	for i := 0; i < int(triangleCount); i++ {
		triangle := def.Triangle{}
		notSolidFlag := dec.Uint16()
		if notSolidFlag != 0 {
			triangle.Flag = 1
		}
		triangle.Index.X = uint32(dec.Uint16())
		triangle.Index.Y = uint32(dec.Uint16())
		triangle.Index.Z = uint32(dec.Uint16())
		mesh.Triangles = append(mesh.Triangles, triangle)
	}

	log.Debugf("mesh.Triangles=%+v count=%d", mesh.Triangles, triangleCount)

	if dec.Error() != nil {
		return fmt.Errorf("meshRead: %v", dec.Error())
	}

	log.Debugf("%+v", mesh)
	e.Fragments[fragmentOffset] = mesh
	return nil
}

func (v *Mesh) build(e *WLD) error {
	name, ok := e.names[-v.NameRef]
	if !ok {
		return fmt.Errorf("offset 0x%x (len %d)", -v.NameRef, len(e.names))
	}

	m := &geo.Mesh{
		Name: name,
	}

	m.Triangles = v.Triangles

	e.meshManager.Add(m)

	dm := &def.Mesh{
		Name: name,
	}

	e.meshes = append(e.meshes, dm)
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

func (e *WLD) meshWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
