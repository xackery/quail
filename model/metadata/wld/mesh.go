package wld

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// TwoDSprite is Sprite2DDef in libeq, Two-Dimensional Object in openzone, 2DSPRITEDEF in wld, Fragment06 in lantern
type TwoDSprite struct {
	NameRef                     int32
	Flags                       uint32
	TextureCount                uint32
	PitchCount                  uint32
	Scale                       common.Vector2
	SphereRef                   uint32
	DepthScale                  float32
	CenterOffset                common.Vector3
	BoundingRadius              float32
	CurrentFrameRef             int32
	Sleep                       uint32
	Headings                    []uint32
	RenderMethod                uint32
	RenderFlags                 uint32
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          common.Vector3
	RenderUVInfoUAxis           common.Vector3
	RenderUVInfoVAxis           common.Vector3
	RenderUVMapEntryCount       uint32
	RenderUVMapEntries          []common.Vector2
}

func (e *TwoDSprite) FragCode() int {
	return 0x06
}

func (e *TwoDSprite) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.TextureCount)
	enc.Uint32(e.PitchCount)
	enc.Float32(e.Scale.X)
	enc.Float32(e.Scale.Y)
	enc.Uint32(e.SphereRef)
	if e.Flags&0x80 == 0x80 {
		enc.Float32(e.DepthScale)
	}
	if e.Flags&0x01 == 0x01 {
		enc.Float32(e.CenterOffset.X)
		enc.Float32(e.CenterOffset.Y)
		enc.Float32(e.CenterOffset.Z)
	}
	if e.Flags&0x02 == 0x02 {
		enc.Float32(e.BoundingRadius)
	}
	if e.Flags&0x04 == 0x04 {
		enc.Int32(e.CurrentFrameRef)
	}
	if e.Flags&0x08 == 0x08 {
		enc.Uint32(e.Sleep)
	}
	for _, heading := range e.Headings {
		enc.Uint32(heading)
	}
	if e.Flags&0x10 == 0x10 {
		enc.Uint32(e.RenderMethod)
		enc.Uint32(e.RenderFlags)
		if e.RenderFlags&0x01 == 0x01 {
			enc.Uint32(e.RenderPen)
		}
		if e.RenderFlags&0x02 == 0x02 {
			enc.Float32(e.RenderBrightness)
		}
		if e.RenderFlags&0x04 == 0x04 {
			enc.Float32(e.RenderScaledAmbient)
		}
		if e.RenderFlags&0x08 == 0x08 {
			enc.Uint32(e.RenderSimpleSpriteReference)
		}
		if e.RenderFlags&0x10 == 0x10 {
			enc.Float32(e.RenderUVInfoOrigin.X)
			enc.Float32(e.RenderUVInfoOrigin.Y)
			enc.Float32(e.RenderUVInfoOrigin.Z)
			enc.Float32(e.RenderUVInfoUAxis.X)
			enc.Float32(e.RenderUVInfoUAxis.Y)
			enc.Float32(e.RenderUVInfoUAxis.Z)
			enc.Float32(e.RenderUVInfoVAxis.X)
			enc.Float32(e.RenderUVInfoVAxis.Y)
			enc.Float32(e.RenderUVInfoVAxis.Z)
		}
		if e.RenderFlags&0x20 == 0x20 {
			enc.Uint32(e.RenderUVMapEntryCount)
			for _, entry := range e.RenderUVMapEntries {
				enc.Float32(entry.X)
				enc.Float32(entry.Y)
			}
		}
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil

}

