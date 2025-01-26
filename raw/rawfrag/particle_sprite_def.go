package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragParticleSpriteDef is ParticleSpriteDef in libeq, empty in openzone, PARTICLESPRITEDEF in wld
type WldFragParticleSpriteDef struct {
	nameRef                     int32
	Flags                       uint32
	VerticesCount               uint32
	Unknown                     uint32
	CenterOffset                model.Vector3
	Radius                      float32
	Vertices                    []model.Vector3
	RenderMethod                uint32
	RenderFlags                 uint32
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          model.Vector3
	RenderUVInfoUAxis           model.Vector3
	RenderUVInfoVAxis           model.Vector3
	RenderUVMapEntryCount       uint32
	RenderUVMapEntries          []model.Vector2
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
		enc.Float32(e.CenterOffset.X)
		enc.Float32(e.CenterOffset.Y)
		enc.Float32(e.CenterOffset.Z)
	}
	if e.Flags&0x02 != 0 { // has radius
		enc.Float32(e.Radius)
	}
	if e.VerticesCount > 0 { // has vertices
		for _, vertex := range e.Vertices {
			enc.Float32(vertex.X)
			enc.Float32(vertex.Y)
			enc.Float32(vertex.Z)
		}
	}
	enc.Uint32(e.RenderMethod)
	enc.Uint32(e.RenderFlags)
	enc.Uint32(e.RenderPen)
	enc.Float32(e.RenderBrightness)
	enc.Float32(e.RenderScaledAmbient)
	enc.Uint32(e.RenderSimpleSpriteReference)
	enc.Float32(e.RenderUVInfoOrigin.X)
	enc.Float32(e.RenderUVInfoOrigin.Y)
	enc.Float32(e.RenderUVInfoOrigin.Z)
	enc.Float32(e.RenderUVInfoUAxis.X)
	enc.Float32(e.RenderUVInfoUAxis.Y)
	enc.Float32(e.RenderUVInfoUAxis.Z)
	enc.Float32(e.RenderUVInfoVAxis.X)
	enc.Float32(e.RenderUVInfoVAxis.Y)
	enc.Float32(e.RenderUVInfoVAxis.Z)
	enc.Uint32(e.RenderUVMapEntryCount)
	for _, entry := range e.RenderUVMapEntries {
		enc.Float32(entry.X)
		enc.Float32(entry.Y)
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
		e.CenterOffset.X = dec.Float32()
		e.CenterOffset.Y = dec.Float32()
		e.CenterOffset.Z = dec.Float32()
	}
	if e.Flags&0x02 != 0 { // has radius
		e.Radius = dec.Float32()
	}
	if e.VerticesCount > 0 { // has vertices
		for i := uint32(0); i < e.VerticesCount; i++ {
			var vertex model.Vector3
			vertex.X = dec.Float32()
			vertex.Y = dec.Float32()
			vertex.Z = dec.Float32()
			e.Vertices = append(e.Vertices, vertex)
		}
	}
	e.RenderMethod = dec.Uint32()
	e.RenderFlags = dec.Uint32()
	e.RenderPen = dec.Uint32()
	e.RenderBrightness = dec.Float32()
	e.RenderScaledAmbient = dec.Float32()
	e.RenderSimpleSpriteReference = dec.Uint32()
	e.RenderUVInfoOrigin.X = dec.Float32()
	e.RenderUVInfoOrigin.Y = dec.Float32()
	e.RenderUVInfoOrigin.Z = dec.Float32()
	e.RenderUVInfoUAxis.X = dec.Float32()
	e.RenderUVInfoUAxis.Y = dec.Float32()
	e.RenderUVInfoUAxis.Z = dec.Float32()
	e.RenderUVInfoVAxis.X = dec.Float32()
	e.RenderUVInfoVAxis.Y = dec.Float32()
	e.RenderUVInfoVAxis.Z = dec.Float32()
	e.RenderUVMapEntryCount = dec.Uint32()
	for i := uint32(0); i < e.RenderUVMapEntryCount; i++ {
		var entry model.Vector2
		entry.X = dec.Float32()
		entry.Y = dec.Float32()
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
