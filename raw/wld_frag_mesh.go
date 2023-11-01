package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragTwoDSprite is Sprite2DDef in libeq, Two-Dimensional Object in openzone, 2DSPRITEDEF in wld, Fragment06 in lantern
type WldFragTwoDSprite struct {
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

func (e *WldFragTwoDSprite) FragCode() int {
	return 0x06
}

func (e *WldFragTwoDSprite) Encode(w io.Writer) error {
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
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil

}

func decodeTwoDSprite(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragTwoDSprite{}
	d.FragName = FragName(d.FragCode())
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
			renderUVMapEntrycount := dec.Uint32()
			for i := uint32(0); i < renderUVMapEntrycount; i++ {
				v := Vector2{}
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

// WldFragTwoDSpriteRef is Sprite2D in libeq, Two-Dimensional Object Reference in openzone, 2DSPRITE (ref) in wld, Fragment07 in lantern
type WldFragTwoDSpriteRef struct {
	FragName      string `yaml:"frag_name"`
	NameRef       int32  `yaml:"name_ref"`
	TwoDSpriteRef uint32 `yaml:"two_d_sprite_ref"`
	Flags         uint32 `yaml:"flags"`
}

func (e *WldFragTwoDSpriteRef) FragCode() int {
	return 0x07
}

func (e *WldFragTwoDSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.TwoDSpriteRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeTwoDSpriteRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragTwoDSpriteRef{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.TwoDSpriteRef = dec.Uint32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragThreeDSprite is Sprite3DDef in libeq, Camera in openzone, 3DSPRITEDEF in wld, Camera in lantern
type WldFragThreeDSprite struct {
	FragName      string    `yaml:"frag_name"`
	NameRef       int32     `yaml:"name_ref"`
	Flags         uint32    `yaml:"flags"`
	SphereListRef uint32    `yaml:"sphere_list_ref"`
	CenterOffset  Vector3   `yaml:"center_offset"`
	Radius        float32   `yaml:"radius"`
	Vertices      []Vector3 `yaml:"vertices"`
	BspNodes      []BspNode `yaml:"bsp_nodes"`
}

type BspNode struct {
	FrontTree                   uint32    `yaml:"front_tree"`
	BackTree                    uint32    `yaml:"back_tree"`
	VertexIndexes               []uint32  `yaml:"vertex_indexes"`
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

func (e *WldFragThreeDSprite) FragCode() int {
	return 0x08
}

func (e *WldFragThreeDSprite) Encode(w io.Writer) error {
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
		enc.Uint32(node.RenderFlags)

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

func decodeThreeDSprite(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragThreeDSprite{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	vertexCount := dec.Uint32()
	bspNodeCount := dec.Uint32()
	d.SphereListRef = dec.Uint32()
	d.CenterOffset.X = dec.Float32()
	d.CenterOffset.Y = dec.Float32()
	d.CenterOffset.Z = dec.Float32()
	d.Radius = dec.Float32()
	for i := 0; i < int(vertexCount); i++ {
		v := Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		d.Vertices = append(d.Vertices, v)
	}
	for i := 0; i < int(bspNodeCount); i++ {
		node := BspNode{}
		vertexIndexCount := dec.Uint32()
		node.FrontTree = dec.Uint32()
		node.BackTree = dec.Uint32()
		for j := 0; j < int(vertexIndexCount); j++ {
			node.VertexIndexes = append(node.VertexIndexes, dec.Uint32())
		}
		node.RenderMethod = dec.Uint32()
		node.RenderFlags = dec.Uint32()

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
				v := Vector2{}
				v.X = dec.Float32()
				v.Y = dec.Float32()
				node.RenderUVMapEntries = append(node.RenderUVMapEntries, v)
			}
		}
		d.BspNodes = append(d.BspNodes, node)
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}

	return d, nil
}

// WldFragThreeDSpriteRef is Sprite3D in libeq, Camera Reference in openzone, 3DSPRITE (ref) in wld, CameraReference in lantern
type WldFragThreeDSpriteRef struct {
	FragName  string `yaml:"frag_name"`
	NameRef   int32
	ThreeDRef int32
	Flags     uint32
}

func (e *WldFragThreeDSpriteRef) FragCode() int {
	return 0x09
}

func (e *WldFragThreeDSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.ThreeDRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeThreeDSpriteRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragThreeDSpriteRef{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.ThreeDRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragFourDSprite is Sprite4DDef in libeq, empty in openzone, 4DSPRITEDEF in wld
type WldFragFourDSprite struct {
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

func (e *WldFragFourDSprite) FragCode() int {
	return 0x0A
}

func (e *WldFragFourDSprite) Encode(w io.Writer) error {
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
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeFourDSprite(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragFourDSprite{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	frameCount := dec.Uint32()
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
		for i := uint32(0); i < frameCount; i++ {
			d.SpriteFragments = append(d.SpriteFragments, dec.Uint32())
		}
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragFourDSpriteRef is Sprite4D in libeq, empty in openzone, 4DSPRITE (ref) in wld
type WldFragFourDSpriteRef struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32  `yaml:"name_ref"`
	FourDRef int32  `yaml:"four_d_ref"`
	Params1  uint32 `yaml:"params_1"`
}

func (e *WldFragFourDSpriteRef) FragCode() int {
	return 0x0B
}

func (e *WldFragFourDSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.FourDRef)
	enc.Uint32(e.Params1)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeFourDSpriteRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragFourDSpriteRef{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.FourDRef = dec.Int32()
	d.Params1 = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragPolyhedron is PolyhedronDef in libeq, Polygon animation in openzone, POLYHEDRONDEFINITION in wld, Fragment17 in lantern
type WldFragPolyhedron struct {
	FragName string     `yaml:"frag_name"`
	NameRef  int32      `yaml:"name_ref"`
	Flags    uint32     `yaml:"flags"`
	Size1    uint32     `yaml:"size_1"`
	Size2    uint32     `yaml:"size_2"`
	Params1  float32    `yaml:"params_1"`
	Params2  float32    `yaml:"params_2"`
	Entries1 []Vector3  `yaml:"entries_1"`
	Entries2 []Entries2 `yaml:"entries_2"`
}

type Entries2 struct {
	Unk1 uint32   `yaml:"unk_1"`
	Unk2 []uint32 `yaml:"unk_2"`
}

func (e *WldFragPolyhedron) FragCode() int {
	return 0x17
}

func (e *WldFragPolyhedron) Encode(w io.Writer) error {
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

func decodePolyhedron(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragPolyhedron{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.Size1 = dec.Uint32()
	d.Size2 = dec.Uint32()
	d.Params1 = dec.Float32()
	d.Params2 = dec.Float32()
	for i := uint32(0); i < d.Size1; i++ {
		v := Vector3{}
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

// WldFragPolyhedronRef is Polyhedron in libeq, Polygon Animation Reference in openzone, POLYHEDRON (ref) in wld, Fragment18 in lantern
type WldFragPolyhedronRef struct {
	FragName    string  `yaml:"frag_name"`
	NameRef     int32   `yaml:"name_ref"`
	FragmentRef int32   `yaml:"fragment_ref"`
	Flags       uint32  `yaml:"flags"`
	Scale       float32 `yaml:"scale"`
}

func (e *WldFragPolyhedronRef) FragCode() int {
	return 0x18
}

func (e *WldFragPolyhedronRef) Encode(w io.Writer) error {
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

func decodePolyhedronRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragPolyhedronRef{}
	d.FragName = FragName(d.FragCode())
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

// WldFragDMSprite is DmSpriteDef in libeq, Alternate Mesh in openzone, DMSPRITEDEF in wld, LegacyMesh in lantern
type WldFragDMSprite struct {
	FragName          string              `yaml:"frag_name"`
	NameRef           int32               `yaml:"name_ref"`
	Flags             uint32              `yaml:"flags"`
	Fragment1Maybe    int16               `yaml:"fragment_1_maybe"`
	MaterialReference uint32              `yaml:"material_reference"`
	Fragment3         uint32              `yaml:"fragment_3"`
	CenterPosition    Vector3             `yaml:"center_position"`
	Params2           uint32              `yaml:"params_2"`
	Something2        uint32              `yaml:"something_2"`
	Something3        uint32              `yaml:"something_3"`
	Verticies         []Vector3           `yaml:"verticies"`
	TexCoords         []Vector3           `yaml:"tex_coords"`
	Normals           []Vector3           `yaml:"normals"`
	Colors            []int32             `yaml:"colors"`
	Polygons          []SpritePolygon     `yaml:"polygons"`
	VertexPieces      []SpriteVertexPiece `yaml:"vertex_pieces"`
	PostVertexFlag    uint32              `yaml:"post_vertex_flag"`
	RenderGroups      []SpriteRenderGroup `yaml:"render_groups"`
	VertexTex         []Vector2           `yaml:"vertex_tex"`
	Size6Pieces       []Size6Entry        `yaml:"size_6_pieces"`
}

type SpritePolygon struct {
	Flag int16 `yaml:"flag"`
	Unk1 int16 `yaml:"unk_1"`
	Unk2 int16 `yaml:"unk_2"`
	Unk3 int16 `yaml:"unk_3"`
	Unk4 int16 `yaml:"unk_4"`
	I1   int16 `yaml:"i_1"`
	I2   int16 `yaml:"i_2"`
	I3   int16 `yaml:"i_3"`
}

type SpriteVertexPiece struct {
	Count  int16 `yaml:"count"`
	Offset int16 `yaml:"offset"`
}

type SpriteRenderGroup struct {
	PolygonCount int16 `yaml:"polygon_count"`
	MaterialId   int16 `yaml:"material_id"`
}

type Size6Entry struct {
	Unk1 uint32 `yaml:"unk_1"`
	Unk2 uint32 `yaml:"unk_2"`
	Unk3 uint32 `yaml:"unk_3"`
	Unk4 uint32 `yaml:"unk_4"`
	Unk5 uint32 `yaml:"unk_5"`
}

func (e *WldFragDMSprite) FragCode() int {
	return 0x2C
}

func (e *WldFragDMSprite) Encode(w io.Writer) error {
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

func decodeDMSprite(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragDMSprite{}
	d.FragName = FragName(d.FragCode())
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
		return nil, fmt.Errorf("vertex count misaligned (%d)", vertexCount)
	}

	if texCoordCount > 999 {
		return nil, fmt.Errorf("tex coord count misaligned (%d)", texCoordCount)
	}

	if normalCount > 999 {
		return nil, fmt.Errorf("normal count misaligned (%d)", normalCount)
	}

	for i := int16(0); i < vertexCount; i++ {
		v := Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		d.Verticies = append(d.Verticies, v)
	}

	for i := uint32(0); i < texCoordCount; i++ {
		v := Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		d.TexCoords = append(d.TexCoords, v)
	}

	for i := uint32(0); i < normalCount; i++ {
		v := Vector3{}
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
			v := Vector2{}
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

// WldFragDMSpriteRef is DmSprite in libeq, Mesh Reference in openzone, empty in wld, MeshReference in lantern
type WldFragDMSpriteRef struct {
	FragName    string `yaml:"frag_name"`
	NameRef     int32  `yaml:"name_ref"`
	DMSpriteRef int32  `yaml:"dm_sprite_ref"`
	Params      uint32 `yaml:"params"`
}

func (e *WldFragDMSpriteRef) FragCode() int {
	return 0x2D
}

func (e *WldFragDMSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.DMSpriteRef)
	enc.Uint32(e.Params)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeDMSpriteRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragDMSpriteRef{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.DMSpriteRef = dec.Int32()
	d.Params = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragMesh is DmSpriteDef2 in libeq, WldFragMesh in openzone, DMSPRITEDEF2 in wld, WldFragMesh in lantern
type WldFragMesh struct {
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
	RawScale          uint16                 `yaml:"raw_scale"`
	MeshopCount       uint16                 `yaml:"meshop_count"`
	Scale             float32                `yaml:"scale"`
	Vertices          [][3]int16             `yaml:"vertices"`
	UVs               [][2]int16             `yaml:"uvs"`
	Normals           [][3]int8              `yaml:"normals"`
	Colors            []RGBA                 `yaml:"colors"`
	Triangles         []MeshTriangleEntry    `yaml:"triangles"`
	TriangleMaterials []MeshTriangleMaterial `yaml:"triangle_materials"`
	VertexPieces      []MeshVertexPiece      `yaml:"vertex_pieces"`
	VertexMaterials   []MeshVertexPiece      `yaml:"vertex_materials"`
	MeshOps           []MeshOpEntry          `yaml:"mesh_ops"`
}

type MeshTriangleEntry struct {
	Flags uint16    `yaml:"flags"`
	Index [3]uint16 `yaml:"indexes"`
}

type MeshVertexPiece struct {
	Count  int16 `yaml:"count"`
	Index1 int16 `yaml:"index_1"`
}

type MeshTriangleMaterial struct {
	Count      uint16 `yaml:"count"`
	MaterialID uint16 `yaml:"material_id"`
}

type MeshOpEntry struct {
	Index1    uint16  `yaml:"index_1"`
	Index2    uint16  `yaml:"index_2"`
	Offset    float32 `yaml:"offset"`
	Param1    uint8   `yaml:"param_1"`
	TypeField uint8   `yaml:"type_field"`
}

func (e *WldFragMesh) FragCode() int {
	return 0x36
}

func (e *WldFragMesh) Encode(w io.Writer) error {
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

	if enc.Error() != nil {
		return enc.Error()
	}

	return nil
}

func decodeMesh(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragMesh{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32() // flags, currently unknown, zone meshes are 0x00018003, placeable objects are 0x00014003

	d.MaterialListRef = dec.Uint32()
	d.AnimationRef = dec.Int32() //used by flags/trees only

	d.Fragment3Ref = dec.Int32() // unknown, usually empty
	d.Fragment4Ref = dec.Int32() // unknown, This usually seems to reference the first [TextureImagesFragment] fragment in the file.

	d.Center.X = dec.Float32() // for zone meshes, x coordinate of the center of the mesh
	d.Center.Y = dec.Float32() // for zone meshes, x coordinate of the center of the mesh
	d.Center.Z = dec.Float32()

	d.Params2.X = dec.Uint32() // unknown, usually empty
	d.Params2.Y = dec.Uint32() // unknown, usually empty
	d.Params2.Z = dec.Uint32() // unknown, usually empty

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

	meshOpCount := dec.Uint16() // number of entries in meshops. Seems to be used only for animated mob models.
	d.RawScale = dec.Uint16()

	// convert scale back to rawscale
	//rawScale = uint16(math.Log2(float64(1 / scale)))

	/// Vertices (x, y, z) belonging to this mesh. Each axis should
	/// be multiplied by (1 shl `scale`) for the final vertex position.
	for i := 0; i < int(vertexCount); i++ {
		d.Vertices = append(d.Vertices, [3]int16{dec.Int16(), dec.Int16(), dec.Int16()})
	}

	for i := 0; i < int(uvCount); i++ {
		d.UVs = append(d.UVs, [2]int16{dec.Int16(), dec.Int16()})
	}

	for i := 0; i < int(normalCount); i++ {
		d.Normals = append(d.Normals, [3]int8{dec.Int8(), dec.Int8(), dec.Int8()})
	}

	for i := 0; i < int(colorCount); i++ {
		color := RGBA{
			R: dec.Uint8(),
			G: dec.Uint8(),
			B: dec.Uint8(),
			A: dec.Uint8(),
		}
		d.Colors = append(d.Colors, color)
	}

	for i := 0; i < int(triangleCount); i++ {
		mte := MeshTriangleEntry{}
		mte.Flags = dec.Uint16()
		mte.Index = [3]uint16{dec.Uint16(), dec.Uint16(), dec.Uint16()}

		d.Triangles = append(d.Triangles, mte)
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

	for i := 0; i < int(meshOpCount); i++ {
		d.MeshOps = append(d.MeshOps, MeshOpEntry{
			Index1:    dec.Uint16(),
			Index2:    dec.Uint16(),
			Offset:    dec.Float32(),
			Param1:    dec.Uint8(),
			TypeField: dec.Uint8(),
		})
	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}

	return d, nil
}

// WldFragMeshAnimated is DmTrackDef2 in libeq, Mesh Animated Vertices in openzone, DMTRACKDEF in wld, MeshAnimatedVertices in lantern
type WldFragMeshAnimated struct {
	FragName    string             `yaml:"frag_name"`
	NameRef     int32              `yaml:"name_ref"`
	Flags       uint32             `yaml:"flags"`
	VertexCount uint16             `yaml:"vertex_count"`
	FrameCount  uint16             `yaml:"frame_count"`
	Param1      uint16             `yaml:"param_1"` // usually contains 100
	Param2      uint16             `yaml:"param_2"` // usually contains 0
	Scale       uint16             `yaml:"scale"`
	Frames      []MeshAnimatedBone `yaml:"frames"`
	Size6       uint32             `yaml:"size_6"`
}

type MeshAnimatedBone struct {
	Position Vector3 `yaml:"position"`
}

func (e *WldFragMeshAnimated) FragCode() int {
	return 0x37
}

func (e *WldFragMeshAnimated) Encode(w io.Writer) error {
	return nil
}

func decodeMeshAnimated(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragMeshAnimated{}
	d.FragName = FragName(d.FragCode())
	return d, nil
}
