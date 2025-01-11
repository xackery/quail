package rawfrag

import (
	"io"

	"github.com/xackery/quail/common"
)

// WldFragDMTrack is DmTrackDef in libeq, empty in openzone, empty in wld
type WldFragDMTrackDef struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
}

func (e *WldFragDMTrackDef) FragCode() int {
	return FragCodeDMTrackDef
}

func (e *WldFragDMTrackDef) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragDMTrackDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}

func (e *WldFragDMTrackDef) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragDMTrackDef) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragDMTrackDef) Tag() string {
	return e.tag
}

func (e *WldFragDMTrackDef) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragDMTrackDef) FragID() int {
	return e.fragID
}

func (e *WldFragDMTrackDef) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragDMTrackDef) FragType() string {
	return "DMTD"
}

func (e *WldFragDMTrackDef) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
