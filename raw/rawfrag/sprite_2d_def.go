package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragSprite2DDef is Sprite2DDef in libeq, Two-Dimensional Object in openzone, 2DSPRITEDEF in wld, Fragment06 in lantern
type WldFragSprite2DDef struct {
	NameRef                     int32           `yaml:"name_ref"`
	Flags                       uint32          `yaml:"flags"`
	TextureCount                uint32          `yaml:"texture_count"`
	PitchCount                  uint32          `yaml:"pitch_count"`
	Scale                       model.Vector2   `yaml:"scale"`
	SphereRef                   uint32          `yaml:"sphere_ref"`
	DepthScale                  float32         `yaml:"depth_scale"`
	CenterOffset                model.Vector3   `yaml:"center_offset"`
	BoundingRadius              float32         `yaml:"bounding_radius"`
	CurrentFrameRef             int32           `yaml:"current_frame_ref"`
	Sleep                       uint32          `yaml:"sleep"`
	Headings                    []uint32        `yaml:"headings"`
	RenderMethod                uint32          `yaml:"render_method"`
	RenderFlags                 uint32          `yaml:"render_flags"`
	RenderPen                   uint32          `yaml:"render_pen"`
	RenderBrightness            float32         `yaml:"render_brightness"`
	RenderScaledAmbient         float32         `yaml:"render_scaled_ambient"`
	RenderSimpleSpriteReference uint32          `yaml:"render_simple_sprite_reference"`
	RenderUVInfoOrigin          model.Vector3   `yaml:"render_uv_info_origin"`
	RenderUVInfoUAxis           model.Vector3   `yaml:"render_uv_info_u_axis"`
	RenderUVInfoVAxis           model.Vector3   `yaml:"render_uv_info_v_axis"`
	RenderUVMapEntries          []model.Vector2 `yaml:"render_uv_map_entries"`
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
				v := model.Vector2{}
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
