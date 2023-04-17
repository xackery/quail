package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// 0x14 model
type model struct {
	nameRef          int32
	flags            uint32
	callbackNameRef  int32
	actionCount      uint32
	fragmentRefCount uint32
	boundsRef        int32
	currentAction    uint32
	offset           geo.Vector3
	rotation         geo.Vector3
	unk1             uint32
	actions          []action
	fragmentRefs     []uint32
	unk2             uint32
}

type action struct {
	lodCount uint32
	unk1     uint32
	lods     []float32
}

func (e *WLD) modelRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &model{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	def.callbackNameRef = dec.Int32()
	def.actionCount = dec.Uint32()
	def.fragmentRefCount = dec.Uint32()
	def.boundsRef = dec.Int32()
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
	for i := 0; i < int(def.actionCount); i++ {
		var action action
		action.lodCount = dec.Uint32()
		action.unk1 = dec.Uint32()
		for j := 0; j < int(action.lodCount); j++ {
			action.lods = append(action.lods, dec.Float32())
		}
		def.actions = append(def.actions, action)
	}
	for i := 0; i < int(def.fragmentRefCount); i++ {
		def.fragmentRefs = append(def.fragmentRefs, dec.Uint32())
	}
	def.unk2 = dec.Uint32()

	if dec.Error() != nil {
		return fmt.Errorf("modelRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *model) build(e *WLD) error {
	return nil
}
