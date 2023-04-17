package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// 0x16 sphere
type sphere struct {
	nameRef int32
	radius  float32
}

func (e *WLD) sphereRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &sphere{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()  //nameRef
	def.radius = dec.Float32() //radius

	if dec.Error() != nil {
		return fmt.Errorf("sphereRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *sphere) build(e *WLD) error {
	return nil
}
