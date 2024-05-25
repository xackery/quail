package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragBlitSprite is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSpriteReference in lantern
type WldFragBlitSprite struct {
}

func (e *WldFragBlitSprite) FragCode() int {
	return FragCodeBlitSprite
}

func (e *WldFragBlitSprite) Write(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func (e *WldFragBlitSprite) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
