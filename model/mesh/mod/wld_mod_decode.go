package mod

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/tag"
)

// DecodeMesh decodes a mesh
func DecodeMesh(model *common.Model, nameRef *int32, isOldWorld bool, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	tag.New()
	*nameRef = dec.Int32()
	flags := dec.Uint32()

	model.FileType = "mod"
	if flags == 0x00018003 {
		model.FileType = "ter"
	}

	if flags == 0x00014003 {
		model.FileType = "mod"
	}

	dec.Uint32() // materialListRef
	dec.Uint32() // animationRef

	dec.Uint32() // unknown, usually empty
	dec.Uint32() // unknown, This usually seems to reference the first [TextureImagesFragment] fragment in the file.
	// for zone meshes, x coordinate of the center of the mesh
	// for placeable objects, this seems to define where the vertices will lie relative to the object’s local origin
	centerX := dec.Float32()
	centerY := dec.Float32()
	centerZ := dec.Float32()

	_ = dec.Uint32() // unknown, usually empty
	_ = dec.Uint32() // unknown, usually empty
	_ = dec.Uint32() // unknown, usually empty

	dec.Float32() // maxDistance, Given the values in center, this seems to contain the maximum distance between any vertex and that position. It seems to define a radius from that position within which the mesh lies.

	dec.Float32() // min x, y, and z coords in absolute coords of any vertex in the model.
	dec.Float32()
	dec.Float32()

	dec.Float32() // max x, y, and z coords in absolute coords of any vertex in the model.
	dec.Float32()
	dec.Float32()

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
	triangleTextureCount := dec.Uint16() // number of triangle texture entries. faces are grouped together by material and polygon material entries. This tells the client the number of faces using a material.
	vertexTextureCount := dec.Uint16()   // number of vertex material entries. Vertices are grouped together by material and vertex material entries tell the client how many vertices there are using a material.

	meshAnimatedBoneCount := dec.Uint16() // number of entries in meshops. Seems to be used only for animated mob models.
	rawScale := dec.Uint16()
	scale := float32(1 / float32(int(1)<<rawScale)) // This allows vertex coordinates to be stored as integral values instead of floating-point values, without losing precision based on mesh size. Vertex values are multiplied by (1 shl `scale`) and stored in the vertex entries. FPSCALE is the internal name.
	// convert scale back to rawscale
	//rawScale = uint16(math.Log2(float64(1 / scale)))

	tag.Add(0, dec.Pos(), "red", "header")
	/// Vertices (x, y, z) belonging to this model. Each axis should
	/// be multiplied by (1 shl `scale`) for the final vertex position.
	for i := 0; i < int(vertexCount); i++ {
		vert := common.Vertex{}
		vert.Position.X = float32(centerX) + (float32(dec.Int16()) * scale)
		vert.Position.Y = float32(centerY) + (float32(dec.Int16()) * scale)
		vert.Position.Z = float32(centerZ) + (float32(dec.Int16()) * scale)
		model.Vertices = append(model.Vertices, vert)
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "%d vertices", vertexCount)

	for i := 0; i < int(uvCount); i++ {
		uv := common.Vector2{}
		if isOldWorld {
			uv.X = float32(dec.Int16()) / 256
			uv.Y = float32(dec.Int16()) / 256
		} else {
			uv.X = float32(dec.Int32()) / 256
			uv.Y = float32(dec.Int32()) / 256
		}
		model.Vertices[i].Uv = uv
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "%d uvs", uvCount)

	for i := 0; i < int(normalCount); i++ {
		normal := common.Vector3{}
		normal.X = float32(dec.Int8()) / 128
		normal.Y = float32(dec.Int8()) / 128
		normal.Z = float32(dec.Int8()) / 128
		if i < len(model.Vertices) {
			model.Vertices[i].Normal = normal
		}
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "%d normals", normalCount)

	for i := 0; i < int(colorCount); i++ {
		color := common.RGBA{}
		color.R = dec.Uint8()
		color.G = dec.Uint8()
		color.B = dec.Uint8()
		color.A = dec.Uint8()
		if i < len(model.Vertices) {
			model.Vertices[i].Tint = color
		}
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "%d colors", colorCount)

	for i := 0; i < int(triangleCount); i++ {
		triangle := common.Triangle{}
		triangle.Flag = uint32(dec.Uint16())
		//dec.Uint16() // remove it's flags
		//triangle.Flag = uint32(dec.Uint16())
		//triangle.Flag = flags
		//fmt.Println("flags", notSolidFlag)
		//if notSolidFlag != 0 {
		//	triangle.Flag = 1
		//}*/

		triangle.Index.X = uint32(dec.Uint16())
		triangle.Index.Y = uint32(dec.Uint16())
		triangle.Index.Z = uint32(dec.Uint16())
		model.Triangles = append(model.Triangles, triangle)
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "%d triangles", triangleCount)

	for i := 0; i < int(vertexPieceCount); i++ {
		// TODO: fix
		dec.Int16()
		dec.Int16()
		/* vertexPiece := MeshVertexPiece{}
		vertexPiece.Count = dec.Int16()
		vertexPiece.Index1 = dec.Int16()

		model.VertexPieces = append(model.VertexPieces, vertexPiece) */
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "%d vertex pieces", vertexPieceCount)

	triangleCounter := 0
	materials := make(map[uint16]*common.Material)
	for i := 0; i < int(triangleTextureCount); i++ {
		count := dec.Uint16() // count
		materialID := dec.Uint16()
		material, ok := materials[materialID]
		if !ok {
			material = &common.Material{
				Name: fmt.Sprintf("material_%d", materialID),
			}
			materials[materialID] = material
		}
		for j := 0; j < int(count); j++ {

			if triangleCounter >= len(model.Triangles) {
				return fmt.Errorf("triangle counter %d out of bounds with %d triangles", triangleCounter, len(model.Triangles))
			}
			if model.Triangles[triangleCounter].Flag&1 > 0 {
				continue
			}
			model.Triangles[triangleCounter].MaterialName = material.Name
			triangleCounter++
		}
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "%d triangle textures", triangleTextureCount)

	for i := 0; i < int(vertexTextureCount); i++ {
		// TODO: fix
		dec.Int16()
		dec.Int16()
		/**
		vertexMat := MeshVertexPiece{}
		vertexMat.Count = dec.Int16()
		vertexMat.Index1 = dec.Int16()
		model.VertexPieces = append(model.VertexPieces, vertexMat)
		*/
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "%d vertex textures", vertexTextureCount)

	for i := 0; i < int(meshAnimatedBoneCount); i++ {
		bone := common.Bone{}
		bone.Pivot.X = float32(centerX) + (float32(dec.Int16()) * scale)
		bone.Pivot.Y = float32(centerY) + (float32(dec.Int16()) * scale)
		bone.Pivot.Z = float32(centerZ) + (float32(dec.Int16()) * scale)
		model.Bones = append(model.Bones, bone)
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "%d bones", meshAnimatedBoneCount)

	if dec.Error() != nil {
		return fmt.Errorf("meshRead: %v", dec.Error())
	}

	//log.Debugf("model %s %d triangles %d material groups", model.Name, triangleCount, triangleTextureCount)
	return nil
}
