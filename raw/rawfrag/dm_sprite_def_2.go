package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragDmSpriteDef2 is DmSpriteDef2 in libeq, WldFragDmSpriteDef2 in openzone, DMSPRITEDEF2 in wld, WldFragDmSpriteDef2 in lantern
type WldFragDmSpriteDef2 struct {
	NameRef int32  `yaml:"name_ref"`
	Flags   uint32 `yaml:"flags"`

	MaterialPaletteRef uint32 `yaml:"material_palette_ref"`
	AnimationRef       int32  `yaml:"animation_ref"`

	Fragment3Ref int32         `yaml:"fragment_3_ref"`
	Fragment4Ref int32         `yaml:"fragment_4_ref"` // unknown, usually ref to first texture
	Center       model.Vector3 `yaml:"center"`
	Params2      model.UIndex3 `yaml:"params_2"`

	MaxDistance float32       `yaml:"max_distance"`
	Min         model.Vector3 `yaml:"min"`
	Max         model.Vector3 `yaml:"max"`
	// vertexCount
	// uvCount
	// normalCount
	// colorCount
	// triangleCount
	// vertexPieceCount
	// triangleMaterialCount
	// vertexMaterialCount
	// meshAnimatedBoneCount
	RawScale          uint16                        `yaml:"raw_scale"`
	MeshopCount       uint16                        `yaml:"meshop_count"`
	Scale             float32                       `yaml:"scale"`
	Vertices          [][3]int16                    `yaml:"vertices"`
	UVs               [][2]int16                    `yaml:"uvs"`
	Normals           [][3]int8                     `yaml:"normals"`
	Colors            []model.RGBA                  `yaml:"colors"`
	Triangles         []WldFragMeshTriangleEntry    `yaml:"triangles"`
	TriangleMaterials []WldFragMeshTriangleMaterial `yaml:"triangle_materials"`
	VertexPieces      []WldFragMeshVertexPiece      `yaml:"vertex_pieces"`
	VertexMaterials   []WldFragMeshVertexPiece      `yaml:"vertex_materials"`
	MeshOps           []WldFragMeshOpEntry          `yaml:"mesh_ops"`
}

type WldFragMeshTriangleEntry struct {
	Flags uint16    `yaml:"flags"`
	Index [3]uint16 `yaml:"indexes"`
}

type WldFragMeshVertexPiece struct {
	Count  int16 `yaml:"count"`
	Index1 int16 `yaml:"index_1"`
}

type WldFragMeshTriangleMaterial struct {
	Count      uint16 `yaml:"count"`
	MaterialID uint16 `yaml:"material_id"`
}

type WldFragMeshOpEntry struct {
	Index1    uint16  `yaml:"index_1"`
	Index2    uint16  `yaml:"index_2"`
	Offset    float32 `yaml:"offset"`
	Param1    uint8   `yaml:"param_1"`
	TypeField uint8   `yaml:"type_field"`
}

func (e *WldFragDmSpriteDef2) FragCode() int {
	return FragCodeDmSpriteDef2
}

