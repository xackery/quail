package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

type fourDSpriteDef struct {
	NameRef         int32
	Flags           uint32
	FrameCount      uint32
	PolyRef         int32
	centerOffset    geo.Vector3 `bin:"-"`
	radius          float32     `bin:"-"`
	currentFrame    uint32      `bin:"-"`
	sleep           uint32      `bin:"-"`
	spriteFragments []uint32
}

func (e *WLD) fourDSpriteDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &fourDSpriteDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.Flags = dec.Uint32()
	def.FrameCount = dec.Uint32()
	def.PolyRef = dec.Int32()
	if def.Flags&0x01 != 0 { // has center offset
		def.centerOffset.X = dec.Float32()
		def.centerOffset.Y = dec.Float32()
		def.centerOffset.Z = dec.Float32()
	}
	if def.Flags&0x02 != 0 { // has radius
		def.radius = dec.Float32()
	}
	if def.Flags&0x04 != 0 { // has current frame
		def.currentFrame = dec.Uint32()
	}
	if def.Flags&0x08 != 0 { // has sleep
		def.sleep = dec.Uint32()
	}
	if def.Flags&0x10 != 0 { // has sprite fragments
		def.spriteFragments = make([]uint32, def.FrameCount)
		for i := uint32(0); i < def.FrameCount; i++ {
			def.spriteFragments[i] = dec.Uint32()
		}
	}

	if dec.Error() != nil {
		return fmt.Errorf("fourDSpriteDefRead: %s", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *fourDSpriteDef) build(e *WLD) error {
	return nil
}
