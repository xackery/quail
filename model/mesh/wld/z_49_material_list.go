package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// materialList 0x31 49
type materialList struct {
	nameRef       int32
	flags         uint32
	materialCount uint32
	materialRefs  []uint32
}

func (e *WLD) materialListRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &materialList{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	def.materialCount = dec.Uint32()
	def.materialRefs = make([]uint32, def.materialCount)
	for i := uint32(0); i < def.materialCount; i++ {
		def.materialRefs[i] = dec.Uint32()
	}

	if dec.Error() != nil {
		return fmt.Errorf("materialListRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *materialList) build(e *WLD) error {
	return nil
}