func (e *WldFragDmSpriteDef2) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)

	enc.Uint32(e.MaterialPaletteRef)
	enc.Int32(e.AnimationRef)

	enc.Int32(e.Fragment3Ref)
	enc.Int32(e.Fragment4Ref)

	enc.Float32(e.Center.X)
	enc.Float32(e.Center.Y)
	enc.Float32(e.Center.Z)

	enc.Uint32(e.Params2.X)
	enc.Uint32(e.Params2.Y)
	enc.Uint32(e.Params2.Z)

	enc.Float32(e.MaxDistance)
	enc.Float32(e.Min.X)
	enc.Float32(e.Min.Y)
	enc.Float32(e.Min.Z)
	enc.Float32(e.Max.X)
	enc.Float32(e.Max.Y)
	enc.Float32(e.Max.Z)

	enc.Uint16(uint16(len(e.Vertices)))
	enc.Uint16(uint16(len(e.UVs)))
	enc.Uint16(uint16(len(e.Normals)))
	enc.Uint16(uint16(len(e.Colors)))
	enc.Uint16(uint16(len(e.Triangles)))
	enc.Uint16(uint16(len(e.VertexPieces)))
	enc.Uint16(uint16(len(e.TriangleMaterials)))
	enc.Uint16(uint16(len(e.VertexMaterials)))
	enc.Uint16(uint16(len(e.MeshOps)))
	enc.Uint16(e.RawScale)

	for _, vertex := range e.Vertices {
		enc.Int16(vertex[0])
		enc.Int16(vertex[1])
		enc.Int16(vertex[2])
	}

	for _, uv := range e.UVs {
		enc.Int16(uv[0])
		enc.Int16(uv[1])
	}

	for _, normal := range e.Normals {
		enc.Int8(normal[0])
		enc.Int8(normal[1])
		enc.Int8(normal[2])
	}

	for _, color := range e.Colors {
		enc.Uint8(color.R)
		enc.Uint8(color.G)
		enc.Uint8(color.B)
		enc.Uint8(color.A)
	}

	for _, triangle := range e.Triangles {
		enc.Uint16(triangle.Flags)
		enc.Uint16(triangle.Index[0])
		enc.Uint16(triangle.Index[1])
		enc.Uint16(triangle.Index[2])
	}

	for _, vertexPiece := range e.VertexPieces {
		enc.Uint16(uint16(vertexPiece.Count))
		enc.Uint16(uint16(vertexPiece.Index1))
	}

	for _, triangleMaterial := range e.TriangleMaterials {
		enc.Uint16(triangleMaterial.Count)
		enc.Uint16(triangleMaterial.MaterialID)
	}

	for _, vertexMaterial := range e.VertexMaterials {
		enc.Uint16(uint16(vertexMaterial.Count))
		enc.Uint16(uint16(vertexMaterial.Index1))
	}

	start := enc.Pos()
	for _, meshOp := range e.MeshOps {
		enc.Uint16(meshOp.Index1)
		enc.Uint16(meshOp.Index2)
		enc.Float32(meshOp.Offset)
		enc.Uint8(meshOp.Param1)
		enc.Uint8(meshOp.TypeField)
	}
	diff := enc.Pos() - start
	paddingSize := (4 - diff%4) % 4

	enc.Bytes(make([]byte, paddingSize))

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragDmSpriteDef2) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32() // flags, currently unknown, zone meshes are 0x00018003, placeable objects are 0x00014003

	e.MaterialPaletteRef = dec.Uint32()
	e.AnimationRef = dec.Int32() //used by flags/trees only

	e.Fragment3Ref = dec.Int32() // unknown, usually empty
	e.Fragment4Ref = dec.Int32() // unknown, This usually seems to reference the first [TextureImagesFragment] fragment in the file.

	e.Center.X = dec.Float32() // for zone meshes, x coordinate of the center of the mesh
	e.Center.Y = dec.Float32() // for zone meshes, x coordinate of the center of the mesh
	e.Center.Z = dec.Float32()

	e.Params2.X = dec.Uint32() // unknown, usually empty
	e.Params2.Y = dec.Uint32() // unknown, usually empty
	e.Params2.Z = dec.Uint32() // unknown, usually empty

	e.MaxDistance = dec.Float32() // Given the values in center, this seems to contain the maximum distance between any vertex and that position. It seems to define a radius from that position within which the mesh lies.
	e.Min.X = dec.Float32()       // min x, y, and z coords in absolute coords of any vertex in the mesh.
	e.Min.Y = dec.Float32()
	e.Min.Z = dec.Float32()
	e.Max.X = dec.Float32() // max x, y, and z coords in absolute coords of any vertex in the mesh.
	e.Max.Y = dec.Float32()
	e.Max.Z = dec.Float32()

	vertexCount := dec.Uint16()   // number of vertices in the mesh (called position_count in libeq)
	uvCount := dec.Uint16()       // number of uv in the mesh (called texture_coordinate_count in libeq)
	normalCount := dec.Uint16()   // number of vertex normal entries in the mesh (called normal_count in libeq)
	colorCount := dec.Uint16()    // number of vertex color entries in the mesh (called color_count in libeq)
	triangleCount := dec.Uint16() // number of triangles in the mesh (called face_count in libeq)
	// This seems to only be used when dealing with animated (mob) models.
	// It contains the number of vertex piece entries. Vertices are grouped together by
	// skeleton piece in this case and vertex piece entries tell the client how
	// many vertices are in each piece. It’s possible that there could be more
	// pieces in the skeleton than are in the meshes it references. Extra pieces have
	// no faces or vertices and I suspect they are there to define attachment points for
	// objects (e.g. weapons or shields).
	vertexPieceCount := dec.Uint16()
	triangleMaterialCount := dec.Uint16() // number of triangle texture entries. faces are grouped together by material and polygon material entries. This tells the client the number of faces using a material.
	vertexMaterialCount := dec.Uint16()   // number of vertex material entries. Vertices are grouped together by material and vertex material entries tell the client how many vertices there are using a material.

	meshOpCount := dec.Uint16() // number of entries in meshops. Seems to be used only for animated mob models.
	e.RawScale = dec.Uint16()

	// convert scale back to rawscale
	//rawScale = uint16(math.Log2(float64(1 / scale)))

	// Vertices (x, y, z) belonging to this mesh. Each axis should
	// be multiplied by (1 shl `scale`) for the final vertex position.
	for i := 0; i < int(vertexCount); i++ {
		e.Vertices = append(e.Vertices, [3]int16{dec.Int16(), dec.Int16(), dec.Int16()})
	}

	for i := 0; i < int(uvCount); i++ {
		e.UVs = append(e.UVs, [2]int16{dec.Int16(), dec.Int16()})
	}

	for i := 0; i < int(normalCount); i++ {
		e.Normals = append(e.Normals, [3]int8{dec.Int8(), dec.Int8(), dec.Int8()})
	}

	for i := 0; i < int(colorCount); i++ {
		color := model.RGBA{
			R: dec.Uint8(),
			G: dec.Uint8(),
			B: dec.Uint8(),
			A: dec.Uint8(),
		}
		e.Colors = append(e.Colors, color)
	}

	for i := 0; i < int(triangleCount); i++ {
		mte := WldFragMeshTriangleEntry{}
		mte.Flags = dec.Uint16()
		mte.Index = [3]uint16{dec.Uint16(), dec.Uint16(), dec.Uint16()}

		e.Triangles = append(e.Triangles, mte)
	}

	for i := 0; i < int(vertexPieceCount); i++ {
		vertexPiece := WldFragMeshVertexPiece{}
		vertexPiece.Count = dec.Int16()
		vertexPiece.Index1 = dec.Int16()

		e.VertexPieces = append(e.VertexPieces, vertexPiece)
	}

	for i := 0; i < int(triangleMaterialCount); i++ {
		e.TriangleMaterials = append(e.TriangleMaterials, WldFragMeshTriangleMaterial{
			Count:      dec.Uint16(),
			MaterialID: dec.Uint16(),
		})
	}

	for i := 0; i < int(vertexMaterialCount); i++ {
		vertexMat := WldFragMeshVertexPiece{}
		vertexMat.Count = dec.Int16()
		vertexMat.Index1 = dec.Int16()
		e.VertexMaterials = append(e.VertexMaterials, vertexMat)
	}

	for i := 0; i < int(meshOpCount); i++ {
		e.MeshOps = append(e.MeshOps, WldFragMeshOpEntry{
			Index1:    dec.Uint16(),
			Index2:    dec.Uint16(),
			Offset:    dec.Float32(),
			Param1:    dec.Uint8(),
			TypeField: dec.Uint8(),
		})
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil

}
