package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDMSprite is DmSprite in libeq, Mesh Reference in openzone, empty in wld, MeshReference in lantern
type WldFragDMSprite struct {
	NameRef     int32  `yaml:"name_ref"`
	DMSpriteRef int32  `yaml:"dm_sprite_ref"`
	Params      uint32 `yaml:"params"`
}

func (e *WldFragDMSprite) FragCode() int {
	return FragCodeDMSprite
}

func (e *WldFragDMSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.DMSpriteRef)
	enc.Uint32(e.Params)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDMSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.DMSpriteRef = dec.Int32()
	e.Params = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
