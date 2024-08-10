package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSimpleSprite is SimpleSprite in libeq, Texture Bitmap Info Reference in openzone, SIMPLESPRITEINST in wld, BitmapInfoReference in lantern
type WldFragSimpleSprite struct {
	NameRef   int32  `yaml:"name_ref"`
	SpriteRef uint32 `yaml:"sprite_ref"`
	Flags     uint32 `yaml:"flags"`
}

func (e *WldFragSimpleSprite) FragCode() int {
	return FragCodeSimpleSprite
}

func (e *WldFragSimpleSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.SpriteRef)
	enc.Uint32(e.Flags)
	enc.Bytes(make([]byte, 2)) // TODO: why 2 extra bytes?
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSimpleSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.SpriteRef = dec.Uint32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