func decodeTwoDSprite(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &TwoDSprite{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.TextureCount = dec.Uint32()
	d.PitchCount = dec.Uint32()
	d.Scale.X = dec.Float32()
	d.Scale.Y = dec.Float32()
	d.SphereRef = dec.Uint32()
	if d.Flags&0x80 == 0x80 {
		d.DepthScale = dec.Float32()
	}
	if d.Flags&0x01 == 0x01 {
		d.CenterOffset.X = dec.Float32()
		d.CenterOffset.Y = dec.Float32()
		d.CenterOffset.Z = dec.Float32()
	}
	if d.Flags&0x02 == 0x02 {
		d.BoundingRadius = dec.Float32()
	}
	if d.Flags&0x04 == 0x04 {
		d.CurrentFrameRef = dec.Int32()
	}
	if d.Flags&0x08 == 0x08 {
		d.Sleep = dec.Uint32()
	}
	d.Headings = make([]uint32, d.PitchCount)
	for i := uint32(0); i < d.PitchCount; i++ {
		d.Headings[i] = dec.Uint32()
	}
	if d.Flags&0x10 == 0x10 {
		d.RenderMethod = dec.Uint32()
		d.RenderFlags = dec.Uint32()
		if d.RenderFlags&0x01 == 0x01 {
			d.RenderPen = dec.Uint32()
		}
		if d.RenderFlags&0x02 == 0x02 {
			d.RenderBrightness = dec.Float32()
		}
		if d.RenderFlags&0x04 == 0x04 {
			d.RenderScaledAmbient = dec.Float32()
		}
		if d.RenderFlags&0x08 == 0x08 {
			d.RenderSimpleSpriteReference = dec.Uint32()
		}
		if d.RenderFlags&0x10 == 0x10 {
			d.RenderUVInfoOrigin.X = dec.Float32()
			d.RenderUVInfoOrigin.Y = dec.Float32()
			d.RenderUVInfoOrigin.Z = dec.Float32()
			d.RenderUVInfoUAxis.X = dec.Float32()
			d.RenderUVInfoUAxis.Y = dec.Float32()
			d.RenderUVInfoUAxis.Z = dec.Float32()
			d.RenderUVInfoVAxis.X = dec.Float32()
			d.RenderUVInfoVAxis.Y = dec.Float32()
			d.RenderUVInfoVAxis.Z = dec.Float32()
		}
		if d.RenderFlags&0x20 == 0x20 {
			d.RenderUVMapEntryCount = dec.Uint32()
			for i := uint32(0); i < d.RenderUVMapEntryCount; i++ {
				v := common.Vector2{}
				v.X = dec.Float32()
				v.Y = dec.Float32()
				d.RenderUVMapEntries = append(d.RenderUVMapEntries, v)
			}
		}
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}

	return d, nil
}

// TwoDSpriteRef is Sprite2D in libeq, Two-Dimensional Object Reference in openzone, 2DSPRITE (ref) in wld, Fragment07 in lantern
type TwoDSpriteRef struct {
	NameRef       int32
	TwoDSpriteRef uint32
	Flags         uint32
}

func (e *TwoDSpriteRef) FragCode() int {
	return 0x07
}

func (e *TwoDSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.TwoDSpriteRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeTwoDSpriteRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &TwoDSpriteRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.TwoDSpriteRef = dec.Uint32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// ThreeDSprite is Sprite3DDef in libeq, Camera in openzone, 3DSPRITEDEF in wld, Camera in lantern
type ThreeDSprite struct {
	NameRef       int32
	Flags         uint32
	VertexCount   uint32
	BspNodeCount  uint32
	SphereListRef uint32
	CenterOffset  common.Vector3
	Radius        float32
	Vertices      []common.Vector3
	BspNodes      []BspNode
}

type BspNode struct {
	VertexCount                 uint32
	FrontTree                   uint32
	BackTree                    uint32
	VertexIndexes               []uint32
	RenderMethod                uint32
	RenderFlags                 uint32
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          common.Vector3
	RenderUVInfoUAxis           common.Vector3
	RenderUVInfoVAxis           common.Vector3
	RenderUVMapEntryCount       uint32
	RenderUVMapEntries          []common.Vector2
}

func (e *ThreeDSprite) FragCode() int {
	return 0x08
}

func (e *ThreeDSprite) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.VertexCount)
	enc.Uint32(e.BspNodeCount)
	enc.Uint32(e.SphereListRef)
	enc.Float32(e.CenterOffset.X)
	enc.Float32(e.CenterOffset.Y)
	enc.Float32(e.CenterOffset.Z)
	enc.Float32(e.Radius)
	for _, vertex := range e.Vertices {
		enc.Float32(vertex.X)
		enc.Float32(vertex.Y)
		enc.Float32(vertex.Z)
	}

	for _, bspNode := range e.BspNodes {
		enc.Uint32(bspNode.VertexCount)
		enc.Uint32(bspNode.FrontTree)
		enc.Uint32(bspNode.BackTree)
		for _, vertexIndex := range bspNode.VertexIndexes {
			enc.Uint32(vertexIndex)
		}
		enc.Uint32(bspNode.RenderMethod)
		enc.Uint32(bspNode.RenderFlags)
		if bspNode.RenderFlags&0x01 == 0x01 {
			enc.Uint32(bspNode.RenderPen)
		}
		if bspNode.RenderFlags&0x02 == 0x02 {
			enc.Float32(bspNode.RenderBrightness)
		}
		if bspNode.RenderFlags&0x04 == 0x04 {
			enc.Float32(bspNode.RenderScaledAmbient)
		}
		if bspNode.RenderFlags&0x08 == 0x08 {
			enc.Uint32(bspNode.RenderSimpleSpriteReference)
		}
		if bspNode.RenderFlags&0x10 == 0x10 {
			enc.Float32(bspNode.RenderUVInfoOrigin.X)
			enc.Float32(bspNode.RenderUVInfoOrigin.Y)
			enc.Float32(bspNode.RenderUVInfoOrigin.Z)
			enc.Float32(bspNode.RenderUVInfoUAxis.X)
			enc.Float32(bspNode.RenderUVInfoUAxis.Y)
			enc.Float32(bspNode.RenderUVInfoUAxis.Z)
			enc.Float32(bspNode.RenderUVInfoVAxis.X)
			enc.Float32(bspNode.RenderUVInfoVAxis.Y)
			enc.Float32(bspNode.RenderUVInfoVAxis.Z)
		}
		if bspNode.RenderFlags&0x20 == 0x20 {
			enc.Uint32(bspNode.RenderUVMapEntryCount)
			for _, entry := range bspNode.RenderUVMapEntries {
				enc.Float32(entry.X)
				enc.Float32(entry.Y)
			}

		}
	}
	if enc.Error() != nil {
		return enc.Error()
	}

	return nil
}

