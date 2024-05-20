package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/tag"
)

// WldFragSprite2DDef is Sprite2DDef in libeq, Two-Dimensional Object in openzone, 2DSPRITEDEF in wld, Fragment06 in lantern
type WldFragSprite2DDef struct {
	FragName                    string    `yaml:"frag_name"`
	NameRef                     int32     `yaml:"name_ref"`
	Flags                       uint32    `yaml:"flags"`
	TextureCount                uint32    `yaml:"texture_count"`
	PitchCount                  uint32    `yaml:"pitch_count"`
	Scale                       Vector2   `yaml:"scale"`
	SphereRef                   uint32    `yaml:"sphere_ref"`
	DepthScale                  float32   `yaml:"depth_scale"`
	CenterOffset                Vector3   `yaml:"center_offset"`
	BoundingRadius              float32   `yaml:"bounding_radius"`
	CurrentFrameRef             int32     `yaml:"current_frame_ref"`
	Sleep                       uint32    `yaml:"sleep"`
	Headings                    []uint32  `yaml:"headings"`
	RenderMethod                uint32    `yaml:"render_method"`
	RenderFlags                 uint32    `yaml:"render_flags"`
	RenderPen                   uint32    `yaml:"render_pen"`
	RenderBrightness            float32   `yaml:"render_brightness"`
	RenderScaledAmbient         float32   `yaml:"render_scaled_ambient"`
	RenderSimpleSpriteReference uint32    `yaml:"render_simple_sprite_reference"`
	RenderUVInfoOrigin          Vector3   `yaml:"render_uv_info_origin"`
	RenderUVInfoUAxis           Vector3   `yaml:"render_uv_info_u_axis"`
	RenderUVInfoVAxis           Vector3   `yaml:"render_uv_info_v_axis"`
	RenderUVMapEntries          []Vector2 `yaml:"render_uv_map_entries"`
}

func (e *WldFragSprite2DDef) FragCode() int {
	return FragCodeSprite2DDef
}

func (e *WldFragSprite2DDef) Write(w io.Writer) error {
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
			enc.Uint32(uint32(len(e.RenderUVMapEntries)))
			for _, entry := range e.RenderUVMapEntries {
				enc.Float32(entry.X)
				enc.Float32(entry.Y)
			}
		}
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil

}

