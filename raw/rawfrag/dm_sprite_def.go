package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// WldFragDMSpriteDef is DmSpriteDef in libeq, Alternate Mesh in openzone, DMSPRITEDEF in wld, LegacyMesh in lantern
type WldFragDMSpriteDef struct {
	NameRef           int32
	Flags             uint32
	Fragment1Maybe    int16
	MaterialReference uint32
	Fragment3         uint32
	CenterPosition    [3]float32
	Params2           uint32
	Something2        uint32
	Something3        uint32
	Vertices          [][3]float32
	TexCoords         [][3]float32
	Normals           [][3]float32
	Colors            []int32
	Polygons          []WldFragDMSpriteSpritePolygon
	VertexPieces      []WldFragDMSpriteVertexPiece
	PostVertexFlag    uint32
	RenderGroups      []WldFragDMSpriteRenderGroup
	VertexTex         [][2]float32
	Size6Pieces       []WldFragDMSpriteSize6Entry
}

type WldFragDMSpriteSpritePolygon struct {
	Flag int16
	Unk1 int16
	Unk2 int16
	Unk3 int16
	Unk4 int16
	I1   int16
	I2   int16
	I3   int16
}

type WldFragDMSpriteVertexPiece struct {
	Count  int16
	Offset int16
}

type WldFragDMSpriteRenderGroup struct {
	PolygonCount int16
	MaterialId   int16
}

type WldFragDMSpriteSize6Entry struct {
	Unk1 uint32
	Unk2 uint32
	Unk3 uint32
	Unk4 uint32
	Unk5 uint32
}

func (e *WldFragDMSpriteDef) FragCode() int {
	return FragCodeDMSpriteDef
}

