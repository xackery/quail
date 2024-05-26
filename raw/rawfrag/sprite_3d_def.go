package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/tag"
)

// WldFragSprite3DDef is Sprite3DDef in libeq, Camera in openzone, 3DSPRITEDEF in wld, Camera in lantern
type WldFragSprite3DDef struct {
	NameRef       int32                        `yaml:"name_ref"`
	Flags         uint32                       `yaml:"flags"`
	SphereListRef uint32                       `yaml:"sphere_list_ref"`
	CenterOffset  model.Vector3                `yaml:"center_offset"`
	Radius        float32                      `yaml:"radius"`
	Vertices      []model.Vector3              `yaml:"vertices"`
	BspNodes      []WldFragThreeDSpriteBspNode `yaml:"bsp_nodes"`
}

type WldFragThreeDSpriteBspNode struct {
	FrontTree                   uint32                             `yaml:"front_tree"`
	BackTree                    uint32                             `yaml:"back_tree"`
	VertexIndexes               []uint32                           `yaml:"vertex_indexes"`
	RenderMethod                uint32                             `yaml:"render_method"`
	RenderFlags                 uint8                              `yaml:"render_flags"`
	RenderPen                   uint32                             `yaml:"render_pen"`
	RenderBrightness            float32                            `yaml:"render_brightness"`
	RenderScaledAmbient         float32                            `yaml:"render_scaled_ambient"`
	RenderSimpleSpriteReference uint32                             `yaml:"render_simple_sprite_reference"`
	RenderUVInfoOrigin          model.Vector3                      `yaml:"render_uv_info_origin"`
	RenderUVInfoUAxis           model.Vector3                      `yaml:"render_uv_info_u_axis"`
	RenderUVInfoVAxis           model.Vector3                      `yaml:"render_uv_info_v_axis"`
	RenderUVMapEntries          []WldFragThreeDSpriteBspNodeUVInfo `yaml:"render_uv_map_entries"`
}

type WldFragThreeDSpriteBspNodeUVInfo struct {
	UvOrigin [3]float32
	UAxis    [3]float32
	VAxis    [3]float32
}

func (e *WldFragSprite3DDef) FragCode() int {
	return FragCodeSprite3DDef
}

func (e *WldFragSprite3DDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Vertices)))
	enc.Uint32(uint32(len(e.BspNodes)))
	enc.Uint32(e.SphereListRef)
	if e.Flags&0x01 == 0x01 {
		enc.Float32(e.CenterOffset.X)
		enc.Float32(e.CenterOffset.Y)
		enc.Float32(e.CenterOffset.Z)
	}
	if e.Flags&0x02 == 0x02 {
		enc.Float32(e.Radius)
	}
	for _, vertex := range e.Vertices {
		enc.Float32(vertex.X)
		enc.Float32(vertex.Y)
		enc.Float32(vertex.Z)
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
			enc.Float32(node.RenderUVInfoOrigin.X)
			enc.Float32(node.RenderUVInfoOrigin.Y)
			enc.Float32(node.RenderUVInfoOrigin.Z)
			enc.Float32(node.RenderUVInfoUAxis.X)
			enc.Float32(node.RenderUVInfoUAxis.Y)
			enc.Float32(node.RenderUVInfoUAxis.Z)
			enc.Float32(node.RenderUVInfoVAxis.X)
			enc.Float32(node.RenderUVInfoVAxis.Y)
			enc.Float32(node.RenderUVInfoVAxis.Z)
		}
		if node.RenderFlags&0x20 == 0x20 {
			enc.Uint32(uint32(len(node.RenderUVMapEntries)))
			for _, entry := range node.RenderUVMapEntries {
				enc.Float32(entry.UvOrigin[0])
				enc.Float32(entry.UvOrigin[1])
				enc.Float32(entry.UvOrigin[2])
				enc.Float32(entry.UAxis[0])
				enc.Float32(entry.UAxis[1])
				enc.Float32(entry.UAxis[2])
				enc.Float32(entry.VAxis[0])
				enc.Float32(entry.VAxis[1])
				enc.Float32(entry.VAxis[2])
			}
		}
		// two sided is 0x40
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

func (e *WldFragSprite3DDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	vertexCount := dec.Uint32()
	bspNodeCount := dec.Uint32()
	e.SphereListRef = dec.Uint32()
	if e.Flags&0x01 == 0x01 {
		e.CenterOffset.X = dec.Float32()
		e.CenterOffset.Y = dec.Float32()
		e.CenterOffset.Z = dec.Float32()
	}
	if e.Flags&0x02 == 0x02 {
		e.Radius = dec.Float32()
	}
	tag.AddRand(tag.LastPos(), dec.Pos(), "header")
	for i := 0; i < int(vertexCount); i++ {
		v := model.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		e.Vertices = append(e.Vertices, v)
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
			node.RenderUVInfoOrigin.X = dec.Float32()
			node.RenderUVInfoOrigin.Y = dec.Float32()
			node.RenderUVInfoOrigin.Z = dec.Float32()
			node.RenderUVInfoUAxis.X = dec.Float32()
			node.RenderUVInfoUAxis.Y = dec.Float32()
			node.RenderUVInfoUAxis.Z = dec.Float32()
			node.RenderUVInfoVAxis.X = dec.Float32()
			node.RenderUVInfoVAxis.Y = dec.Float32()
			node.RenderUVInfoVAxis.Z = dec.Float32()
		}
		if node.RenderFlags&0x20 == 0x20 {
			renderUVMapEntryCount := dec.Uint32()
			for j := 0; j < int(renderUVMapEntryCount); j++ {
				v := WldFragThreeDSpriteBspNodeUVInfo{
					UvOrigin: [3]float32{dec.Float32(), dec.Float32(), dec.Float32()},
					UAxis:    [3]float32{dec.Float32(), dec.Float32(), dec.Float32()},
					VAxis:    [3]float32{dec.Float32(), dec.Float32(), dec.Float32()},
				}
				node.RenderUVMapEntries = append(node.RenderUVMapEntries, v)
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