func (e *WldFragSprite2DDef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.TextureCount = dec.Uint32()
	e.PitchCount = dec.Uint32()
	e.Scale.X = dec.Float32()
	e.Scale.Y = dec.Float32()
	e.SphereRef = dec.Uint32()
	if e.Flags&0x80 == 0x80 {
		e.DepthScale = dec.Float32()
	}
	if e.Flags&0x01 == 0x01 {
		e.CenterOffset.X = dec.Float32()
		e.CenterOffset.Y = dec.Float32()
		e.CenterOffset.Z = dec.Float32()
	}
	if e.Flags&0x02 == 0x02 {
		e.BoundingRadius = dec.Float32()
	}
	if e.Flags&0x04 == 0x04 {
		e.CurrentFrameRef = dec.Int32()
	}
	if e.Flags&0x08 == 0x08 {
		e.Sleep = dec.Uint32()
	}
	e.Headings = make([]uint32, e.PitchCount)
	for i := uint32(0); i < e.PitchCount; i++ {
		e.Headings[i] = dec.Uint32()
	}
	if e.Flags&0x10 == 0x10 {
		e.RenderMethod = dec.Uint32()
		e.RenderFlags = dec.Uint32()
		if e.RenderFlags&0x01 == 0x01 {
			e.RenderPen = dec.Uint32()
		}
		if e.RenderFlags&0x02 == 0x02 {
			e.RenderBrightness = dec.Float32()
		}
		if e.RenderFlags&0x04 == 0x04 {
			e.RenderScaledAmbient = dec.Float32()
		}
		if e.RenderFlags&0x08 == 0x08 {
			e.RenderSimpleSpriteReference = dec.Uint32()
		}
		if e.RenderFlags&0x10 == 0x10 {
			e.RenderUVInfoOrigin.X = dec.Float32()
			e.RenderUVInfoOrigin.Y = dec.Float32()
			e.RenderUVInfoOrigin.Z = dec.Float32()
			e.RenderUVInfoUAxis.X = dec.Float32()
			e.RenderUVInfoUAxis.Y = dec.Float32()
			e.RenderUVInfoUAxis.Z = dec.Float32()
			e.RenderUVInfoVAxis.X = dec.Float32()
			e.RenderUVInfoVAxis.Y = dec.Float32()
			e.RenderUVInfoVAxis.Z = dec.Float32()
		}
		if e.RenderFlags&0x20 == 0x20 {
			renderUVMapEntrycount := dec.Uint32()
			for i := uint32(0); i < renderUVMapEntrycount; i++ {
				v := Vector2{}
				v.X = dec.Float32()
				v.Y = dec.Float32()
				e.RenderUVMapEntries = append(e.RenderUVMapEntries, v)
			}
		}
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil

}

// WldFragSprite2D is Sprite2D in libeq, Two-Dimensional Object Reference in openzone, 2DSPRITE (ref) in wld, Fragment07 in lantern
type WldFragSprite2D struct {
	FragName      string `yaml:"frag_name"`
	NameRef       int32  `yaml:"name_ref"`
	TwoDSpriteRef uint32 `yaml:"two_d_sprite_ref"`
	Flags         uint32 `yaml:"flags"`
}

func (e *WldFragSprite2D) FragCode() int {
	return FragCodeSprite2D
}

func (e *WldFragSprite2D) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.TwoDSpriteRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite2D) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.TwoDSpriteRef = dec.Uint32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragSprite3DDef is Sprite3DDef in libeq, Camera in openzone, 3DSPRITEDEF in wld, Camera in lantern
type WldFragSprite3DDef struct {
	FragName      string                       `yaml:"frag_name"`
	NameRef       int32                        `yaml:"name_ref"`
	Flags         uint32                       `yaml:"flags"`
	SphereListRef uint32                       `yaml:"sphere_list_ref"`
	CenterOffset  Vector3                      `yaml:"center_offset"`
	Radius        float32                      `yaml:"radius"`
	Vertices      []Vector3                    `yaml:"vertices"`
	BspNodes      []WldFragThreeDSpriteBspNode `yaml:"bsp_nodes"`
}

