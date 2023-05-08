package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type compositeSprite struct {
}

func (e *WLD) compositeSpriteRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &compositeSprite{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return fmt.Errorf("compositeSpriteRead: %s", dec.Error())
	}

	log.Debugf("compositeSprite: %+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *compositeSprite) build(e *WLD) error {
	return nil
}

func (e *WLD) compositeSpriteWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
