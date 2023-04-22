package wld

import (
	"fmt"
	"io"
)

type compositeSprite struct {
}

func (e *WLD) compositeSpriteRead(r io.ReadSeeker, fragmentOffset int) error {
	return fmt.Errorf("compositeSpriteRead not implemented")
}

func (v *compositeSprite) build(e *WLD) error {
	return nil
}
