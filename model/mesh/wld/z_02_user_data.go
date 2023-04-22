package wld

import (
	"fmt"
	"io"
)

type userData struct {
}

func (e *WLD) userDataRead(r io.ReadSeeker, fragmentOffset int) error {
	return fmt.Errorf("userDataRead: not implemented")
}

func (v *userData) build(e *WLD) error {
	return nil
}
