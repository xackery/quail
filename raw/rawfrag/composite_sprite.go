package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragCompositeSprite is empty in libeq, empty in openzone, COMPOSITESPRITE (ref) in wld
type WldFragCompositeSprite struct {
	parents               []common.TreeLinker
	children              []common.TreeLinker
	fragID                int
	tag                   string
	NameRef               int32  `yaml:"name_ref"`
	CompositeSpriteDefRef int32  `yaml:"composite_sprite_def_ref"`
	Flags                 uint32 `yaml:"flags"`
}

func (e *WldFragCompositeSprite) FragCode() int {
	return FragCodeCompositeSprite
}

func (e *WldFragCompositeSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.CompositeSpriteDefRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragCompositeSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.CompositeSpriteDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragCompositeSprite) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragCompositeSprite) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragCompositeSprite) Tag() string {
	return e.tag
}

func (e *WldFragCompositeSprite) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragCompositeSprite) FragID() int {
	return e.fragID
}

func (e *WldFragCompositeSprite) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragCompositeSprite) FragType() string {
	return "COSI"
}

func (e *WldFragCompositeSprite) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
