package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// 0x1C light
type light struct {
	nameRef     int32
	lightDefRef int32
	flags       uint32
}

func (e *WLD) lightRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &light{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.lightDefRef = dec.Int32()
	def.flags = dec.Uint32()

	if dec.Error() != nil {
		return fmt.Errorf("lightRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *light) build(e *WLD) error {
	return nil
}

func (e *WLD) lightWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
