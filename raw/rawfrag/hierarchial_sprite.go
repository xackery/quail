package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragHierarchialSprite is HierarchialSprite in libeq, SkeletonTrackSetReference in openzone, HIERARCHIALSPRITE (ref) in wld, SkeletonHierarchyReference in lantern
type WldFragHierarchialSprite struct {
	NameRef              int16  `yaml:"name_ref"`
	HierarchialSpriteRef int16  `yaml:"hierarchial_sprite_ref"`
	Flags                uint32 `yaml:"flags"`
}

func (e *WldFragHierarchialSprite) FragCode() int {
	return FragCodeHierarchialSprite
}

func (e *WldFragHierarchialSprite) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int16(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Int16(e.HierarchialSpriteRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragHierarchialSprite) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int16()
	e.Flags = dec.Uint32()
	e.HierarchialSpriteRef = dec.Int16()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
