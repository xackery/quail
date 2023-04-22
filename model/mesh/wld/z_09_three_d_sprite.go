package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// 0x09
type threeDSprite struct {
	NameRef   int32
	ThreeDRef int32
	Flags     uint32
}

func (e *WLD) threeDSpriteRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &threeDSprite{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.ThreeDRef = dec.Int32()
	def.Flags = dec.Uint32()
	if dec.Error() != nil {
		return fmt.Errorf("threeDSpriteRead: %s", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *threeDSprite) build(e *WLD) error {
	return nil
}
