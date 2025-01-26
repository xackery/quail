package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragParticleSpriteDef is ParticleSpriteDef in libeq, empty in openzone, PARTICLESPRITEDEF in wld
type WldFragParticleSpriteDef struct {
	nameRef                     int32
	Flags                       uint32
	VerticesCount               uint32
	Unknown                     uint32
	CenterOffset                [3]float32
	Radius                      float32
	Vertices                    [][3]float32
	RenderMethod                uint32
	RenderFlags                 uint32
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          [3]float32
	RenderUVInfoUAxis           [3]float32
	RenderUVInfoVAxis           [3]float32
	RenderUVMapEntryCount       uint32
	RenderUVMapEntries          [][2]float32
	Pen                         []uint32
}

func (e *WldFragParticleSpriteDef) FragCode() int {
	return FragCodeParticleSpriteDef
}

func (e *WldFragParticleSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.VerticesCount)
	enc.Uint32(e.Unknown)
	if e.Flags&0x01 != 0 { // has center offset
		enc.Float32(e.CenterOffset[0])
		enc.Float32(e.CenterOffset[1])
		enc.Float32(e.CenterOffset[2])
	}
	if e.Flags&0x02 != 0 { // has radius
		enc.Float32(e.Radius)
	}
	if e.VerticesCount > 0 { // has vertices
		for _, vertex := range e.Vertices {
			enc.Float32(vertex[0])
			enc.Float32(vertex[1])
			enc.Float32(vertex[2])
		}
	}
	enc.Uint32(e.RenderMethod)
	enc.Uint32(e.RenderFlags)
	enc.Uint32(e.RenderPen)
	enc.Float32(e.RenderBrightness)
	enc.Float32(e.RenderScaledAmbient)
	enc.Uint32(e.RenderSimpleSpriteReference)
	enc.Float32(e.RenderUVInfoOrigin[0])
	enc.Float32(e.RenderUVInfoOrigin[1])
	enc.Float32(e.RenderUVInfoOrigin[2])
	enc.Float32(e.RenderUVInfoUAxis[0])
	enc.Float32(e.RenderUVInfoUAxis[1])
	enc.Float32(e.RenderUVInfoUAxis[2])
	enc.Float32(e.RenderUVInfoVAxis[0])
	enc.Float32(e.RenderUVInfoVAxis[1])
	enc.Float32(e.RenderUVInfoVAxis[2])
	enc.Uint32(e.RenderUVMapEntryCount)
	for _, entry := range e.RenderUVMapEntries {
		enc.Float32(entry[0])
		enc.Float32(entry[1])
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragParticleSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.VerticesCount = dec.Uint32()
	e.Unknown = dec.Uint32()
	if e.Flags&0x01 != 0 { // has center offset
		e.CenterOffset[0] = dec.Float32()
		e.CenterOffset[1] = dec.Float32()
		e.CenterOffset[2] = dec.Float32()
	}
	if e.Flags&0x02 != 0 { // has radius
		e.Radius = dec.Float32()
	}
	if e.VerticesCount > 0 { // has vertices
		for i := uint32(0); i < e.VerticesCount; i++ {
			var vertex [3]float32
			vertex[0] = dec.Float32()
			vertex[1] = dec.Float32()
			vertex[2] = dec.Float32()
			e.Vertices = append(e.Vertices, vertex)
		}
	}
	e.RenderMethod = dec.Uint32()
	e.RenderFlags = dec.Uint32()
	e.RenderPen = dec.Uint32()
	e.RenderBrightness = dec.Float32()
	e.RenderScaledAmbient = dec.Float32()
	e.RenderSimpleSpriteReference = dec.Uint32()
	e.RenderUVInfoOrigin[0] = dec.Float32()
	e.RenderUVInfoOrigin[1] = dec.Float32()
	e.RenderUVInfoOrigin[2] = dec.Float32()
	e.RenderUVInfoUAxis[0] = dec.Float32()
	e.RenderUVInfoUAxis[1] = dec.Float32()
	e.RenderUVInfoUAxis[2] = dec.Float32()
	e.RenderUVInfoVAxis[0] = dec.Float32()
	e.RenderUVInfoVAxis[1] = dec.Float32()
	e.RenderUVInfoVAxis[2] = dec.Float32()
	e.RenderUVMapEntryCount = dec.Uint32()
	for i := uint32(0); i < e.RenderUVMapEntryCount; i++ {
		var entry [2]float32
		entry[0] = dec.Float32()
		entry[1] = dec.Float32()
		e.RenderUVMapEntries = append(e.RenderUVMapEntries, entry)
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragParticleSpriteDef) NameRef() int32 {
	return e.nameRef
}
