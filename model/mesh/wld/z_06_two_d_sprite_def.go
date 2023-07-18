package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// twoDSpriteDef is a 2D sprite definition 0x06
type twoDSpriteDef struct {
	nameRef         int32
	flags           uint32
	textureCount    uint32
	pitchCount      uint32
	scale           geo.Vector2
	sphereRef       uint32
	depthScale      float32
	centerOffset    geo.Vector3
	boundingRadius  float32
	currentFrameRef int32
	sleep           uint32
	headings        []uint32
	renderMethod    uint32
	renderInfo      renderInfo
}

func (e *WLD) twoDSpriteDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &twoDSpriteDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	def.textureCount = dec.Uint32()
	def.pitchCount = dec.Uint32()
	def.scale.X = dec.Float32()
	def.scale.Y = dec.Float32()
	def.sphereRef = dec.Uint32()
	if def.flags&0x80 == 0x80 {
		def.depthScale = dec.Float32()
	}
	if def.flags&0x01 == 0x01 {
		def.centerOffset.X = dec.Float32()
		def.centerOffset.Y = dec.Float32()
		def.centerOffset.Z = dec.Float32()
	}
	if def.flags&0x02 == 0x02 {
		def.boundingRadius = dec.Float32()
	}
	if def.flags&0x04 == 0x04 {
		def.currentFrameRef = dec.Int32()
	}
	if def.flags&0x08 == 0x08 {
		def.sleep = dec.Uint32()
	}
	def.headings = make([]uint32, def.pitchCount)
	for i := uint32(0); i < def.pitchCount; i++ {
		def.headings[i] = dec.Uint32()
	}

	if def.flags&0x10 == 0x10 {
		def.renderMethod = dec.Uint32()
	}
	if def.flags&0x20 == 0x20 {
		def.renderInfo.uvInfo.origin.X = dec.Float32()
		def.renderInfo.uvInfo.origin.Y = dec.Float32()
		def.renderInfo.uvInfo.origin.Z = dec.Float32()
		def.renderInfo.uvInfo.uAxis.X = dec.Float32()
		def.renderInfo.uvInfo.uAxis.Y = dec.Float32()
		def.renderInfo.uvInfo.uAxis.Z = dec.Float32()
		def.renderInfo.uvInfo.vAxis.X = dec.Float32()
		def.renderInfo.uvInfo.vAxis.Y = dec.Float32()
		def.renderInfo.uvInfo.vAxis.Z = dec.Float32()
	}

	log.Debugf("%+v\n", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *twoDSpriteDef) build(e *WLD) error {
	return nil
}

func (e *WLD) twoDSpriteDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
