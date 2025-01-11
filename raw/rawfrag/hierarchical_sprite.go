package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragHierarchicalSprite is HierarchicalSprite in libeq, SkeletonTrackSetReference in openzone, HIERARCHICALSPRITE (ref) in wld, SkeletonHierarchyReference in lantern
type WldFragHierarchicalSprite struct {
	parents               []common.TreeLinker
	children              []common.TreeLinker
	fragID                int
	tag                   string
	NameRef               uint32
	HierarchicalSpriteRef uint32
	Param                 uint32
}

func (e *WldFragHierarchicalSprite) FragCode() int {
	return FragCodeHierarchicalSprite
}

func (e *WldFragHierarchicalSprite) Write(w io.Writer, isNewWorld bool) error {
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

func (e *WldFragHierarchicalSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
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

func (e *WldFragHierarchicalSprite) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragHierarchicalSprite) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragHierarchicalSprite) Tag() string {
	return e.tag
}

func (e *WldFragHierarchicalSprite) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragHierarchicalSprite) FragID() int {
	return e.fragID
}

func (e *WldFragHierarchicalSprite) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragHierarchicalSprite) FragType() string {
	return "HISI"
}

func (e *WldFragHierarchicalSprite) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
