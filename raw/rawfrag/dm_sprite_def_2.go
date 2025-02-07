package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/xackery/encdec"
)

// WldFragDmSpriteDef2 is DmSpriteDef2 in libeq, WldFragDmSpriteDef2 in openzone, DMSPRITEDEF2 in wld, WldFragDmSpriteDef2 in lantern
type WldFragDmSpriteDef2 struct {
	nameRef              int32
	Flags                uint32
	MaterialPaletteRef   uint32
	DMTrackRef           int32 // only used for flags/trees
	Fragment3Ref         int32
	Fragment4Ref         int32 // unknown, usually ref to first texture
	CenterOffset         [3]float32
	Params2              [3]uint32
	BoundingRadius       float32 // AKA max_distance in libeq
	BoundingBoxMin       [3]float32
	BoundingBoxMax       [3]float32
	Scale                uint16
	Vertices             [][3]int16
	UVs                  [][2]float32
	VertexNormals        [][3]int8
	Colors               [][4]uint8
	Faces                []WldFragMeshFaceEntry
	FaceMaterialGroups   [][2]uint16
	SkinAssignmentGroups [][2]int16
	VertexMaterialGroups [][2]int16
	MeshOps              []WldFragMeshOpEntry
}

type WldFragMeshFaceEntry struct {
	Flags uint16
	Index [3]uint16
}

type WldFragMeshSkinAssignmentGroup struct {
	Count  int16
	Index1 int16
}

type WldFragMeshFaceMaterialGroup struct {
	Count      uint16
	MaterialID uint16
}

type WldFragMeshOpEntry struct {
	Index1    uint16
	Index2    uint16
	Offset    float32
	Param1    uint8
	TypeField uint8
}

func (e *WldFragDmSpriteDef2) FragCode() int {
	return FragCodeDmSpriteDef2
}