type WldFragThreeDSpriteBspNode struct {
	FrontTree                   uint32    `yaml:"front_tree"`
	BackTree                    uint32    `yaml:"back_tree"`
	VertexIndexes               []uint32  `yaml:"vertex_indexes"`
	RenderMethod                uint8     `yaml:"render_method"`
	RenderFlags                 uint8     `yaml:"render_flags"`
	RenderPen                   uint8     `yaml:"render_pen"`
	RenderBrightness            uint8     `yaml:"render_brightness"`
	RenderScaledAmbient         uint8     `yaml:"render_scaled_ambient"`
	RenderSimpleSpriteReference uint8     `yaml:"render_simple_sprite_reference"`
	RenderUVInfoOrigin          Vector3   `yaml:"render_uv_info_origin"`
	RenderUVInfoUAxis           Vector3   `yaml:"render_uv_info_u_axis"`
	RenderUVInfoVAxis           Vector3   `yaml:"render_uv_info_v_axis"`
	RenderUVMapEntries          []Vector2 `yaml:"render_uv_map_entries"`
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
	enc.Float32(e.CenterOffset.X)
	enc.Float32(e.CenterOffset.Y)
	enc.Float32(e.CenterOffset.Z)
	enc.Float32(e.Radius)
	tag.AddRand(tag.LastPos(), enc.Pos(), "header")
	for _, vertex := range e.Vertices {
		enc.Float32(vertex.X)
		enc.Float32(vertex.Y)
		enc.Float32(vertex.Z)
	}
	tag.AddRandf(tag.LastPos(), enc.Pos(), "verts=%d", len(e.Vertices))
	for _, node := range e.BspNodes {
		enc.Uint32(uint32(len(node.VertexIndexes)))
		enc.Uint32(node.FrontTree)
		enc.Uint32(node.BackTree)
		for _, vertexIndex := range node.VertexIndexes {
			enc.Uint32(vertexIndex)
		}

		enc.Uint8(node.RenderMethod)
		enc.Uint8(node.RenderFlags)
		tag.AddRandf(tag.LastPos(), enc.Pos(), "renderFlags=%d", node.RenderFlags)

		if node.RenderFlags&0x01 == 0x01 {
			enc.Uint8(node.RenderPen)
		}
		if node.RenderFlags&0x02 == 0x02 {
			enc.Uint8(node.RenderBrightness)
		}
		if node.RenderFlags&0x04 == 0x04 {
			enc.Uint8(node.RenderScaledAmbient)
		}
		if node.RenderFlags&0x08 == 0x08 {
			enc.Uint8(node.RenderSimpleSpriteReference)
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
			tag.AddRandf(tag.LastPos(), enc.Pos(), "renderUVInfoOrigin=%f,%f,%f", node.RenderUVInfoOrigin.X, node.RenderUVInfoOrigin.Y, node.RenderUVInfoOrigin.Z)
		}
		if node.RenderFlags&0x20 == 0x20 {
			enc.Uint32(uint32(len(node.RenderUVMapEntries)))
			for _, entry := range node.RenderUVMapEntries {
				enc.Float32(entry.X)
				enc.Float32(entry.Y)
			}
			tag.AddRandf(tag.LastPos(), enc.Pos(), "renderUVMapEntryCount=%d", len(node.RenderUVMapEntries))
		}
		tag.AddRandf(tag.LastPos(), enc.Pos(), "bspNode")
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragSprite3DDef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	vertexCount := dec.Uint32()
	bspNodeCount := dec.Uint32()
	e.SphereListRef = dec.Uint32()
	e.CenterOffset.X = dec.Float32()
	e.CenterOffset.Y = dec.Float32()
	e.CenterOffset.Z = dec.Float32()
	e.Radius = dec.Float32()
	tag.AddRand(tag.LastPos(), dec.Pos(), "header")
	for i := 0; i < int(vertexCount); i++ {
		v := Vector3{}
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
		node.RenderMethod = dec.Uint8()
		node.RenderFlags = dec.Uint8()

		if node.RenderFlags&0x01 == 0x01 {
			node.RenderPen = dec.Uint8()
		}
		if node.RenderFlags&0x02 == 0x02 {
			node.RenderBrightness = dec.Uint8()
		}
		if node.RenderFlags&0x04 == 0x04 {
			node.RenderScaledAmbient = dec.Uint8()
		}
		if node.RenderFlags&0x08 == 0x08 {
			node.RenderSimpleSpriteReference = dec.Uint8()
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
				v := Vector2{}
				v.X = dec.Float32()
				v.Y = dec.Float32()
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

// WldFragSprite3D is Sprite3D in libeq, Camera Reference in openzone, 3DSPRITE (ref) in wld, CameraReference in lantern
type WldFragSprite3D struct {
	FragName  string `yaml:"frag_name"`
	NameRef   int32
	ThreeDRef int32
	Flags     uint32
}

func (e *WldFragSprite3D) FragCode() int {
	return FragCodeSprite3D
}

func (e *WldFragSprite3D) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.ThreeDRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite3D) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.ThreeDRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragSprite4DDef is Sprite4DDef in libeq, empty in openzone, 4DSPRITEDEF in wld
type WldFragSprite4DDef struct {
	FragName        string   `yaml:"frag_name"`
	NameRef         int32    `yaml:"name_ref"`
	Flags           uint32   `yaml:"flags"`
	PolyRef         int32    `yaml:"poly_ref"`
	CenterOffset    Vector3  `yaml:"center_offset"`
	Radius          float32  `yaml:"radius"`
	CurrentFrame    uint32   `yaml:"current_frame"`
	Sleep           uint32   `yaml:"sleep"`
	SpriteFragments []uint32 `yaml:"sprite_fragments"`
}

func (e *WldFragSprite4DDef) FragCode() int {
	return FragCodeSprite4DDef
}

func (e *WldFragSprite4DDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.SpriteFragments)))
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
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite4DDef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	frameCount := dec.Uint32()
	e.PolyRef = dec.Int32()
	if e.Flags&0x01 != 0 {
		e.CenterOffset.X = dec.Float32()
		e.CenterOffset.Y = dec.Float32()
		e.CenterOffset.Z = dec.Float32()
	}
	if e.Flags&0x02 != 0 {
		e.Radius = dec.Float32()
	}
	if e.Flags&0x04 != 0 {
		e.CurrentFrame = dec.Uint32()
	}
	if e.Flags&0x08 != 0 {
		e.Sleep = dec.Uint32()
	}
	if e.Flags&0x10 != 0 {
		for i := uint32(0); i < frameCount; i++ {
			e.SpriteFragments = append(e.SpriteFragments, dec.Uint32())
		}
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragSprite4D is Sprite4D in libeq, empty in openzone, 4DSPRITE (ref) in wld
type WldFragSprite4D struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32  `yaml:"name_ref"`
	FourDRef int32  `yaml:"four_d_ref"`
	Params1  uint32 `yaml:"params_1"`
}

func (e *WldFragSprite4D) FragCode() int {
	return FragCodeSprite4D
}

func (e *WldFragSprite4D) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.FourDRef)
	enc.Uint32(e.Params1)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite4D) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.FourDRef = dec.Int32()
	e.Params1 = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragPolyhedronDef is PolyhedronDef in libeq, Polygon animation in openzone, POLYHEDRONDEFINITION in wld, Fragment17 in lantern
