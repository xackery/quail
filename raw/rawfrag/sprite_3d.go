package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSprite3D is Sprite3D in libeq, Camera Reference in openzone, 3DSPRITE (ref) in wld, CameraReference in lantern
type WldFragSprite3D struct {
	NameRef        int32
	Sprite3DDefRef int32
	Flags          uint32
}

func (e *WldFragSprite3D) FragCode() int {
	return FragCodeSprite3D
}

func (e *WldFragSprite3D) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.Sprite3DDefRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite3D) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Sprite3DDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
