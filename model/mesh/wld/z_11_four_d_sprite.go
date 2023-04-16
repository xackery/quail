package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type fourDSprite struct {
	NameRef  int32
	FourDRef int32
	Params1  uint32
}

func (e *WLD) fourDSpriteRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &fourDSprite{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.FourDRef = dec.Int32()
	def.Params1 = dec.Uint32()
	if dec.Error() != nil {
		return fmt.Errorf("fourDSpriteRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *fourDSprite) build(e *WLD) error {
	return nil
}
