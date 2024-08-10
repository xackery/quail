package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDirectionalLight is DirectionalLight in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type WldFragDirectionalLight struct {
}

func (e *WldFragDirectionalLight) FragCode() int {
	return FragCodeDirectionalLight
}

func (e *WldFragDirectionalLight) Write(w io.Writer, isNewWorld bool) error {
	return fmt.Errorf("not implemented")
}

func (e *WldFragDirectionalLight) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
