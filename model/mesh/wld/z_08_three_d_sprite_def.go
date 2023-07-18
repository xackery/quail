package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// 0x08 8 typically a camera
type threeDSpriteDef struct {
	NameRef       int32
	Flags         uint32
	VertexCount   uint32
	BspNodeCount  uint32
	SphereListRef uint32
	CenterOffset  geo.Vector3
	Radius        float32
	Vertices      []geo.Vector3
	BspNodes      []bspNode
}

type bspNode struct {
	vertexCount   uint32
	FrontTree     uint32
	BackTree      uint32
	VertexIndexes []uint32
	RenderMethod  uint32
	RenderInfo    renderInfo
}

type renderInfo struct {
	flags                 uint32
	pen                   uint32  // only exists if bit 0 of flags is set
	brightness            float32 // only exists if bit 1 of flags is set
	scaledAmbient         float32 // only exists if bit 2 of flags is set
	simpleSpriteReference uint32  // only exists if bit 3 of flags is set
	uvInfo                uvInfo  // only exists if bit 4 of flags is set
	uvMap                 uvMap   // only exists if bit 5 of flags is set
}

type uvInfo struct {
	origin geo.Vector3
	uAxis  geo.Vector3
	vAxis  geo.Vector3
}

type uvMap struct {
	entryCount uint32
	entries    []geo.Vector2
}

func (e *WLD) threeDSpriteDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &threeDSpriteDef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.Flags = dec.Uint32()
	def.VertexCount = dec.Uint32()
	def.BspNodeCount = dec.Uint32()
	def.SphereListRef = dec.Uint32()
	def.CenterOffset.X = dec.Float32()
	def.CenterOffset.Y = dec.Float32()
	def.CenterOffset.Z = dec.Float32()
	def.Radius = dec.Float32()
	for i := 0; i < int(def.VertexCount); i++ {
		v := geo.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		def.Vertices = append(def.Vertices, v)
	}
	for i := 0; i < int(def.BspNodeCount); i++ {
		v := bspNode{}
		v.vertexCount = dec.Uint32()
		v.FrontTree = dec.Uint32()
		v.BackTree = dec.Uint32()
		for j := 0; j < int(v.vertexCount); j++ {
			v.VertexIndexes = append(v.VertexIndexes, dec.Uint32())
		}
		v.RenderMethod = dec.Uint32()
		v.RenderInfo.flags = dec.Uint32()
		if v.RenderInfo.flags&0x01 == 0x01 {
			v.RenderInfo.pen = dec.Uint32()
		}
		if v.RenderInfo.flags&0x02 == 0x02 {
			v.RenderInfo.brightness = dec.Float32()
		}
		if v.RenderInfo.flags&0x04 == 0x04 {
			v.RenderInfo.scaledAmbient = dec.Float32()
		}
		if v.RenderInfo.flags&0x08 == 0x08 {
			v.RenderInfo.simpleSpriteReference = dec.Uint32()
		}
		if v.RenderInfo.flags&0x10 == 0x10 {
			v.RenderInfo.uvInfo.origin.X = dec.Float32()
			v.RenderInfo.uvInfo.origin.Y = dec.Float32()
			v.RenderInfo.uvInfo.origin.Z = dec.Float32()
			v.RenderInfo.uvInfo.uAxis.X = dec.Float32()
			v.RenderInfo.uvInfo.uAxis.Y = dec.Float32()
			v.RenderInfo.uvInfo.uAxis.Z = dec.Float32()
			v.RenderInfo.uvInfo.vAxis.X = dec.Float32()
			v.RenderInfo.uvInfo.vAxis.Y = dec.Float32()
			v.RenderInfo.uvInfo.vAxis.Z = dec.Float32()
		}
		if v.RenderInfo.flags&0x20 == 0x20 {
			v.RenderInfo.uvMap.entryCount = dec.Uint32()
			for j := 0; j < int(v.RenderInfo.uvMap.entryCount); j++ {
				vv := geo.Vector2{}
				vv.X = dec.Float32()
				vv.Y = dec.Float32()
				v.RenderInfo.uvMap.entries = append(v.RenderInfo.uvMap.entries, vv)
			}
		}

		def.BspNodes = append(def.BspNodes, v)
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("threeDSpriteRead: %v", err)
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *threeDSpriteDef) build(e *WLD) error {
	return nil
}

func (e *WLD) threeDSpriteDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
