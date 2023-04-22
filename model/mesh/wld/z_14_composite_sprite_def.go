package wld

import (
	"fmt"
	"io"
)

type compositeSpriteDef struct {
}

func (e *WLD) compositeSpriteDefRead(r io.ReadSeeker, fragmentOffset int) error {
	return fmt.Errorf("compositeSpriteDefRead not implemented")
}

func (v *compositeSpriteDef) build(e *WLD) error {
	return nil
}
