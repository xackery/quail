package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type dmSpriteDef struct {
}

func (e *WLD) dmSpriteDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &dmSpriteDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("dmSpriteDefRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *dmSpriteDef) build(e *WLD) error {
	return nil
}
