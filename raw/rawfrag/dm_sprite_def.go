package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model"
)

// WldFragDMSpriteDef is DmSpriteDef in libeq, Alternate Mesh in openzone, DMSPRITEDEF in wld, LegacyMesh in lantern
type WldFragDMSpriteDef struct {
	NameRef           int32                          `yaml:"name_ref"`
	Flags             uint32                         `yaml:"flags"`
	Fragment1Maybe    int16                          `yaml:"fragment_1_maybe"`
	MaterialReference uint32                         `yaml:"material_reference"`
	Fragment3         uint32                         `yaml:"fragment_3"`
	CenterPosition    model.Vector3                  `yaml:"center_position"`
	Params2           uint32                         `yaml:"params_2"`
	Something2        uint32                         `yaml:"something_2"`
	Something3        uint32                         `yaml:"something_3"`
	Vertices          []model.Vector3                `yaml:"verticies"`
	TexCoords         []model.Vector3                `yaml:"tex_coords"`
	Normals           []model.Vector3                `yaml:"normals"`
	Colors            []int32                        `yaml:"colors"`
	Polygons          []WldFragDMSpriteSpritePolygon `yaml:"polygons"`
	VertexPieces      []WldFragDMSpriteVertexPiece   `yaml:"vertex_pieces"`
	PostVertexFlag    uint32                         `yaml:"post_vertex_flag"`
	RenderGroups      []WldFragDMSpriteRenderGroup   `yaml:"render_groups"`
	VertexTex         []model.Vector2                `yaml:"vertex_tex"`
	Size6Pieces       []WldFragDMSpriteSize6Entry    `yaml:"size_6_pieces"`
}

type WldFragDMSpriteSpritePolygon struct {
	Flag int16 `yaml:"flag"`
	Unk1 int16 `yaml:"unk_1"`
	Unk2 int16 `yaml:"unk_2"`
	Unk3 int16 `yaml:"unk_3"`
	Unk4 int16 `yaml:"unk_4"`
	I1   int16 `yaml:"i_1"`
	I2   int16 `yaml:"i_2"`
	I3   int16 `yaml:"i_3"`
}

type WldFragDMSpriteVertexPiece struct {
	Count  int16 `yaml:"count"`
	Offset int16 `yaml:"offset"`
}

type WldFragDMSpriteRenderGroup struct {
	PolygonCount int16 `yaml:"polygon_count"`
	MaterialId   int16 `yaml:"material_id"`
}

type WldFragDMSpriteSize6Entry struct {
	Unk1 uint32 `yaml:"unk_1"`
	Unk2 uint32 `yaml:"unk_2"`
	Unk3 uint32 `yaml:"unk_3"`
	Unk4 uint32 `yaml:"unk_4"`
	Unk5 uint32 `yaml:"unk_5"`
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
	enc.Float32(e.CenterPosition.X)
	enc.Float32(e.CenterPosition.Y)
	enc.Float32(e.CenterPosition.Z)
	enc.Uint32(e.Params2)
	enc.Uint32(e.Something2)
	enc.Uint32(e.Something3)
	for _, vertex := range e.Vertices {
		enc.Float32(vertex.X)
		enc.Float32(vertex.Y)
		enc.Float32(vertex.Z)
	}
	for _, texCoord := range e.TexCoords {
		enc.Float32(texCoord.X)
		enc.Float32(texCoord.Y)
		enc.Float32(texCoord.Z)
	}
	for _, normal := range e.Normals {
		enc.Float32(normal.X)
		enc.Float32(normal.Y)
		enc.Float32(normal.Z)
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
		enc.Float32(vertexTex.X)
		enc.Float32(vertexTex.Y)
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
	e.CenterPosition.X = dec.Float32()
	e.CenterPosition.Y = dec.Float32()
	e.CenterPosition.Z = dec.Float32()
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
		v := model.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		e.Vertices = append(e.Vertices, v)
	}

	for i := uint32(0); i < texCoordCount; i++ {
		v := model.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		e.TexCoords = append(e.TexCoords, v)
	}

	for i := uint32(0); i < normalCount; i++ {
		v := model.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		e.Normals = append(e.Normals, v)
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
			v := model.Vector2{}
			v.X = dec.Float32()
			v.Y = dec.Float32()
			e.VertexTex = append(e.VertexTex, v)
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