type WldFragPolyhedronDef struct {
	FragName string                      `yaml:"frag_name"`
	NameRef  int32                       `yaml:"name_ref"`
	Flags    uint32                      `yaml:"flags"`
	Size1    uint32                      `yaml:"size_1"`
	Size2    uint32                      `yaml:"size_2"`
	Params1  float32                     `yaml:"params_1"`
	Params2  float32                     `yaml:"params_2"`
	Entries1 []Vector3                   `yaml:"entries_1"`
	Entries2 []WldFragPolyhedronEntries2 `yaml:"entries_2"`
}

type WldFragPolyhedronEntries2 struct {
	Unk1 uint32   `yaml:"unk_1"`
	Unk2 []uint32 `yaml:"unk_2"`
}

func (e *WldFragPolyhedronDef) FragCode() int {
	return FragCodePolyhedronDef
}

func (e *WldFragPolyhedronDef) Write(w io.Writer) error {
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
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPolyhedronDef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.Size1 = dec.Uint32()
	e.Size2 = dec.Uint32()
	e.Params1 = dec.Float32()
	e.Params2 = dec.Float32()
	for i := uint32(0); i < e.Size1; i++ {
		v := Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		e.Entries1 = append(e.Entries1, v)
	}
	for i := uint32(0); i < e.Size2; i++ {
		entry := WldFragPolyhedronEntries2{}
		entry.Unk1 = dec.Uint32()
		for j := uint32(0); j < e.Size1; j++ {
			entry.Unk2 = append(entry.Unk2, dec.Uint32())
		}
		e.Entries2 = append(e.Entries2, entry)
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragPolyhedron is Polyhedron in libeq, Polygon Animation Reference in openzone, POLYHEDRON (ref) in wld, Fragment18 in lantern
type WldFragPolyhedron struct {
	FragName    string  `yaml:"frag_name"`
	NameRef     int32   `yaml:"name_ref"`
	FragmentRef int32   `yaml:"fragment_ref"`
	Flags       uint32  `yaml:"flags"`
	Scale       float32 `yaml:"scale"`
}

func (e *WldFragPolyhedron) FragCode() int {
	return FragCodePolyhedron
}

func (e *WldFragPolyhedron) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.FragmentRef)
	enc.Uint32(e.Flags)
	enc.Float32(e.Scale)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPolyhedron) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.FragmentRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.Scale = dec.Float32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragDMSpriteDef is DmSpriteDef in libeq, Alternate Mesh in openzone, DMSPRITEDEF in wld, LegacyMesh in lantern
