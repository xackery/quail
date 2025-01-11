package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSprite2DDef is Sprite2DDef in libeq, Two-Dimensional Object in openzone, 2DSPRITEDEF in wld, Fragment06 in lantern
type WldFragSprite2DDef struct {
	nameRef                     int32
	Flags                       uint32
	Scale                       [2]float32
	SphereListRef               uint32
	DepthScale                  float32
	CenterOffset                [3]float32
	BoundingRadius              float32
	CurrentFrameRef             int32
	Sleep                       uint32
	Pitches                     []*WldFragSprite2DPitch
	RenderMethod                uint32
	RenderFlags                 uint32
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          [3]float32
	RenderUVInfoUAxis           [3]float32
	RenderUVInfoVAxis           [3]float32
	Uvs                         [][2]float32
}

type WldFragSprite2DPitch struct {
	PitchCap        int32
	TopOrBottomView uint32
	Headings        []*WldFragSprite2DHeading
}

type WldFragSprite2DHeading struct {
	HeadingCap int32
	FrameRefs  []int32
}

func (e *WldFragSprite2DDef) FragCode() int {
	return FragCodeSprite2DDef
}

func (e *WldFragSprite2DDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)

	textureCount := uint32(0)
	if len(e.Pitches) < 1 {
		return fmt.Errorf("no pitches found")
	}
	if len(e.Pitches[0].Headings) < 1 {
		return fmt.Errorf("no headings found")
	}
	if len(e.Pitches[0].Headings[0].FrameRefs) < 1 {
		return fmt.Errorf("no frame refs found")
	}
	textureCount = uint32(len(e.Pitches[0].Headings[0].FrameRefs))
	enc.Uint32(textureCount)
	enc.Uint32(uint32(len(e.Pitches)))

	enc.Float32(e.Scale[0])
	enc.Float32(e.Scale[1])
	enc.Uint32(e.SphereListRef)
	if e.Flags&0x80 == 0x80 {
		enc.Float32(e.DepthScale)
	}
	if e.Flags&0x01 == 0x01 {
		enc.Float32(e.CenterOffset[0])
		enc.Float32(e.CenterOffset[1])
		enc.Float32(e.CenterOffset[2])
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

	for _, pitch := range e.Pitches {
		enc.Int32(pitch.PitchCap)
		enc.Uint32((uint32(pitch.TopOrBottomView) << 31) | (uint32(len(pitch.Headings)) & 0x7FFFFFFF))
		for _, heading := range pitch.Headings {
			enc.Int32(heading.HeadingCap)
			for _, frameRef := range heading.FrameRefs {
				enc.Int32(frameRef)
			}
		}
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
			enc.Float32(e.RenderUVInfoOrigin[0])
			enc.Float32(e.RenderUVInfoOrigin[1])
			enc.Float32(e.RenderUVInfoOrigin[2])
			enc.Float32(e.RenderUVInfoUAxis[0])
			enc.Float32(e.RenderUVInfoUAxis[1])
			enc.Float32(e.RenderUVInfoUAxis[2])
			enc.Float32(e.RenderUVInfoVAxis[0])
			enc.Float32(e.RenderUVInfoVAxis[1])
			enc.Float32(e.RenderUVInfoVAxis[2])
		}
		if e.RenderFlags&0x20 == 0x20 {
			enc.Uint32(uint32(len(e.Uvs)))
			for _, uv := range e.Uvs {
				enc.Float32(uv[0])
				enc.Float32(uv[1])
			}
		}
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil

}

func (e *WldFragSprite2DDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	textureCount := dec.Uint32()
	pitchCount := dec.Uint32()
	e.Scale[0] = dec.Float32()
	e.Scale[1] = dec.Float32()
	e.SphereListRef = dec.Uint32()
	if e.Flags&0x80 == 0x80 {
		e.DepthScale = dec.Float32()
	}
	if e.Flags&0x01 == 0x01 {
		e.CenterOffset[0] = dec.Float32()
		e.CenterOffset[1] = dec.Float32()
		e.CenterOffset[2] = dec.Float32()
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
	e.Pitches = []*WldFragSprite2DPitch{}
	for i := uint32(0); i < pitchCount; i++ {
		pitch := &WldFragSprite2DPitch{
			PitchCap: dec.Int32(),
		}
		weirdFlagCount := dec.Uint32()
		pitch.TopOrBottomView = uint32((weirdFlagCount >> 31) & 0x1)
		headingCount := uint32(weirdFlagCount & 0x7FFFFFFF)

		pitch.Headings = []*WldFragSprite2DHeading{}
		for j := uint32(0); j < headingCount; j++ {
			heading := &WldFragSprite2DHeading{
				HeadingCap: dec.Int32(),
			}

			heading.FrameRefs = make([]int32, textureCount)
			for k := uint32(0); k < textureCount; k++ {
				heading.FrameRefs[k] = dec.Int32()
			}

			pitch.Headings = append(pitch.Headings, heading)
		}

		e.Pitches = append(e.Pitches, pitch)
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
			e.RenderUVInfoOrigin[0] = dec.Float32()
			e.RenderUVInfoOrigin[1] = dec.Float32()
			e.RenderUVInfoOrigin[2] = dec.Float32()
			e.RenderUVInfoUAxis[0] = dec.Float32()
			e.RenderUVInfoUAxis[1] = dec.Float32()
			e.RenderUVInfoUAxis[2] = dec.Float32()
			e.RenderUVInfoVAxis[0] = dec.Float32()
			e.RenderUVInfoVAxis[1] = dec.Float32()
			e.RenderUVInfoVAxis[2] = dec.Float32()
		}
		if e.RenderFlags&0x20 == 0x20 {
			renderUVMapEntryCount := dec.Uint32()
			for j := 0; j < int(renderUVMapEntryCount); j++ {
				u := [2]float32{dec.Float32(), dec.Float32()}
				e.Uvs = append(e.Uvs, u)
			}
		}
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil

}

func (e *WldFragSprite2DDef) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragSprite2DDef) SetNameRef(nameRef int32) {
	e.nameRef = nameRef
}
