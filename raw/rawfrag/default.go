package rawfrag

import (
	"io"

	"github.com/xackery/quail/common"
)

// WldFragDefault is empty in libeq, empty in openzone, DEFAULT?? in wld
type WldFragDefault struct {
	children []common.TreeLinker
	fragID   int
	tag      string
	parents  []common.TreeLinker
}

func (e *WldFragDefault) FragCode() int {
	return FragCodeDefault
}

func (e *WldFragDefault) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragDefault) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}

func (e *WldFragDefault) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragDefault) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragDefault) Tag() string {
	return e.tag
}

func (e *WldFragDefault) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragDefault) FragID() int {
	return e.fragID
}

func (e *WldFragDefault) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragDefault) FragType() string {
	return "DEFD"
}

func (e *WldFragDefault) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
