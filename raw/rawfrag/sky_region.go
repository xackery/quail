package rawfrag

import (
	"io"

	"github.com/xackery/quail/common"
)

// WldFragSkyRegion is empty in libeq, empty in openzone, SKYREGION in wld
type WldFragSkyRegion struct {
	children []common.TreeLinker
	fragID   int
	tag      string
	parents  []common.TreeLinker
}

func (e *WldFragSkyRegion) FragCode() int {
	return FragCodeSkyRegion
}

func (e *WldFragSkyRegion) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragSkyRegion) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}

func (e *WldFragSkyRegion) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragSkyRegion) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragSkyRegion) Tag() string {
	return e.tag
}

func (e *WldFragSkyRegion) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragSkyRegion) FragID() int {
	return e.fragID
}

func (e *WldFragSkyRegion) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragSkyRegion) FragType() string {
	return "SKYD"
}

func (e *WldFragSkyRegion) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
