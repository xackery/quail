package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type polyhedron struct {
	nameRef     int32
	fragmentRef int32
	flags       uint32
	scale       float32
}

func (e *WLD) polyhedronRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &polyhedron{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.fragmentRef = dec.Int32()
	def.flags = dec.Uint32()
	def.scale = dec.Float32()
	if dec.Error() != nil {
		return fmt.Errorf("polyhedronRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *polyhedron) build(e *WLD) error {
	return nil
}