func (e *WldFragDMSpriteDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Int16(e.Fragment1Maybe)
	enc.Uint32(e.MaterialReference)
	enc.Uint32(e.Fragment3)
	enc.Float32(e.CenterPosition[0])
	enc.Float32(e.CenterPosition[1])
	enc.Float32(e.CenterPosition[2])
	enc.Uint32(e.Params2)
	enc.Uint32(e.Something2)
	enc.Uint32(e.Something3)
	for _, vertex := range e.Vertices {
		enc.Float32(vertex[0])
		enc.Float32(vertex[1])
		enc.Float32(vertex[2])
	}
	for _, texCoord := range e.TexCoords {
		enc.Float32(texCoord[0])
		enc.Float32(texCoord[1])
		enc.Float32(texCoord[2])
	}
	for _, normal := range e.Normals {
		enc.Float32(normal[0])
		enc.Float32(normal[1])
		enc.Float32(normal[2])
	}
	for _, color := range e.Colors {
		enc.Int32(color)
	}
	for _, polygon := range e.Polygons {
		enc.Int16(polygon.Flag)
		enc.Int16(polygon.Unk1)
		enc.Int16(polygon.Unk2)
		enc.Int16(polygon.Unk3)
		enc.Int16(polygon.Unk4)
		enc.Int16(polygon.I1)
		enc.Int16(polygon.I2)
		enc.Int16(polygon.I3)
	}

	for _, sizePiece := range e.Size6Pieces {
		enc.Uint32(sizePiece.Unk1)
		enc.Uint32(sizePiece.Unk2)
		enc.Uint32(sizePiece.Unk3)
		enc.Uint32(sizePiece.Unk4)
		enc.Uint32(sizePiece.Unk5)
	}

	for _, vertexPiece := range e.VertexPieces {
		enc.Int16(vertexPiece.Count)
		enc.Int16(vertexPiece.Offset)
	}

	if e.Flags&9 != 0 {
		enc.Uint32(e.PostVertexFlag)
	}

	for _, renderGroup := range e.RenderGroups {
		enc.Int16(renderGroup.PolygonCount)
		enc.Int16(renderGroup.MaterialId)
	}
	for _, vertexTex := range e.VertexTex {
		enc.Float32(vertexTex[0])
		enc.Float32(vertexTex[1])
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDMSpriteDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	vertexCount := dec.Int16()
	texCoordCount := dec.Uint32()
	normalCount := dec.Uint32()
	colorCount := dec.Uint32()
	polygonCount := dec.Uint32()
	size6 := dec.Uint32()
	e.Fragment1Maybe = dec.Int16()
	vertexPieceCount := dec.Uint32()
	e.MaterialReference = dec.Uint32()
	e.Fragment3 = dec.Uint32()
	e.CenterPosition[0] = dec.Float32()
	e.CenterPosition[1] = dec.Float32()
	e.CenterPosition[2] = dec.Float32()
	e.Params2 = dec.Uint32()
	e.Something2 = dec.Uint32()
	e.Something3 = dec.Uint32()

	if vertexCount > 999 {
		return fmt.Errorf("vertex count misaligned (%d)", vertexCount)
	}

	if texCoordCount > 999 {
		log.Warnf("texCoordCount > 999: %d", texCoordCount)
		return nil
		//return fmt.Errorf("tex coord count misaligned (%d)", texCoordCount)
	}

	if normalCount > 999 {
		return fmt.Errorf("normal count misaligned (%d)", normalCount)
	}

	for i := int16(0); i < vertexCount; i++ {
		e.Vertices = append(e.Vertices, [3]float32{dec.Float32(), dec.Float32(), dec.Float32()})
	}

	for i := uint32(0); i < texCoordCount; i++ {
		e.TexCoords = append(e.TexCoords, [3]float32{dec.Float32(), dec.Float32(), dec.Float32()})
	}

	for i := uint32(0); i < normalCount; i++ {
		e.Normals = append(e.Normals, [3]float32{dec.Float32(), dec.Float32(), dec.Float32()})
	}

	for i := uint32(0); i < colorCount; i++ {
		e.Colors = append(e.Colors, dec.Int32())
	}

	for i := uint32(0); i < polygonCount; i++ {
		p := WldFragDMSpriteSpritePolygon{}
		p.Flag = dec.Int16()
		p.Unk1 = dec.Int16()
		p.Unk2 = dec.Int16()
		p.Unk3 = dec.Int16()
		p.Unk4 = dec.Int16()
		p.I1 = dec.Int16()
		p.I2 = dec.Int16()
		p.I3 = dec.Int16()
		e.Polygons = append(e.Polygons, p)
	}

	for i := uint32(0); i < size6; i++ {
		s := WldFragDMSpriteSize6Entry{}
		s.Unk1 = dec.Uint32()
		s.Unk2 = dec.Uint32()
		s.Unk3 = dec.Uint32()
		s.Unk4 = dec.Uint32()
		s.Unk5 = dec.Uint32()
		e.Size6Pieces = append(e.Size6Pieces, s)
	}

	for i := uint32(0); i < vertexPieceCount; i++ {
		v := WldFragDMSpriteVertexPiece{}
		v.Count = dec.Int16()
		v.Offset = dec.Int16()
		e.VertexPieces = append(e.VertexPieces, v)
	}

	if e.Flags&9 != 0 {
		e.PostVertexFlag = dec.Uint32()
	}

	if e.Flags&11 != 0 {
		spriteRenderGroupCount := dec.Uint32()
		for i := uint32(0); i < spriteRenderGroupCount; i++ {
			s := WldFragDMSpriteRenderGroup{}
			s.PolygonCount = dec.Int16()
			s.MaterialId = dec.Int16()
			e.RenderGroups = append(e.RenderGroups, s)
		}
	}

	if e.Flags&12 != 0 {
		spriteVertexCount := dec.Uint32()
		for i := uint32(0); i < spriteVertexCount; i++ {
			e.VertexTex = append(e.VertexTex, [2]float32{dec.Float32(), dec.Float32()})
		}
	}

	if e.Flags&13 != 0 {
		dec.Uint32()
		dec.Uint32()
		dec.Uint32()
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil

}
