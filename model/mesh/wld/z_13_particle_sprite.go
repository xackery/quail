package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type particleSprite struct {
	nameRef              int32
	particleSpriteDefRef int32
	flags                uint32
}

func (e *WLD) particleSpriteRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &particleSprite{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.particleSpriteDefRef = dec.Int32()
	def.flags = dec.Uint32()
	if dec.Error() != nil {
		return fmt.Errorf("particleSpriteRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *particleSprite) build(e *WLD) error {
	return nil
}

func (e *WLD) particleSpriteWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