type WldFragDMSpriteDef struct {
	FragName          string                         `yaml:"frag_name"`
	NameRef           int32                          `yaml:"name_ref"`
	Flags             uint32                         `yaml:"flags"`
	Fragment1Maybe    int16                          `yaml:"fragment_1_maybe"`
	MaterialReference uint32                         `yaml:"material_reference"`
	Fragment3         uint32                         `yaml:"fragment_3"`
	CenterPosition    Vector3                        `yaml:"center_position"`
	Params2           uint32                         `yaml:"params_2"`
	Something2        uint32                         `yaml:"something_2"`
	Something3        uint32                         `yaml:"something_3"`
	Verticies         []Vector3                      `yaml:"verticies"`
	TexCoords         []Vector3                      `yaml:"tex_coords"`
	Normals           []Vector3                      `yaml:"normals"`
	Colors            []int32                        `yaml:"colors"`
	Polygons          []WldFragDMSpriteSpritePolygon `yaml:"polygons"`
	VertexPieces      []WldFragDMSpriteVertexPiece   `yaml:"vertex_pieces"`
	PostVertexFlag    uint32                         `yaml:"post_vertex_flag"`
	RenderGroups      []WldFragDMSpriteRenderGroup   `yaml:"render_groups"`
	VertexTex         []Vector2                      `yaml:"vertex_tex"`
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
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDMSpriteDef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
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
		v := Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		e.Verticies = append(e.Verticies, v)
	}

	for i := uint32(0); i < texCoordCount; i++ {
		v := Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		e.TexCoords = append(e.TexCoords, v)
	}

	for i := uint32(0); i < normalCount; i++ {
		v := Vector3{}
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
			v := Vector2{}
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

// WldFragDMSprite is DmSprite in libeq, Mesh Reference in openzone, empty in wld, MeshReference in lantern
type WldFragDMSprite struct {
	FragName    string `yaml:"frag_name"`
	NameRef     int32  `yaml:"name_ref"`
	DMSpriteRef int32  `yaml:"dm_sprite_ref"`
	Params      uint32 `yaml:"params"`
}

func (e *WldFragDMSprite) FragCode() int {
	return FragCodeDMSprite
}

func (e *WldFragDMSprite) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.DMSpriteRef)
	enc.Uint32(e.Params)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDMSprite) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.DMSpriteRef = dec.Int32()
	e.Params = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragDmSpriteDef2 is DmSpriteDef2 in libeq, WldFragDmSpriteDef2 in openzone, DMSPRITEDEF2 in wld, WldFragDmSpriteDef2 in lantern
