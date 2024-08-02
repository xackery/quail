package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragHierarchicalSprite is HierarchicalSprite in libeq, SkeletonTrackSetReference in openzone, HIERARCHICALSPRITE (ref) in wld, SkeletonHierarchyReference in lantern
type WldFragHierarchicalSprite struct {
	NameRef               uint32
	HierarchicalSpriteRef uint32
	Param                 uint32
}

func (e *WldFragHierarchicalSprite) FragCode() int {
	return FragCodeHierarchicalSprite
}

func (e *WldFragHierarchicalSprite) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Uint32(e.NameRef)
	enc.Uint32(e.HierarchicalSpriteRef)
	enc.Uint32(e.Param)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragHierarchicalSprite) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Uint32()
	e.HierarchicalSpriteRef = dec.Uint32()
	e.Param = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
