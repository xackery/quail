package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragBlitSprite is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSpriteReference in lantern
type WldFragBlitSprite struct {
	parents       []common.TreeLinker
	children      []common.TreeLinker
	fragID        int
	tag           string
	NameRef       int32
	BlitSpriteRef int32
	Unknown       int32
}

func (e *WldFragBlitSprite) FragCode() int {
	return FragCodeBlitSprite
}

func (e *WldFragBlitSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.BlitSpriteRef)
	enc.Int32(e.Unknown)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragBlitSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.BlitSpriteRef = dec.Int32()
	e.Unknown = dec.Int32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragBlitSprite) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragBlitSprite) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragBlitSprite) Tag() string {
	return e.tag
}

func (e *WldFragBlitSprite) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragBlitSprite) FragID() int {
	return e.fragID
}

func (e *WldFragBlitSprite) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragBlitSprite) FragType() string {
	return "BLSI"
}

func (e *WldFragBlitSprite) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
