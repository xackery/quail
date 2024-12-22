package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSprite2DDef is Sprite2DDef in libeq, Two-Dimensional Object in openzone, 2DSPRITEDEF in wld, Fragment06 in lantern
type WldFragSprite2DDef struct {
	NameRef                     int32
	Flags                       uint32
	TextureCount                uint32
	PitchCount                  uint32
	Scale                       [2]float32
	SphereListRef               uint32
	DepthScale                  float32
	CenterOffset                [3]float32
	BoundingRadius              float32
	CurrentFrameRef             int32
	Sleep                       uint32
	Pitches                     []Pitch
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

type Pitch struct {
	PitchCap     int32
	Flag         bool
	HeadingCount uint32
	Headings     []Heading
}

type Heading struct {
	HeadingCap int32
	Frames     []int32
}

func (e *WldFragSprite2DDef) FragCode() int {
	return FragCodeSprite2DDef
}

func (e *WldFragSprite2DDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.TextureCount)
	enc.Uint32(e.PitchCount)

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

		// Combine Flag and HeadingCount into one DWORD
		rawFlag := uint32(0)
		if pitch.Flag {
			rawFlag |= 0x80000000 // Set MSB if Flag is true
		}
		rawFlag |= uint32(pitch.HeadingCount)
		enc.Uint32(rawFlag)

		for _, heading := range pitch.Headings {
			enc.Int32(heading.HeadingCap)
			for _, frame := range heading.Frames {
				enc.Int32(frame)
			}
		}
	}
	if e.Flags&0x10 == 0x10 {
		enc.Uint32(e.RenderMethod)
		enc.Uint8(e.RenderFlags)

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
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.TextureCount = dec.Uint32()
	e.PitchCount = dec.Uint32()
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
	e.Pitches = make([]Pitch, e.PitchCount)
	for i := uint32(0); i < e.PitchCount; i++ {
		var pitch Pitch
		pitch.PitchCap = dec.Int32()

		// Extract the most significant bit as Flag and the rest as HeadingCount
		rawFlag := dec.Uint32()
		pitch.Flag = (rawFlag & 0x80000000) != 0 // MSB
		pitch.HeadingCount = uint32(rawFlag & 0x7FFFFFFF)

		// Read Headings
		pitch.Headings = make([]Heading, pitch.HeadingCount)
		for j := uint32(0); j < pitch.HeadingCount; j++ {
			var heading Heading
			heading.HeadingCap = dec.Int32()
			heading.Frames = make([]int32, e.TextureCount)
			for k := uint32(0); k < e.TextureCount; k++ {
				heading.Frames[k] = dec.Int32()
			}
			pitch.Headings[j] = heading
		}
		e.Pitches[i] = pitch
	}
	if e.Flags&0x10 == 0x10 {
		e.RenderMethod = dec.Uint32()
		e.RenderFlags = dec.Uint8()

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
