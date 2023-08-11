package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// 0x1B lightDef
type lightDef struct {
	nameRef         int32
	flags           uint32
	frameCurrentRef uint32
	sleep           uint32
	lightLevels     []float32
	colors          []geo.Vector3
}

func (e *WLD) lightDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &lightDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	frameCount := dec.Uint32()
	if def.flags&0x1 != 0 {
		def.frameCurrentRef = dec.Uint32()
	}
	if def.flags&0x2 != 0 {
		def.sleep = dec.Uint32()
	}
	if def.flags&0x4 != 0 {
		for i := uint32(0); i < frameCount; i++ {
			def.lightLevels = append(def.lightLevels, dec.Float32())
		}
	}
	// 0x08 is skip frames, unused
	if def.flags&0x10 != 0 {
		for i := uint32(0); i < frameCount; i++ {
			var color geo.Vector3
			color.X = dec.Float32()
			color.Y = dec.Float32()
			color.Z = dec.Float32()
			def.colors = append(def.colors, color)
		}
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *lightDef) build(e *WLD) error {
	return nil
}

func (e *WLD) lightDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
