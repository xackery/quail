package rawfrag

import (
	"io"

	"github.com/xackery/quail/common"
)

// WldFragUserData is empty in libeq, empty in openzone, USERDATA in wld
type WldFragUserData struct {
	children []common.TreeLinker
	fragID   int
	tag      string
	parents  []common.TreeLinker
}

func (e *WldFragUserData) FragCode() int {
	return FragCodeUserData
}

func (e *WldFragUserData) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragUserData) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}

func (e *WldFragUserData) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragUserData) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragUserData) Tag() string {
	return e.tag
}

func (e *WldFragUserData) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragUserData) FragID() int {
	return e.fragID
}

func (e *WldFragUserData) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragUserData) FragType() string {
	return "USRD"
}

func (e *WldFragUserData) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
