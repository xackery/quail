package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type compositeSpriteDef struct {
}

func (e *WLD) compositeSpriteDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &compositeSpriteDef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return fmt.Errorf("decode: %s", dec.Error())
	}

	log.Debugf("compositeSpriteDef: %+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *compositeSpriteDef) build(e *WLD) error {
	return nil
}

func (e *WLD) compositeSpriteDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
