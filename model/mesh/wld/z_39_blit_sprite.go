package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type blitSprite struct {
}

func (e *WLD) blitSpriteRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &blitSprite{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("blitSpriteRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *blitSprite) build(e *WLD) error {
	return nil
}

func (e *WLD) blitSpriteWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