type WldFragDmSpriteDef2 struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32  `yaml:"name_ref"`
	Flags    uint32 `yaml:"flags"`

	MaterialListRef uint32 `yaml:"material_list_ref"`
	AnimationRef    int32  `yaml:"animation_ref"`

	Fragment3Ref int32   `yaml:"fragment_3_ref"`
	Fragment4Ref int32   `yaml:"fragment_4_ref"` // unknown, usually ref to first texture
	Center       Vector3 `yaml:"center"`
	Params2      UIndex3 `yaml:"params_2"`

	MaxDistance float32 `yaml:"max_distance"`
	Min         Vector3 `yaml:"min"`
	Max         Vector3 `yaml:"max"`
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
	Colors            []RGBA                        `yaml:"colors"`
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

	enc.Uint32(e.MaterialListRef)
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
	enc.Uint16(uint16(len(e.Vertices)))
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

	for _, meshOp := range e.MeshOps {
		enc.Uint16(meshOp.Index1)
		enc.Uint16(meshOp.Index2)
		enc.Float32(meshOp.Offset)
		enc.Uint8(meshOp.Param1)
		enc.Uint8(meshOp.TypeField)
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragDmSpriteDef2) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32() // flags, currently unknown, zone meshes are 0x00018003, placeable objects are 0x00014003

	e.MaterialListRef = dec.Uint32()
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

	meshOpCount := dec.Uint16() // number of entries in meshops. Seems to be used only for animated mob models.
	e.RawScale = dec.Uint16()

	// convert scale back to rawscale
	//rawScale = uint16(math.Log2(float64(1 / scale)))

	/// Vertices (x, y, z) belonging to this mesh. Each axis should
	/// be multiplied by (1 shl `scale`) for the final vertex position.
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
		color := RGBA{
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

// WldFragDmTrackDef2 is DmTrackDef2 in libeq, Mesh Animated Vertices in openzone, DMTRACKDEF in wld, MeshAnimatedVertices in lantern
type WldFragDmTrackDef2 struct {
	FragName    string                    `yaml:"frag_name"`
	NameRef     int32                     `yaml:"name_ref"`
	Flags       uint32                    `yaml:"flags"`
	VertexCount uint16                    `yaml:"vertex_count"`
	FrameCount  uint16                    `yaml:"frame_count"`
	Param1      uint16                    `yaml:"param_1"` // usually contains 100
	Param2      uint16                    `yaml:"param_2"` // usually contains 0
	Scale       uint16                    `yaml:"scale"`
	Frames      []WldFragMeshAnimatedBone `yaml:"frames"`
	Size6       uint32                    `yaml:"size_6"`
}

type WldFragMeshAnimatedBone struct {
	Position Vector3 `yaml:"position"`
}

func (e *WldFragDmTrackDef2) FragCode() int {
	return FragCodeDmTrackDef2
}

func (e *WldFragDmTrackDef2) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDmTrackDef2) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	return nil
}

// WldFragWorldTree is WorldTree in libeq, BSP Tree in openzone, WORLDTREE in wld, BspTree in lantern
// For serialization, refer to here: https://github.com/knervous/LanternExtractor2/blob/knervous/merged/LanternExtractor/EQ/Wld/DataTypes/BspNode.cs
// For constructing, refer to here: https://github.com/knervous/LanternExtractor2/blob/920541d15958e90aa91f7446a74226cbf26b829a/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs#L304
type WldFragWorldTree struct {
	FragName  string          `yaml:"frag_name"`
	NameRef   int32           `yaml:"name_ref"`
	NodeCount uint32          `yaml:"node_count"`
	Nodes     []WorldTreeNode `yaml:"nodes"`
}

type WorldTreeNode struct {
	Normal    Vector3 `yaml:"normal"`
	Distance  float32 `yaml:"distance"`
	RegionRef int32   `yaml:"region_ref"`
	FrontRef  int32   `yaml:"front_ref"`
	BackRef   int32   `yaml:"back_ref"`
}

func (e *WldFragWorldTree) FragCode() int {
	return FragCodeWorldTree
}

