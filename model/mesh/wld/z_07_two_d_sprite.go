package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type twoDSprite struct {
	NameRef       int32
	TwoDSpriteRef uint32
	Flags         uint32
}

func (e *WLD) twoDSpriteRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &twoDSprite{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.TwoDSpriteRef = dec.Uint32()
	def.Flags = dec.Uint32()
	if dec.Error() != nil {
		return fmt.Errorf("twoDSpriteRead: %s", dec.Error())
	}

	log.Debugf("twoDSprite: %+v\n", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *twoDSprite) build(e *WLD) error {
	return nil
}

func (e *WLD) twoDSpriteWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
