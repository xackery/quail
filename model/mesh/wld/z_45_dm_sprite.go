package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type dmSprite struct {
}

func (e *WLD) dmSpriteRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &dmSprite{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("dmSpriteRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *dmSprite) build(e *WLD) error {
	return nil
}

func (e *WLD) dmSpriteWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
