package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ghostiam/binstruct"
	"github.com/xackery/quail/log"
)

type compositeSprite struct {
}

func (e *WLD) compositeSpriteRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &compositeSprite{}

	dec := binstruct.NewDecoder(r, binary.LittleEndian)
	err := dec.Decode(def)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *compositeSprite) build(e *WLD) error {
	return nil
}
