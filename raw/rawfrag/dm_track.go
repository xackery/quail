package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragDMTrack is DmTrack in libeq, Mesh Animated Vertices Reference in openzone, empty in wld, MeshAnimatedVerticesReference in lantern
type WldFragDMTrack struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
	NameRef  int32
	TrackRef int32
	Flags    uint32
}

func (e *WldFragDMTrack) FragCode() int {
	return FragCodeDMTrack
}

func (e *WldFragDMTrack) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.TrackRef)
	enc.Uint32(e.Flags)

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDMTrack) Read(r io.ReadSeeker, isNewWorld bool) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.TrackRef = dec.Int32()
	e.Flags = dec.Uint32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDMTrack) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragDMTrack) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragDMTrack) Tag() string {
	return e.tag
}

func (e *WldFragDMTrack) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragDMTrack) FragID() int {
	return e.fragID
}

func (e *WldFragDMTrack) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragDMTrack) FragType() string {
	return "DMTI"
}

func (e *WldFragDMTrack) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
