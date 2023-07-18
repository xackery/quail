package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// 0x1A sphereList
type sphereList struct {
	nameRef          int32
	sphereListDefRef int32
	params1          uint32
}

func (e *WLD) sphereListRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &sphereList{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.sphereListDefRef = dec.Int32()
	def.params1 = dec.Uint32()

	if dec.Error() != nil {
		return fmt.Errorf("sphereListRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *sphereList) build(e *WLD) error {
	return nil
}

func (e *WLD) sphereListWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
