package wld

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
)

// Mesh 0x36 54
type Mesh struct {
	NameRef  int32
	Name     string
	FileType string
	Flags    uint32
	// A reference to a [MaterialListFragment] fragment. This tells the client which materials this mesh uses.
	// For zone meshes the [MaterialListFragment] contains all the materials used in the entire zone.
	// For placeable objects the [MaterialListFragment] contains all of the materials used in that object.
	MaterialListRef   uint32
	AnimationRef      int32
	Fragment3Ref      int32 // unknown, usually empty
	Fragment4Ref      int32 // unknown, usually ref to first texture
	Center            common.Vector3
	Params2           common.UIndex3 // unknown, usually 0,0,0
	MaxDistance       float32        // radius from center, max distance from center
	Min               common.Vector3 // min x,y,z
	Max               common.Vector3 // max x,y,z
	MeshopCount       uint16         // used for animated mshes
	Scale             float32        `bin:"ScaleUnmarshal,le"`
	Vertices          []common.Vertex
	Uvs               []common.Vector2 `bin:"UvsUnmarshal,le"`
	Normals           []common.UIndex3 `bin:"NormalsUnmarshal,le"`
	Colors            []common.RGBA    `bin:"len:ColorCount"`
	TriangleMaterials []MeshTriangleMaterial
	Triangles         []common.Triangle
	VertexPieces      []MeshVertexPiece
	AnimatedBones     []MeshAnimatedBone
}

type MeshVertexPiece struct {
	Count  int16
	Index1 int16
}

type MeshTriangleMaterial struct {
	Count      uint16
	MaterialID uint16
}

type MeshAnimatedBone struct {
	Position common.Vector3
}

func (e *WLD) meshRead(r io.ReadSeeker, fragmentOffset int) error {
	mesh := &Mesh{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	nameRef := dec.Int32()
	//nameLength := dec.Uint32()
	name, ok := e.names[-nameRef]
	if !ok {
		return fmt.Errorf("unknown name ref %d", nameRef)
	}
	mesh.Name = strings.ToLower(strings.TrimSuffix(name, "_DMSPRITEDEF"))
	flags := dec.Uint32() // flags, currently unknown, zone meshes are 0x00018003, placeable objects are 0x00014003

	if flags == 0x00018003 {
		mesh.FileType = "ter"
	}

	if flags == 0x00014003 {
		mesh.FileType = "mod"
	}

	mesh.MaterialListRef = dec.Uint32()
	mesh.AnimationRef = dec.Int32() //used by flags/trees only

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

	mesh.MaxDistance = dec.Float32() // Given the values in center, this seems to contain the maximum distance between any vertex and that position. It seems to define a radius from that position within which the mesh lies.
	mesh.Min.X = dec.Float32()       // min x, y, and z coords in absolute coords of any vertex in the mesh.
	mesh.Min.Y = dec.Float32()
	mesh.Min.Z = dec.Float32()

	mesh.Max.X = dec.Float32() // max x, y, and z coords in absolute coords of any vertex in the mesh.
	mesh.Max.Y = dec.Float32()
	mesh.Max.Z = dec.Float32()

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
	vertexPieceCount := dec.Uint16()
	polygonTextureCount := dec.Uint16() // number of triangle texture entries. faces are grouped together by material and polygon material entries. This tells the client the number of faces using a material.
	vertexTextureCount := dec.Uint16()  // number of vertex material entries. Vertices are grouped together by material and vertex material entries tell the client how many vertices there are using a material.

	meshAnimatedBoneCount := dec.Uint16() // number of entries in meshops. Seems to be used only for animated mob models.
	rawScale := dec.Uint16()
	fmt.Println(rawScale)
	scale := float32(1 / float32(int(1)<<rawScale)) // This allows vertex coordinates to be stored as integral values instead of floating-point values, without losing precision based on mesh size. Vertex values are multiplied by (1 shl `scale`) and stored in the vertex entries. FPSCALE is the internal name.
	// convert scale back to rawscale
	//rawScale = uint16(math.Log2(float64(1 / scale)))

	fmt.Println("sclae:", scale)
	/// Vertices (x, y, z) belonging to this mesh. Each axis should
	/// be multiplied by (1 shl `scale`) for the final vertex position.
	for i := 0; i < int(vertexCount); i++ {
		vert := common.Vertex{}
		vert.Position.X = float32(centerX) + (float32(dec.Int16()) * scale)
		vert.Position.Y = float32(centerY) + (float32(dec.Int16()) * scale)
		vert.Position.Z = float32(centerZ) + (float32(dec.Int16()) * scale)
		mesh.Vertices = append(mesh.Vertices, vert)
	}

	for i := 0; i < int(uvCount); i++ {
		uv := common.Vector2{}
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
		normal := common.Vector3{}
		normal.X = float32(dec.Int8()) / 128
		normal.Y = float32(dec.Int8()) / 128
		normal.Z = float32(dec.Int8()) / 128
		if i < len(mesh.Vertices) {
			mesh.Vertices[i].Normal = normal
		}
	}

	for i := 0; i < int(colorCount); i++ {
		color := common.RGBA{}
		color.R = dec.Uint8()
		color.G = dec.Uint8()
		color.B = dec.Uint8()
		color.A = dec.Uint8()
		if i < len(mesh.Vertices) {
			mesh.Vertices[i].Tint = color
		}
	}

	for i := 0; i < int(triangleCount); i++ {
		triangle := common.Triangle{}
		notSolidFlag := dec.Uint16()
		if notSolidFlag != 0 {
			triangle.Flag = 1
		}
		triangle.Index.X = uint32(dec.Uint16())
		triangle.Index.Y = uint32(dec.Uint16())
		triangle.Index.Z = uint32(dec.Uint16())
		mesh.Triangles = append(mesh.Triangles, triangle)
	}

	for i := 0; i < int(vertexPieceCount); i++ {
		vertexPiece := MeshVertexPiece{}
		vertexPiece.Count = dec.Int16()
		vertexPiece.Index1 = dec.Int16()

		mesh.VertexPieces = append(mesh.VertexPieces, vertexPiece)
	}

	for i := 0; i < int(polygonTextureCount); i++ {
		mesh.TriangleMaterials = append(mesh.TriangleMaterials, MeshTriangleMaterial{
			Count:      dec.Uint16(),
			MaterialID: dec.Uint16(),
		})
	}

	for i := 0; i < int(vertexTextureCount); i++ {
		vertexMat := MeshVertexPiece{}
		vertexMat.Count = dec.Int16()
		vertexMat.Index1 = dec.Int16()
		mesh.VertexPieces = append(mesh.VertexPieces, vertexMat)
	}

	for i := 0; i < int(meshAnimatedBoneCount); i++ {
		mab := MeshAnimatedBone{}
		mab.Position.X = float32(centerX) + (float32(dec.Int16()) * scale)
		mab.Position.Y = float32(centerY) + (float32(dec.Int16()) * scale)
		mab.Position.Z = float32(centerZ) + (float32(dec.Int16()) * scale)
		mesh.AnimatedBones = append(mesh.AnimatedBones, mab)
	}

	if dec.Error() != nil {
		return fmt.Errorf("meshRead: %v", dec.Error())
	}

	log.Debugf("mesh %s %d triangles %d material groups", mesh.Name, triangleCount, polygonTextureCount)
	e.Fragments[fragmentOffset] = mesh
	return nil
}

func (e *WLD) meshWrite(w io.Writer, fragmentOffset int) error {
	return nil
}
