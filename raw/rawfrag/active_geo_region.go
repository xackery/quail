package rawfrag

import (
	"io"

	"github.com/xackery/quail/common"
)

// WldFragActiveGeoRegion is empty in libeq, empty in openzone, ACTIVEGEOMETRYREGION in wld
type WldFragActiveGeoRegion struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
}

func (e *WldFragActiveGeoRegion) FragCode() int {
	return FragCodeActiveGeoRegion
}

func (e *WldFragActiveGeoRegion) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragActiveGeoRegion) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}

func (e *WldFragActiveGeoRegion) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragActiveGeoRegion) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragActiveGeoRegion) Tag() string {
	return e.tag
}

func (e *WldFragActiveGeoRegion) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragActiveGeoRegion) FragID() int {
	return e.fragID
}

func (e *WldFragActiveGeoRegion) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragActiveGeoRegion) FragType() string {
	return "AGR"
}

func (e *WldFragActiveGeoRegion) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
