package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragDMSprite is DmSprite in libeq, Mesh Reference in openzone, empty in wld, MeshReference in lantern
type WldFragDMSprite struct {
	parents     []common.TreeLinker
	children    []common.TreeLinker
	fragID      int
	tag         string
	NameRef     int32  `yaml:"name_ref"`
	DMSpriteRef int32  `yaml:"dm_sprite_ref"`
	Params      uint32 `yaml:"params"`
}

func (e *WldFragDMSprite) FragCode() int {
	return FragCodeDMSprite
}

func (e *WldFragDMSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.DMSpriteRef)
	enc.Uint32(e.Params)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDMSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.DMSpriteRef = dec.Int32()
	e.Params = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragDMSprite) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragDMSprite) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragDMSprite) Tag() string {
	return e.tag
}

func (e *WldFragDMSprite) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragDMSprite) FragID() int {
	return e.fragID
}

func (e *WldFragDMSprite) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragDMSprite) FragType() string {
	return "DMSI"
}

func (e *WldFragDMSprite) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