func decodeThreeDSprite(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &ThreeDSprite{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.VertexCount = dec.Uint32()
	d.BspNodeCount = dec.Uint32()
	d.SphereListRef = dec.Uint32()
	d.CenterOffset.X = dec.Float32()
	d.CenterOffset.Y = dec.Float32()
	d.CenterOffset.Z = dec.Float32()
	d.Radius = dec.Float32()
	for i := 0; i < int(d.VertexCount); i++ {
		v := common.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		d.Vertices = append(d.Vertices, v)
	}
	for i := 0; i < int(d.BspNodeCount); i++ {
		b := BspNode{}
		b.VertexCount = dec.Uint32()
		b.FrontTree = dec.Uint32()
		b.BackTree = dec.Uint32()
		for j := 0; j < int(b.VertexCount); j++ {
			b.VertexIndexes = append(b.VertexIndexes, dec.Uint32())
		}
		b.RenderMethod = dec.Uint32()
		b.RenderFlags = dec.Uint32()
		if b.RenderFlags&0x01 == 0x01 {
			b.RenderPen = dec.Uint32()
		}
		if b.RenderFlags&0x02 == 0x02 {
			b.RenderBrightness = dec.Float32()
		}
		if b.RenderFlags&0x04 == 0x04 {
			b.RenderScaledAmbient = dec.Float32()
		}
		if b.RenderFlags&0x08 == 0x08 {
			b.RenderSimpleSpriteReference = dec.Uint32()
		}
		if b.RenderFlags&0x10 == 0x10 {
			b.RenderUVInfoOrigin.X = dec.Float32()
			b.RenderUVInfoOrigin.Y = dec.Float32()
			b.RenderUVInfoOrigin.Z = dec.Float32()
			b.RenderUVInfoUAxis.X = dec.Float32()
			b.RenderUVInfoUAxis.Y = dec.Float32()
			b.RenderUVInfoUAxis.Z = dec.Float32()
			b.RenderUVInfoVAxis.X = dec.Float32()
			b.RenderUVInfoVAxis.Y = dec.Float32()
			b.RenderUVInfoVAxis.Z = dec.Float32()
		}
		if b.RenderFlags&0x20 == 0x20 {
			b.RenderUVMapEntryCount = dec.Uint32()
			for j := 0; j < int(b.RenderUVMapEntryCount); j++ {
				v := common.Vector2{}
				v.X = dec.Float32()
				v.Y = dec.Float32()
				b.RenderUVMapEntries = append(b.RenderUVMapEntries, v)
			}
		}
		d.BspNodes = append(d.BspNodes, b)
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}

	return d, nil
}

// ThreeDSpriteRef is Sprite3D in libeq, Camera Reference in openzone, 3DSPRITE (ref) in wld, CameraReference in lantern
type ThreeDSpriteRef struct {
	NameRef   int32
	ThreeDRef int32
	Flags     uint32
}

func (e *ThreeDSpriteRef) FragCode() int {
	return 0x09
}

func (e *ThreeDSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.ThreeDRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeThreeDSpriteRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &ThreeDSpriteRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.ThreeDRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// FourDSprite is Sprite4DDef in libeq, empty in openzone, 4DSPRITEDEF in wld
type FourDSprite struct {
	NameRef         int32
	Flags           uint32
	FrameCount      uint32
	PolyRef         int32
	CenterOffset    common.Vector3
	Radius          float32
	CurrentFrame    uint32
	Sleep           uint32
	SpriteFragments []uint32
}

func (e *FourDSprite) FragCode() int {
	return 0x0A
}

func (e *FourDSprite) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.FrameCount)
	enc.Int32(e.PolyRef)
	if e.Flags&0x01 != 0 {
		enc.Float32(e.CenterOffset.X)
		enc.Float32(e.CenterOffset.Y)
		enc.Float32(e.CenterOffset.Z)
	}
	if e.Flags&0x02 != 0 {
		enc.Float32(e.Radius)
	}
	if e.Flags&0x04 != 0 {
		enc.Uint32(e.CurrentFrame)
	}
	if e.Flags&0x08 != 0 {
		enc.Uint32(e.Sleep)
	}
	if e.Flags&0x10 != 0 {
		for _, spriteFragment := range e.SpriteFragments {
			enc.Uint32(spriteFragment)
		}
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeFourDSprite(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &FourDSprite{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.FrameCount = dec.Uint32()
	d.PolyRef = dec.Int32()
	if d.Flags&0x01 != 0 {
		d.CenterOffset.X = dec.Float32()
		d.CenterOffset.Y = dec.Float32()
		d.CenterOffset.Z = dec.Float32()
	}
	if d.Flags&0x02 != 0 {
		d.Radius = dec.Float32()
	}
	if d.Flags&0x04 != 0 {
		d.CurrentFrame = dec.Uint32()
	}
	if d.Flags&0x08 != 0 {
		d.Sleep = dec.Uint32()
	}
	if d.Flags&0x10 != 0 {
		for i := uint32(0); i < d.FrameCount; i++ {
			d.SpriteFragments = append(d.SpriteFragments, dec.Uint32())
		}
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// FourDSpriteRef is Sprite4D in libeq, empty in openzone, 4DSPRITE (ref) in wld
type FourDSpriteRef struct {
	NameRef  int32
	FourDRef int32
	Params1  uint32
}

func (e *FourDSpriteRef) FragCode() int {
	return 0x0B
}

func (e *FourDSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.FourDRef)
	enc.Uint32(e.Params1)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeFourDSpriteRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &FourDSpriteRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.FourDRef = dec.Int32()
	d.Params1 = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// Polyhedron is PolyhedronDef in libeq, Polygon animation in openzone, POLYHEDRONDEFINITION in wld, Fragment17 in lantern
type Polyhedron struct {
	NameRef  int32
	Flags    uint32
	Size1    uint32
	Size2    uint32
	Params1  float32
	Params2  float32
	Entries1 []common.Vector3
	Entries2 []Entries2
}

type Entries2 struct {
	Unk1 uint32
	Unk2 []uint32
}

func (e *Polyhedron) FragCode() int {
	return 0x17
}

func (e *Polyhedron) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.Size1)
	enc.Uint32(e.Size2)
	enc.Float32(e.Params1)
	enc.Float32(e.Params2)
	for _, entry := range e.Entries1 {
		enc.Float32(entry.X)
		enc.Float32(entry.Y)
		enc.Float32(entry.Z)
	}
	for _, entry := range e.Entries2 {
		enc.Uint32(entry.Unk1)
		for _, unk2 := range entry.Unk2 {
			enc.Uint32(unk2)
		}
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodePolyhedron(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Polyhedron{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.Size1 = dec.Uint32()
	d.Size2 = dec.Uint32()
	d.Params1 = dec.Float32()
	d.Params2 = dec.Float32()
	for i := uint32(0); i < d.Size1; i++ {
		v := common.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		d.Entries1 = append(d.Entries1, v)
	}
	for i := uint32(0); i < d.Size2; i++ {
		e := Entries2{}
		e.Unk1 = dec.Uint32()
		for j := uint32(0); j < d.Size1; j++ {
			e.Unk2 = append(e.Unk2, dec.Uint32())
		}
		d.Entries2 = append(d.Entries2, e)
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// PolyhedronRef is Polyhedron in libeq, Polygon Animation Reference in openzone, POLYHEDRON (ref) in wld, Fragment18 in lantern
type PolyhedronRef struct {
	NameRef     int32
	FragmentRef int32
	Flags       uint32
	Scale       float32
}

func (e *PolyhedronRef) FragCode() int {
	return 0x18
}

func (e *PolyhedronRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.FragmentRef)
	enc.Uint32(e.Flags)
	enc.Float32(e.Scale)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodePolyhedronRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &PolyhedronRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.FragmentRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.Scale = dec.Float32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// DMSprite is DmSpriteDef in libeq, Alternate Mesh in openzone, DMSPRITEDEF in wld, LegacyMesh in lantern
type DMSprite struct {
	NameRef           int32
	Flags             uint32
	Fragment1Maybe    int16
	MaterialReference uint32
	Fragment3         uint32
	CenterPosition    common.Vector3
	Params2           uint32
	Something2        uint32
	Something3        uint32
	Verticies         []common.Vector3
	TexCoords         []common.Vector3
	Normals           []common.Vector3
	Colors            []int32
	Polygons          []SpritePolygon
	VertexPieces      []SpriteVertexPiece
	PostVertexFlag    uint32
	RenderGroups      []SpriteRenderGroup
	VertexTex         []common.Vector2
	Size6Pieces       []Size6Entry
}

type SpritePolygon struct {
	Flag int16
	Unk1 int16
	Unk2 int16
	Unk3 int16
	Unk4 int16
	I1   int16
	I2   int16
	I3   int16
}

type SpriteVertexPiece struct {
	Count  int16
	Offset int16
}

type SpriteRenderGroup struct {
	PolygonCount int16
	MaterialId   int16
}

type Size6Entry struct {
	Unk1 uint32
	Unk2 uint32
	Unk3 uint32
	Unk4 uint32
	Unk5 uint32
}

func (e *DMSprite) FragCode() int {
	return 0x2C
}

func (e *DMSprite) Encode(w io.Writer) error {
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
	for _, vertex := range e.Verticies {
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
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeDMSprite(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &DMSprite{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	vertexCount := dec.Int16()
	texCoordCount := dec.Uint32()
	normalCount := dec.Uint32()
	colorCount := dec.Uint32()
	polygonCount := dec.Uint32()
	size6 := dec.Uint32()
	d.Fragment1Maybe = dec.Int16()
	vertexPieceCount := dec.Uint32()
	d.MaterialReference = dec.Uint32()
	d.Fragment3 = dec.Uint32()
	d.CenterPosition.X = dec.Float32()
	d.CenterPosition.Y = dec.Float32()
	d.CenterPosition.Z = dec.Float32()
	d.Params2 = dec.Uint32()
	d.Something2 = dec.Uint32()
	d.Something3 = dec.Uint32()

	if vertexCount > 999 {
		return nil, fmt.Errorf("vertex count misaligned")
	}

	if texCoordCount > 999 {
		return nil, fmt.Errorf("tex coord count misaligned")
	}

	if normalCount > 999 {
		return nil, fmt.Errorf("normal count misaligned")
	}

	for i := int16(0); i < vertexCount; i++ {
		v := common.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		d.Verticies = append(d.Verticies, v)
	}

	for i := uint32(0); i < texCoordCount; i++ {
		v := common.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		d.TexCoords = append(d.TexCoords, v)
	}

	for i := uint32(0); i < normalCount; i++ {
		v := common.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		d.Normals = append(d.Normals, v)
	}

	for i := uint32(0); i < colorCount; i++ {
		d.Colors = append(d.Colors, dec.Int32())
	}

	for i := uint32(0); i < polygonCount; i++ {
		p := SpritePolygon{}
		p.Flag = dec.Int16()
		p.Unk1 = dec.Int16()
		p.Unk2 = dec.Int16()
		p.Unk3 = dec.Int16()
		p.Unk4 = dec.Int16()
		p.I1 = dec.Int16()
		p.I2 = dec.Int16()
		p.I3 = dec.Int16()
		d.Polygons = append(d.Polygons, p)
	}

	for i := uint32(0); i < size6; i++ {
		s := Size6Entry{}
		s.Unk1 = dec.Uint32()
		s.Unk2 = dec.Uint32()
		s.Unk3 = dec.Uint32()
		s.Unk4 = dec.Uint32()
		s.Unk5 = dec.Uint32()
		d.Size6Pieces = append(d.Size6Pieces, s)
	}

	for i := uint32(0); i < vertexPieceCount; i++ {
		v := SpriteVertexPiece{}
		v.Count = dec.Int16()
		v.Offset = dec.Int16()
		d.VertexPieces = append(d.VertexPieces, v)
	}

	if d.Flags&9 != 0 {
		d.PostVertexFlag = dec.Uint32()
	}

	if d.Flags&11 != 0 {
		spriteRenderGroupCount := dec.Uint32()
		for i := uint32(0); i < spriteRenderGroupCount; i++ {
			s := SpriteRenderGroup{}
			s.PolygonCount = dec.Int16()
			s.MaterialId = dec.Int16()
			d.RenderGroups = append(d.RenderGroups, s)
		}
	}

	if d.Flags&12 != 0 {
		spriteVertexCount := dec.Uint32()
		for i := uint32(0); i < spriteVertexCount; i++ {
			v := common.Vector2{}
			v.X = dec.Float32()
			v.Y = dec.Float32()
			d.VertexTex = append(d.VertexTex, v)
		}
	}

	if d.Flags&13 != 0 {
		dec.Uint32()
		dec.Uint32()
		dec.Uint32()
	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}

	return d, nil
}

// DMSpriteRef is DmSprite in libeq, Mesh Reference in openzone, empty in wld, MeshReference in lantern
type DMSpriteRef struct {
	NameRef     int32  `yaml:"name_ref"`
	DMSpriteRef int32  `yaml:"dm_sprite_ref"`
	Params      uint32 `yaml:"params"`
}

func (e *DMSpriteRef) FragCode() int {
	return 0x2D
}

func (e *DMSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.DMSpriteRef)
	enc.Uint32(e.Params)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeDMSpriteRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &DMSpriteRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.DMSpriteRef = dec.Int32()
	d.Params = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// Mesh is DmSpriteDef2 in libeq, Mesh in openzone, DMSPRITEDEF2 in wld, Mesh in lantern
type Mesh struct {
	NameRef  int32
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
	VertexMaterials   []MeshVertexPiece
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

func (e *Mesh) FragCode() int {
	return 0x36
}

func (e *Mesh) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.MaterialListRef)
	enc.Int32(e.AnimationRef)
	enc.Uint32(0)
	enc.Uint32(0)
	enc.Float32(e.Center.X)
	enc.Float32(e.Center.Y)
	enc.Float32(e.Center.Z)
	enc.Uint32(0)
	enc.Uint32(0)
	enc.Uint32(0)
	enc.Float32(e.MaxDistance)
	enc.Float32(e.Min.X)
	enc.Float32(e.Min.Y)
	enc.Float32(e.Min.Z)
	enc.Float32(e.Max.X)
	enc.Float32(e.Max.Y)
	enc.Float32(e.Max.Z)
	enc.Uint16(uint16(len(e.Vertices)))
	enc.Uint16(uint16(len(e.Uvs)))
	enc.Uint16(uint16(len(e.Normals)))
	enc.Uint16(uint16(len(e.Colors)))
	enc.Uint16(uint16(len(e.Triangles)))
	enc.Uint16(uint16(len(e.VertexPieces)))
	enc.Uint16(uint16(len(e.TriangleMaterials)))
	enc.Uint16(uint16(len(e.VertexMaterials)))
	enc.Uint16(uint16(len(e.AnimatedBones)))
	rawScale := uint16(math.Log2(float64(1 / e.Scale)))
	enc.Uint16(rawScale)
	for _, vertex := range e.Vertices {
		enc.Int16(int16(int(vertex.Position.X-e.Center.X) * (1 << rawScale)))
		enc.Int16(int16(int(vertex.Position.Y-e.Center.Y) * (1 << rawScale)))
		enc.Int16(int16(int(vertex.Position.Z-e.Center.Z) * (1 << rawScale)))
	}
	for _, uv := range e.Uvs {
		if isOldWorld {
			enc.Int16(int16(uv.X * 256))
			enc.Int16(int16(uv.Y * 256))
		} else {
			enc.Int32(int32(uv.X * 256))
			enc.Int32(int32(uv.Y * 256))
		}
	}
	for _, normal := range e.Normals {
		enc.Int8(int8(normal.X * 128))
		enc.Int8(int8(normal.Y * 128))
		enc.Int8(int8(normal.Z * 128))
	}

	for _, color := range e.Colors {
		enc.Uint8(color.R)
		enc.Uint8(color.G)
		enc.Uint8(color.B)
		enc.Uint8(color.A)
	}

	for _, triangle := range e.Triangles {
		enc.Uint16(uint16(triangle.Flag & 1))
		enc.Uint16(uint16(triangle.Index.X))
		enc.Uint16(uint16(triangle.Index.Y))
		enc.Uint16(uint16(triangle.Index.Z))
	}

	for _, vertexPiece := range e.VertexPieces {
		enc.Int16(vertexPiece.Count)
		enc.Int16(vertexPiece.Index1)
	}

	for _, triangleMaterial := range e.TriangleMaterials {
		enc.Uint16(triangleMaterial.Count)
		enc.Uint16(triangleMaterial.MaterialID)
	}

	for _, vertexMaterial := range e.VertexMaterials {
		enc.Uint16(uint16(vertexMaterial.Count))
		enc.Uint16(uint16(vertexMaterial.Index1))
	}

	for _, animatedBone := range e.AnimatedBones {
		enc.Float32(animatedBone.Position.X)
		enc.Float32(animatedBone.Position.Y)
		enc.Float32(animatedBone.Position.Z)
	}

	if enc.Error() != nil {
		return enc.Error()
	}

	return nil
}

func decodeMesh(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Mesh{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32() // flags, currently unknown, zone meshes are 0x00018003, placeable objects are 0x00014003

	if d.Flags == 0x00018003 {
		d.FileType = "ter"
	}

	if d.Flags == 0x00014003 {
		d.FileType = "mod"
	}

	d.MaterialListRef = dec.Uint32()
	d.AnimationRef = dec.Int32() //used by flags/trees only

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

	d.MaxDistance = dec.Float32() // Given the values in center, this seems to contain the maximum distance between any vertex and that position. It seems to define a radius from that position within which the mesh lies.
	d.Min.X = dec.Float32()       // min x, y, and z coords in absolute coords of any vertex in the mesh.
	d.Min.Y = dec.Float32()
	d.Min.Z = dec.Float32()

	d.Max.X = dec.Float32() // max x, y, and z coords in absolute coords of any vertex in the mesh.
	d.Max.Y = dec.Float32()
	d.Max.Z = dec.Float32()

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
	triangleMaterialCount := dec.Uint16() // number of triangle texture entries. faces are grouped together by material and polygon material entries. This tells the client the number of faces using a material.
	vertexMaterialCount := dec.Uint16()   // number of vertex material entries. Vertices are grouped together by material and vertex material entries tell the client how many vertices there are using a material.

	meshAnimatedBoneCount := dec.Uint16() // number of entries in meshops. Seems to be used only for animated mob models.
	rawScale := dec.Uint16()
	scale := float32(1 / float32(int(1)<<rawScale)) // This allows vertex coordinates to be stored as integral values instead of floating-point values, without losing precision based on mesh size. Vertex values are multiplied by (1 shl `scale`) and stored in the vertex entries. FPSCALE is the internal name.
	// convert scale back to rawscale
	//rawScale = uint16(math.Log2(float64(1 / scale)))

	/// Vertices (x, y, z) belonging to this mesh. Each axis should
	/// be multiplied by (1 shl `scale`) for the final vertex position.
	for i := 0; i < int(vertexCount); i++ {
		vert := common.Vertex{}
		vert.Position.X = float32(centerX) + (float32(dec.Int16()) * scale)
		vert.Position.Y = float32(centerY) + (float32(dec.Int16()) * scale)
		vert.Position.Z = float32(centerZ) + (float32(dec.Int16()) * scale)
		d.Vertices = append(d.Vertices, vert)
	}

	for i := 0; i < int(uvCount); i++ {
		uv := common.Vector2{}
		if isOldWorld {
			uv.X = float32(dec.Int16()) / 256
			uv.Y = float32(dec.Int16()) / 256
		} else {
			uv.X = float32(dec.Int32()) / 256
			uv.Y = float32(dec.Int32()) / 256
		}
		d.Vertices[i].Uv = uv
	}

	for i := 0; i < int(normalCount); i++ {
		normal := common.Vector3{}
		normal.X = float32(dec.Int8()) / 128
		normal.Y = float32(dec.Int8()) / 128
		normal.Z = float32(dec.Int8()) / 128
		if i < len(d.Vertices) {
			d.Vertices[i].Normal = normal
		}
	}

	for i := 0; i < int(colorCount); i++ {
		color := common.RGBA{}
		color.R = dec.Uint8()
		color.G = dec.Uint8()
		color.B = dec.Uint8()
		color.A = dec.Uint8()
		if i < len(d.Vertices) {
			d.Vertices[i].Tint = color
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
		d.Triangles = append(d.Triangles, triangle)
	}

	for i := 0; i < int(vertexPieceCount); i++ {
		vertexPiece := MeshVertexPiece{}
		vertexPiece.Count = dec.Int16()
		vertexPiece.Index1 = dec.Int16()

		d.VertexPieces = append(d.VertexPieces, vertexPiece)
	}

	for i := 0; i < int(triangleMaterialCount); i++ {
		d.TriangleMaterials = append(d.TriangleMaterials, MeshTriangleMaterial{
			Count:      dec.Uint16(),
			MaterialID: dec.Uint16(),
		})
	}

	for i := 0; i < int(vertexMaterialCount); i++ {
		vertexMat := MeshVertexPiece{}
		vertexMat.Count = dec.Int16()
		vertexMat.Index1 = dec.Int16()
		d.VertexMaterials = append(d.VertexMaterials, vertexMat)
	}

	for i := 0; i < int(meshAnimatedBoneCount); i++ {
		mab := MeshAnimatedBone{}
		mab.Position.X = float32(centerX) + (float32(dec.Int16()) * scale)
		mab.Position.Y = float32(centerY) + (float32(dec.Int16()) * scale)
		mab.Position.Z = float32(centerZ) + (float32(dec.Int16()) * scale)
		d.AnimatedBones = append(d.AnimatedBones, mab)
	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}

	return d, nil
}

// MeshAnimated is DmTrackDef2 in libeq, Mesh Animated Vertices in openzone, DMTRACKDEF in wld, MeshAnimatedVertices in lantern
type MeshAnimated struct {
	NameRef     int32
	Flags       uint32
	VertexCount uint16
	FrameCount  uint16
	Param1      uint16 // usually contains 100
	Param2      uint16 // usually contains 0
	Scale       uint16
	Frames      []MeshAnimatedBone
	Size6       uint32
}

func (e *MeshAnimated) FragCode() int {
	return 0x37
}

func (e *MeshAnimated) Encode(w io.Writer) error {
	return nil
}

func decodeMeshAnimated(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &MeshAnimated{}
	return d, nil
}