func (e *WldFragWorldTree) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.NodeCount)
	for _, node := range e.Nodes {
		enc.Float32(node.Normal.X)
		enc.Float32(node.Normal.Y)
		enc.Float32(node.Normal.Z)
		enc.Float32(node.Distance)
		enc.Int32(node.RegionRef)
		enc.Int32(node.FrontRef)
		enc.Int32(node.BackRef)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragWorldTree) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.NodeCount = dec.Uint32()
	for i := uint32(0); i < e.NodeCount; i++ {
		node := WorldTreeNode{}
		node.Normal.X = dec.Float32()
		node.Normal.Y = dec.Float32()
		node.Normal.Z = dec.Float32()
		node.Distance = dec.Float32()
		node.RegionRef = dec.Int32()
		node.FrontRef = dec.Int32()
		node.BackRef = dec.Int32()
		e.Nodes = append(e.Nodes, node)
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragRegion is Region in libeq, Bsp WldFragRegion in openzone, REGION in wld, BspRegion in lantern
type WldFragRegion struct {
	FragName             string    `yaml:"frag_name"`
	NameRef              int32     `yaml:"name_ref"`
	Flags                uint32    `yaml:"flags"`
	AmbientLightRef      int32     `yaml:"ambient_light_ref"`
	RegionVertexCount    uint32    `yaml:"region_vertex_count"`
	RegionProximalCount  uint32    `yaml:"region_proximal_count"`
	RenderVertexCount    uint32    `yaml:"render_vertex_count"`
	WallCount            uint32    `yaml:"wall_count"`
	ObstacleCount        uint32    `yaml:"obstacle_count"`
	CuttingObstacleCount uint32    `yaml:"cutting_obstacle_count"`
	VisibleNodeCount     uint32    `yaml:"visible_node_count"`
	RegionVertices       []Vector3 `yaml:"region_vertices"`
	RegionProximals      []Vector2 `yaml:"region_proximals"`
	RenderVertices       []Vector3 `yaml:"render_vertices"`
	Walls                []Wall    `yaml:"walls"`
}

type Wall struct {
	Flags                       uint32    `yaml:"flags"`
	VertexCount                 uint32    `yaml:"vertex_count"`
	RenderMethod                uint32    `yaml:"render_method"`
	RenderFlags                 uint32    `yaml:"render_flags"`
	RenderPen                   uint32    `yaml:"render_pen"`
	RenderBrightness            float32   `yaml:"render_brightness"`
	RenderScaledAmbient         float32   `yaml:"render_scaled_ambient"`
	RenderSimpleSpriteReference uint32    `yaml:"render_simple_sprite_reference"`
	RenderUVInfoOrigin          Vector3   `yaml:"render_uv_info_origin"`
	RenderUVInfoUAxis           Vector3   `yaml:"render_uv_info_u_axis"`
	RenderUVInfoVAxis           Vector3   `yaml:"render_uv_info_v_axis"`
	RenderUVMapEntryCount       uint32    `yaml:"render_uv_map_entry_count"`
	RenderUVMapEntries          []Vector2 `yaml:"render_uv_map_entries"`
	Normal                      Quad4     `yaml:"normal"`
	Vertices                    []uint32  `yaml:"vertices"`
}

func (e *WldFragRegion) FragCode() int {
	return FragCodeRegion
}

func (e *WldFragRegion) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Int32(e.AmbientLightRef)
	enc.Uint32(e.RegionVertexCount)
	enc.Uint32(e.RegionProximalCount)
	enc.Uint32(e.RenderVertexCount)
	enc.Uint32(e.WallCount)
	enc.Uint32(e.ObstacleCount)
	enc.Uint32(e.CuttingObstacleCount)
	enc.Uint32(e.VisibleNodeCount)
	for _, regionVertex := range e.RegionVertices {
		enc.Float32(regionVertex.X)
		enc.Float32(regionVertex.Y)
		enc.Float32(regionVertex.Z)
	}
	for _, regionProximal := range e.RegionProximals {
		enc.Float32(regionProximal.X)
		enc.Float32(regionProximal.Y)
	}
	for _, renderVertex := range e.RenderVertices {
		enc.Float32(renderVertex.X)
		enc.Float32(renderVertex.Y)
		enc.Float32(renderVertex.Z)
	}
	for _, wall := range e.Walls {
		enc.Uint32(wall.Flags)
		enc.Uint32(wall.VertexCount)
		enc.Uint32(wall.RenderMethod)
		enc.Uint32(wall.RenderFlags)
		enc.Uint32(wall.RenderPen)
		enc.Float32(wall.RenderBrightness)
		enc.Float32(wall.RenderScaledAmbient)
		enc.Uint32(wall.RenderSimpleSpriteReference)
		enc.Float32(wall.RenderUVInfoOrigin.X)
		enc.Float32(wall.RenderUVInfoOrigin.Y)
		enc.Float32(wall.RenderUVInfoOrigin.Z)
		enc.Float32(wall.RenderUVInfoUAxis.X)
		enc.Float32(wall.RenderUVInfoUAxis.Y)
		enc.Float32(wall.RenderUVInfoUAxis.Z)
		enc.Float32(wall.RenderUVInfoVAxis.X)
		enc.Float32(wall.RenderUVInfoVAxis.Y)
		enc.Float32(wall.RenderUVInfoVAxis.Z)
		enc.Uint32(wall.RenderUVMapEntryCount)
		for _, renderUVMapEntry := range wall.RenderUVMapEntries {
			enc.Float32(renderUVMapEntry.X)
			enc.Float32(renderUVMapEntry.Y)
		}
		enc.Float32(wall.Normal.X)
		enc.Float32(wall.Normal.Y)
		enc.Float32(wall.Normal.Z)
		enc.Float32(wall.Normal.W)
		for _, vertex := range wall.Vertices {
			enc.Uint32(vertex)
		}
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragRegion) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.AmbientLightRef = dec.Int32()
	e.RegionVertexCount = dec.Uint32()
	e.RegionProximalCount = dec.Uint32()
	e.RenderVertexCount = dec.Uint32()
	e.WallCount = dec.Uint32()
	e.ObstacleCount = dec.Uint32()
	e.CuttingObstacleCount = dec.Uint32()
	e.VisibleNodeCount = dec.Uint32()
	e.RegionVertices = make([]Vector3, e.RegionVertexCount)
	for i := uint32(0); i < e.RegionVertexCount; i++ {
		e.RegionVertices[i] = Vector3{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
		}
	}
	e.RegionProximals = make([]Vector2, e.RegionProximalCount)
	for i := uint32(0); i < e.RegionProximalCount; i++ {
		e.RegionProximals[i] = Vector2{
			X: dec.Float32(),
			Y: dec.Float32(),
		}
	}
	if e.WallCount != 0 {
		e.RenderVertexCount = 0
	}

	e.RenderVertices = make([]Vector3, e.RenderVertexCount)
	for i := uint32(0); i < e.RenderVertexCount; i++ {
		e.RenderVertices[i] = Vector3{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
		}
	}

	e.Walls = make([]Wall, e.WallCount)
	for i := uint32(0); i < e.WallCount; i++ {
		wall := Wall{}
		wall.Flags = dec.Uint32()
		wall.VertexCount = dec.Uint32()
		wall.RenderMethod = dec.Uint32()
		wall.RenderFlags = dec.Uint32()
		wall.RenderPen = dec.Uint32()
		wall.RenderBrightness = dec.Float32()
		wall.RenderScaledAmbient = dec.Float32()
		wall.RenderSimpleSpriteReference = dec.Uint32()
		wall.RenderUVInfoOrigin.X = dec.Float32()
		wall.RenderUVInfoOrigin.Y = dec.Float32()
		wall.RenderUVInfoOrigin.Z = dec.Float32()
		wall.RenderUVInfoUAxis.X = dec.Float32()
		wall.RenderUVInfoUAxis.Y = dec.Float32()
		wall.RenderUVInfoUAxis.Z = dec.Float32()
		wall.RenderUVInfoVAxis.X = dec.Float32()
		wall.RenderUVInfoVAxis.Y = dec.Float32()
		wall.RenderUVInfoVAxis.Z = dec.Float32()
		wall.RenderUVMapEntryCount = dec.Uint32()
		for i := uint32(0); i < wall.RenderUVMapEntryCount; i++ {
			wall.RenderUVMapEntries = append(wall.RenderUVMapEntries, Vector2{
				X: dec.Float32(),
				Y: dec.Float32(),
			})
		}
		wall.Normal.X = dec.Float32()
		wall.Normal.Y = dec.Float32()
		wall.Normal.Z = dec.Float32()
		wall.Normal.W = dec.Float32()
		wall.Vertices = make([]uint32, wall.VertexCount)
		for i := uint32(0); i < wall.VertexCount; i++ {
			wall.Vertices[i] = dec.Uint32()
		}
		e.Walls[i] = wall
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
