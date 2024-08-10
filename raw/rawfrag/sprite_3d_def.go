package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/tag"
)

// WldFragSprite3DDef is Sprite3DDef in libeq, Camera in openzone, 3DSPRITEDEF in wld, Camera in lantern
type WldFragSprite3DDef struct {
	NameRef        int32
	Flags          uint32
	SphereListRef  uint32
	CenterOffset   [3]float32
	BoundingRadius float32
	Vertices       [][3]float32
	BspNodes       []WldFragThreeDSpriteBspNode
}

type WldFragThreeDSpriteBspNode struct {
	FrontTree                   uint32
	BackTree                    uint32
	VertexIndexes               []uint32
	RenderMethod                uint32
	RenderFlags                 uint8
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          [3]float32
	RenderUVInfoUAxis           [3]float32
	RenderUVInfoVAxis           [3]float32
	Uvs                         [][2]float32
}

func (e *WldFragSprite3DDef) FragCode() int {
	return FragCodeSprite3DDef
}

func (e *WldFragSprite3DDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Vertices)))
	enc.Uint32(uint32(len(e.BspNodes)))
	enc.Uint32(e.SphereListRef)
	if e.Flags&0x01 == 0x01 {
		enc.Float32(e.CenterOffset[0])
		enc.Float32(e.CenterOffset[1])
		enc.Float32(e.CenterOffset[2])
	}
	if e.Flags&0x02 == 0x02 {
		enc.Float32(e.BoundingRadius)
	}
	for _, vertex := range e.Vertices {
		enc.Float32(vertex[0])
		enc.Float32(vertex[1])
		enc.Float32(vertex[2])
	}
	for _, node := range e.BspNodes {
		enc.Uint32(uint32(len(node.VertexIndexes)))
		enc.Uint32(node.FrontTree)
		enc.Uint32(node.BackTree)
		for _, vertexIndex := range node.VertexIndexes {
			enc.Uint32(vertexIndex)
		}

		enc.Uint32(node.RenderMethod)
		enc.Uint8(node.RenderFlags)

		if node.RenderFlags&0x01 == 0x01 {
			enc.Uint32(node.RenderPen)
		}
		if node.RenderFlags&0x02 == 0x02 {
			enc.Float32(node.RenderBrightness)
		}
		if node.RenderFlags&0x04 == 0x04 {
			enc.Float32(node.RenderScaledAmbient)
		}
		if node.RenderFlags&0x08 == 0x08 {
			enc.Uint32(node.RenderSimpleSpriteReference)
		}
		if node.RenderFlags&0x10 == 0x10 {
			enc.Float32(node.RenderUVInfoOrigin[0])
			enc.Float32(node.RenderUVInfoOrigin[1])
			enc.Float32(node.RenderUVInfoOrigin[2])
			enc.Float32(node.RenderUVInfoUAxis[0])
			enc.Float32(node.RenderUVInfoUAxis[1])
			enc.Float32(node.RenderUVInfoUAxis[2])
			enc.Float32(node.RenderUVInfoVAxis[0])
			enc.Float32(node.RenderUVInfoVAxis[1])
			enc.Float32(node.RenderUVInfoVAxis[2])
		}
		if node.RenderFlags&0x20 == 0x20 {
			enc.Uint32(uint32(len(node.Uvs)))
			for _, uv := range node.Uvs {
				enc.Float32(uv[0])
				enc.Float32(uv[1])
			}
		}
		// two sided is 0x40 on flag, not needed to write
	}
	enc.Byte(0x00)
	enc.Byte(0x00)
	enc.Byte(0x00)

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragSprite3DDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	vertexCount := dec.Uint32()
	bspNodeCount := dec.Uint32()
	e.SphereListRef = dec.Uint32()
	if e.Flags&0x01 == 0x01 {
		e.CenterOffset[0] = dec.Float32()
		e.CenterOffset[1] = dec.Float32()
		e.CenterOffset[2] = dec.Float32()
	}
	if e.Flags&0x02 == 0x02 {
		e.BoundingRadius = dec.Float32()
	}
	tag.AddRand(tag.LastPos(), dec.Pos(), "header")
	for i := 0; i < int(vertexCount); i++ {
		e.Vertices = append(e.Vertices, [3]float32{dec.Float32(), dec.Float32(), dec.Float32()})
	}
	tag.AddRandf(tag.LastPos(), dec.Pos(), "verts=%d", vertexCount)
	for i := 0; i < int(bspNodeCount); i++ {
		node := WldFragThreeDSpriteBspNode{}
		vertexIndexCount := dec.Uint32()
		node.FrontTree = dec.Uint32()
		node.BackTree = dec.Uint32()
		for j := 0; j < int(vertexIndexCount); j++ {
			node.VertexIndexes = append(node.VertexIndexes, dec.Uint32())
		}
		node.RenderMethod = dec.Uint32()
		node.RenderFlags = dec.Uint8()

		if node.RenderFlags&0x01 == 0x01 {
			node.RenderPen = dec.Uint32()
		}
		if node.RenderFlags&0x02 == 0x02 {
			node.RenderBrightness = dec.Float32()
		}
		if node.RenderFlags&0x04 == 0x04 {
			node.RenderScaledAmbient = dec.Float32()
		}
		if node.RenderFlags&0x08 == 0x08 {
			node.RenderSimpleSpriteReference = dec.Uint32()
		}
		if node.RenderFlags&0x10 == 0x10 {
			node.RenderUVInfoOrigin[0] = dec.Float32()
			node.RenderUVInfoOrigin[1] = dec.Float32()
			node.RenderUVInfoOrigin[2] = dec.Float32()
			node.RenderUVInfoUAxis[0] = dec.Float32()
			node.RenderUVInfoUAxis[1] = dec.Float32()
			node.RenderUVInfoUAxis[2] = dec.Float32()
			node.RenderUVInfoVAxis[0] = dec.Float32()
			node.RenderUVInfoVAxis[1] = dec.Float32()
			node.RenderUVInfoVAxis[2] = dec.Float32()
		}
		if node.RenderFlags&0x20 == 0x20 {
			renderUVMapEntryCount := dec.Uint32()
			for j := 0; j < int(renderUVMapEntryCount); j++ {
				u := [2]float32{dec.Float32(), dec.Float32()}
				node.Uvs = append(node.Uvs, u)
			}
		}
		e.BspNodes = append(e.BspNodes, node)
		tag.AddRandf(tag.LastPos(), dec.Pos(), "%d bspNode", i)
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil

}
