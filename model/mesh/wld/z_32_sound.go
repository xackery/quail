package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// 0x20 sound
type sound struct {
	nameRef int32
	flags   uint32
}

func (e *WLD) soundRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &sound{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	if dec.Error() != nil {
		return fmt.Errorf("soundRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *sound) build(e *WLD) error {
	return nil
}

func (e *WLD) soundWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
