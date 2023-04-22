package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

type dmSpriteDef struct {
	nameRef           int16
	flags             uint32
	fragment1Maybe    int16
	materialReference uint32
	fragment3         uint32
	centerPosition    geo.Vector3
	params2           uint32
	something2        uint32
	something3        uint32
	verticies         []geo.Vector3
	texCoords         []geo.Vector3
	normals           []geo.Vector3
	colors            []int32
	polygons          []spritePolygon
	vertexPieces      []spriteVertexPiece
	renderGroups      []spriteRenderGroup
	vertexTex         []geo.Vector2
}

type spritePolygon struct {
	flag int16
	unk1 int16
	unk2 int16
	unk3 int16
	unk4 int16
	i1   int16
	i2   int16
	i3   int16
}

type spriteVertexPiece struct {
	count  int16
	offset int16
}

type spriteRenderGroup struct {
	polygonCount int16
	materialID   int16
}

func (e *WLD) dmSpriteDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &dmSpriteDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int16()
	def.flags = dec.Uint32()
	vertexCount := dec.Uint32()
	texCoordCount := dec.Uint32()
	normalCount := dec.Uint32()
	colorCount := dec.Uint32()
	polygonCount := dec.Uint32()
	size6 := dec.Int16()
	def.fragment1Maybe = dec.Int16()
	vertexPieceCount := dec.Uint32()
	def.materialReference = dec.Uint32()
	def.fragment3 = dec.Uint32()
	def.centerPosition.X = dec.Float32()
	def.centerPosition.Y = dec.Float32()
	def.centerPosition.Z = dec.Float32()
	def.params2 = dec.Uint32()
	def.something2 = dec.Uint32()
	def.something3 = dec.Uint32()

	// TODO: fix alignments
	if vertexCount > 999 {
		log.Debugf("this is misaligned and needs work")
		return nil
	}

	for i := uint32(0); i < vertexCount; i++ {
		var vertex geo.Vector3
		vertex.X = dec.Float32()
		vertex.Y = dec.Float32()
		vertex.Z = dec.Float32()
		def.verticies = append(def.verticies, vertex)
	}

	for i := uint32(0); i < texCoordCount; i++ {
		var texCoord geo.Vector3
		texCoord.X = dec.Float32()
		texCoord.Y = dec.Float32()
		texCoord.Z = dec.Float32()
		def.texCoords = append(def.texCoords, texCoord)
	}

	for i := uint32(0); i < normalCount; i++ {
		var normal geo.Vector3
		normal.X = dec.Float32()
		normal.Y = dec.Float32()
		normal.Z = dec.Float32()
		def.normals = append(def.normals, normal)
	}

	for i := uint32(0); i < colorCount; i++ {
		def.colors = append(def.colors, dec.Int32())
	}

	for i := uint32(0); i < polygonCount; i++ {
		var polygon spritePolygon
		polygon.flag = dec.Int16()
		polygon.unk1 = dec.Int16()
		polygon.unk2 = dec.Int16()
		polygon.unk3 = dec.Int16()
		polygon.unk4 = dec.Int16()
		polygon.i1 = dec.Int16()
		polygon.i2 = dec.Int16()
		polygon.i3 = dec.Int16()
		def.polygons = append(def.polygons, polygon)
	}

	for i := uint32(0); i < uint32(size6); i++ {
		dec.Uint32()
		dec.Uint32()
		dec.Uint32()
		dec.Uint32()
		dec.Uint32()
	}

	for i := uint32(0); i < vertexPieceCount; i++ {
		var vertexPiece spriteVertexPiece
		vertexPiece.count = dec.Int16()
		vertexPiece.offset = dec.Int16()
		def.vertexPieces = append(def.vertexPieces, vertexPiece)
	}

	if def.flags&9 != 0 {
		dec.Uint32()
	}

	if def.flags&11 != 0 {
		spriteRenderGroupCount := dec.Uint32()
		for i := uint32(0); i < spriteRenderGroupCount; i++ {
			var renderGroup spriteRenderGroup
			renderGroup.polygonCount = dec.Int16()
			renderGroup.materialID = dec.Int16()
			def.renderGroups = append(def.renderGroups, renderGroup)
		}
	}

	if def.flags&12 != 0 {
		spriteVertexCount := dec.Uint32()
		for i := uint32(0); i < spriteVertexCount; i++ {
			var vertexTex geo.Vector2
			vertexTex.X = dec.Float32()
			vertexTex.Y = dec.Float32()
			def.vertexTex = append(def.vertexTex, vertexTex)
		}
	}

	if def.flags&13 != 0 {
		dec.Uint32()
		dec.Uint32()
		dec.Uint32()
	}

	if dec.Error() != nil {
		return fmt.Errorf("dmSpriteDefRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *dmSpriteDef) build(e *WLD) error {
	return nil
}