func (e *WldFragDmSpriteDef2) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	padStart := enc.Pos()
	enc.Int32(e.nameRef)

	enc.Uint32(e.Flags)

	enc.Uint32(e.MaterialPaletteRef)
	enc.Int32(e.DMTrackRef)

	enc.Int32(e.Fragment3Ref)
	enc.Int32(e.Fragment4Ref)

	enc.Float32(e.CenterOffset[0])
	enc.Float32(e.CenterOffset[1])
	enc.Float32(e.CenterOffset[2])

	enc.Uint32(e.Params2[0])
	enc.Uint32(e.Params2[1])
	enc.Uint32(e.Params2[2])

	enc.Float32(e.BoundingRadius)
	enc.Float32(e.BoundingBoxMin[0])
	enc.Float32(e.BoundingBoxMin[1])
	enc.Float32(e.BoundingBoxMin[2])
	enc.Float32(e.BoundingBoxMax[0])
	enc.Float32(e.BoundingBoxMax[1])
	enc.Float32(e.BoundingBoxMax[2])

	enc.Uint16(uint16(len(e.Vertices)))
	enc.Uint16(uint16(len(e.UVs)))
	enc.Uint16(uint16(len(e.VertexNormals)))
	enc.Uint16(uint16(len(e.Colors)))
	enc.Uint16(uint16(len(e.Faces)))
	enc.Uint16(uint16(len(e.SkinAssignmentGroups)))
	enc.Uint16(uint16(len(e.FaceMaterialGroups)))
	enc.Uint16(uint16(len(e.VertexMaterialGroups)))
	enc.Uint16(uint16(len(e.MeshOps)))
	enc.Uint16(e.Scale)

	for _, vertex := range e.Vertices {
		enc.Int16(vertex[0])
		enc.Int16(vertex[1])
		enc.Int16(vertex[2])
	}

	for _, uv := range e.UVs {
		if isNewWorld {
			enc.Float32(float32(uv[0]))
			enc.Float32(float32(uv[1]))
		} else {
			enc.Int16(int16(uv[0]))
			enc.Int16(int16(uv[1]))
		}
	}

	for _, normal := range e.VertexNormals {
		enc.Int8(normal[0])
		enc.Int8(normal[1])
		enc.Int8(normal[2])
	}

	for _, color := range e.Colors {
		for i := 0; i < 4; i++ {
			enc.Uint8(color[i])
		}
	}

	for _, face := range e.Faces {
		enc.Uint16(face.Flags)
		enc.Uint16(face.Index[0])
		enc.Uint16(face.Index[1])
		enc.Uint16(face.Index[2])
	}

	for _, vertexPiece := range e.SkinAssignmentGroups {
		enc.Uint16(uint16(vertexPiece[0]))
		enc.Uint16(uint16(vertexPiece[1]))
	}

	for _, triangleMaterial := range e.FaceMaterialGroups {
		enc.Uint16(triangleMaterial[0])
		enc.Uint16(triangleMaterial[1])
	}

	for _, vertexMaterial := range e.VertexMaterialGroups {
		enc.Uint16(uint16(vertexMaterial[0]))
		enc.Uint16(uint16(vertexMaterial[1]))
	}

	for _, meshOp := range e.MeshOps {
		if meshOp.TypeField != 4 {
			enc.Uint16(meshOp.Index1)
			enc.Uint16(meshOp.Index2)
		} else {
			enc.Float32(meshOp.Offset)
		}
		enc.Uint8(meshOp.Param1)
		enc.Uint8(meshOp.TypeField)
	}
	diff := enc.Pos() - padStart
	paddingSize := (4 - diff%4) % 4
	if paddingSize > 0 {
		enc.Bytes(make([]byte, paddingSize))
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragDmSpriteDef2) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32() // flags, currently unknown, zone meshes are 0x00018003, placeable objects are 0x00014003

	e.MaterialPaletteRef = dec.Uint32()
	e.DMTrackRef = dec.Int32() //used by flags/trees only

	e.Fragment3Ref = dec.Int32() // has polyhedron ref
	e.Fragment4Ref = dec.Int32() // unknown, This usually seems to reference the first [TextureImagesFragment] fragment in the file.

	e.CenterOffset[0] = dec.Float32() // for zone meshes, x coordinate of the center of the mesh
	e.CenterOffset[1] = dec.Float32() // for zone meshes, y coordinate of the center of the mesh
	e.CenterOffset[2] = dec.Float32() // for zone meshes, z coordinate of the center of the mesh

	e.Params2[0] = dec.Uint32() // unknown, usually empty
	e.Params2[1] = dec.Uint32() // unknown, usually empty
	e.Params2[2] = dec.Uint32() // unknown, usually empty

	e.BoundingRadius = dec.Float32()    // Given the values in center, this seems to contain the maximum distance between any vertex and that position. It seems to define a radius from that position within which the mesh lies.
	e.BoundingBoxMin[0] = dec.Float32() // min x, y, and z coords in absolute coords of any vertex in the mesh.
	e.BoundingBoxMin[1] = dec.Float32()
	e.BoundingBoxMin[2] = dec.Float32()
	e.BoundingBoxMax[0] = dec.Float32() // max x, y, and z coords in absolute coords of any vertex in the mesh.
	e.BoundingBoxMax[1] = dec.Float32()
	e.BoundingBoxMax[2] = dec.Float32()

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
	e.Scale = dec.Uint16()

	// convert scale back to rawscale
	//rawScale = uint16(math.Log2(float64(1 / scale)))

	// Vertices (x, y, z) belonging to this mesh. Each axis should
	// be multiplied by (1 shl `scale`) for the final vertex position.
	for i := 0; i < int(vertexCount); i++ {
		e.Vertices = append(e.Vertices, [3]int16{dec.Int16(), dec.Int16(), dec.Int16()})
	}

	for i := 0; i < int(uvCount); i++ {
		if isNewWorld {
			e.UVs = append(e.UVs, [2]float32{dec.Float32(), dec.Float32()})
		} else {
			e.UVs = append(e.UVs, [2]float32{float32(dec.Int16()), float32(dec.Int16())})
		}
	}

	for i := 0; i < int(normalCount); i++ {
		e.VertexNormals = append(e.VertexNormals, [3]int8{dec.Int8(), dec.Int8(), dec.Int8()})
	}

	for i := 0; i < int(colorCount); i++ {
		color := [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()}
		e.Colors = append(e.Colors, color)
	}

	for i := 0; i < int(triangleCount); i++ {
		mte := WldFragMeshFaceEntry{}
		mte.Flags = dec.Uint16()
		mte.Index = [3]uint16{dec.Uint16(), dec.Uint16(), dec.Uint16()}

		e.Faces = append(e.Faces, mte)
	}

	for i := 0; i < int(vertexPieceCount); i++ {
		e.SkinAssignmentGroups = append(e.SkinAssignmentGroups, [2]int16{dec.Int16(), dec.Int16()})
	}

	for i := 0; i < int(triangleMaterialCount); i++ {
		e.FaceMaterialGroups = append(e.FaceMaterialGroups, [2]uint16{dec.Uint16(), dec.Uint16()})
	}

	for i := 0; i < int(vertexMaterialCount); i++ {
		e.VertexMaterialGroups = append(e.VertexMaterialGroups, [2]int16{dec.Int16(), dec.Int16()})
	}

	for i := 0; i < int(meshOpCount); i++ {
		val := dec.Bytes(4)
		entry := WldFragMeshOpEntry{
			Param1:    dec.Uint8(),
			TypeField: dec.Uint8(),
		}
		if entry.TypeField != 4 {
			entry.Index1 = binary.LittleEndian.Uint16(val)
			entry.Index2 = binary.LittleEndian.Uint16(val[2:])
		} else {
			bits := binary.LittleEndian.Uint32(val)
			entry.Offset = math.Float32frombits(bits)
		}

		e.MeshOps = append(e.MeshOps, entry)
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil

}

func (e *WldFragDmSpriteDef2) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragDmSpriteDef2) SetNameRef(nameRef int32) {
	e.nameRef = nameRef
}
