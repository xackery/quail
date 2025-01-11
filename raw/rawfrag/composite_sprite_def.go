package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragCompositeSpriteDef is empty in libeq, empty in openzone, COMPOSITESPRITEDEF in wld, Actor in lantern
type WldFragCompositeSpriteDef struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
	NameRef  int32  `yaml:"name_ref"`
	Flags    uint32 `yaml:"flags"`
}

func (e *WldFragCompositeSpriteDef) FragCode() int {
	return FragCodeCompositeSpriteDef
}

func (e *WldFragCompositeSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragCompositeSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragCompositeSpriteDef) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragCompositeSpriteDef) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragCompositeSpriteDef) Tag() string {
	return e.tag
}

func (e *WldFragCompositeSpriteDef) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragCompositeSpriteDef) FragID() int {
	return e.fragID
}

func (e *WldFragCompositeSpriteDef) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragCompositeSpriteDef) FragType() string {
	return "COSD"
}

func (e *WldFragCompositeSpriteDef) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
