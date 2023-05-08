package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// 0x1F soundDef
type soundDef struct {
	nameRef int32
	flags   uint32
}

func (e *WLD) soundDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &soundDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	if dec.Error() != nil {
		return fmt.Errorf("soundDefRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *soundDef) build(e *WLD) error {
	return nil
}

func (e *WLD) soundDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
