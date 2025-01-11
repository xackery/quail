package rawfrag

import (
	"io"

	"github.com/xackery/quail/common"
)

// DirectionalLigtOld is empty in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type WldFragDirectionalLightOld struct {
	children []common.TreeLinker
	fragID   int
	tag      string
	parents  []common.TreeLinker
}

func (e *WldFragDirectionalLightOld) FragCode() int {
	return FragCodeDirectionalLightOld
}

func (e *WldFragDirectionalLightOld) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragDirectionalLightOld) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}

func (e *WldFragDirectionalLightOld) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragDirectionalLightOld) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragDirectionalLightOld) Tag() string {
	return e.tag
}

func (e *WldFragDirectionalLightOld) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragDirectionalLightOld) FragID() int {
	return e.fragID
}

func (e *WldFragDirectionalLightOld) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragDirectionalLightOld) FragType() string {
	return "DLOD"
}

func (e *WldFragDirectionalLightOld) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
