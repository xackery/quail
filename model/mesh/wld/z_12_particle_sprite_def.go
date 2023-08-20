package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
)

// 0x0c particleSpriteDef (particle sprite definition)
type particleSpriteDef struct {
	nameRef       int32
	flags         uint32
	verticesCount uint32
	unknown       uint32
	centerOffset  common.Vector3
	radius        float32
	vertices      []common.Vector3
	renderMethod  uint32
	renderInfo    renderInfo
	pen           []uint32
}

func (e *WLD) particleSpriteDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &particleSpriteDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	def.verticesCount = dec.Uint32()
	def.unknown = dec.Uint32()
	if def.flags&0x01 != 0 { // has center offset
		def.centerOffset.X = dec.Float32()
		def.centerOffset.Y = dec.Float32()
		def.centerOffset.Z = dec.Float32()
	}
	if def.flags&0x02 != 0 { // has radius
		def.radius = dec.Float32()
	}
	if def.verticesCount > 0 { // has vertices
		def.vertices = make([]common.Vector3, def.verticesCount)
		for i := uint32(0); i < def.verticesCount; i++ {
			def.vertices[i].X = dec.Float32()
			def.vertices[i].Y = dec.Float32()
			def.vertices[i].Z = dec.Float32()
		}
	}
	def.renderMethod = dec.Uint32()
	def.renderInfo.uvInfo.origin.X = dec.Float32()
	def.renderInfo.uvInfo.origin.Y = dec.Float32()
	def.renderInfo.uvInfo.origin.Z = dec.Float32()
	def.renderInfo.uvInfo.uAxis.X = dec.Float32()
	def.renderInfo.uvInfo.uAxis.Y = dec.Float32()
	def.renderInfo.uvInfo.uAxis.Z = dec.Float32()
	def.renderInfo.uvInfo.vAxis.X = dec.Float32()
	def.renderInfo.uvInfo.vAxis.Y = dec.Float32()
	def.renderInfo.uvInfo.vAxis.Z = dec.Float32()
	for i := 0; i < int(def.verticesCount); i++ {
		def.pen = append(def.pen, dec.Uint32())
	}

	if dec.Error() != nil {
		return fmt.Errorf("particleSpriteDefRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *particleSpriteDef) build(e *WLD) error {
	return nil
}

func (e *WLD) particleSpriteDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
