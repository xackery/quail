package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragSprite2D is Sprite2D in libeq, Two-Dimensional Object Reference in openzone, 2DSPRITE (ref) in wld, Fragment07 in lantern
type WldFragSprite2D struct {
	parents       []common.TreeLinker
	children      []common.TreeLinker
	fragID        int
	tag           string
	NameRef       int32  `yaml:"name_ref"`
	TwoDSpriteRef uint32 `yaml:"two_d_sprite_ref"`
	Flags         uint32 `yaml:"flags"`
}

func (e *WldFragSprite2D) FragCode() int {
	return FragCodeSprite2D
}

func (e *WldFragSprite2D) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.TwoDSpriteRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite2D) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.TwoDSpriteRef = dec.Uint32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSprite2D) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragSprite2D) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragSprite2D) Tag() string {
	return e.tag
}

func (e *WldFragSprite2D) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragSprite2D) FragID() int {
	return e.fragID
}

func (e *WldFragSprite2D) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragSprite2D) FragType() string {
	return "S2DI"
}

func (e *WldFragSprite2D) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
