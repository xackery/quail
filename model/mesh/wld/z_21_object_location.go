package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
)

type objectLocation struct {
	nameRef        int32
	actorDefRef    int32
	flags          uint32
	sphereRef      uint32
	currentAction  uint32
	offset         common.Vector3
	rotation       common.Vector3
	unk1           uint32
	boundingRadius float32
	scale          float32
	soundNameRef   int32
	unk2           int32
}

func (e *WLD) objectLocationRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &objectLocation{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.actorDefRef = dec.Int32()
	def.flags = dec.Uint32()
	def.sphereRef = dec.Uint32()
	if def.flags&0x1 == 0x1 {
		def.currentAction = dec.Uint32()
	}
	if def.flags&0x2 == 0x2 {
		def.offset.X = dec.Float32()
		def.offset.Y = dec.Float32()
		def.offset.Z = dec.Float32()
		def.rotation.X = dec.Float32()
		def.rotation.Y = dec.Float32()
		def.rotation.Z = dec.Float32()
		def.unk1 = dec.Uint32()
	}
	if def.flags&0x4 == 0x4 {
		def.boundingRadius = dec.Float32()
	}
	if def.flags&0x8 == 0x8 {
		def.scale = dec.Float32()
	}
	if def.flags&0x10 == 0x10 {
		def.soundNameRef = dec.Int32()
	}
	def.unk2 = dec.Int32()

	if dec.Error() != nil {
		return fmt.Errorf("objectLocationRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *objectLocation) build(e *WLD) error {
	return nil
}

func (e *WLD) objectLocationWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
