package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// MaterialList 0x31 49
type MaterialList struct {
	NameRef       int32
	Flags         uint32
	MaterialCount uint32
	MaterialRefs  []uint32
}

func (e *WLD) materialListRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &MaterialList{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.Flags = dec.Uint32()
	def.MaterialCount = dec.Uint32()
	def.MaterialRefs = make([]uint32, def.MaterialCount)
	for i := uint32(0); i < def.MaterialCount; i++ {
		def.MaterialRefs[i] = dec.Uint32()
	}

	if dec.Error() != nil {
		return fmt.Errorf("materialListRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *MaterialList) build(e *WLD) error {
	return nil
}

func (e *WLD) materialListWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
