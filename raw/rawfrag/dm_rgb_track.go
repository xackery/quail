package rawfrag

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragDmRGBTrack is DmRGBTrack in libeq, Vertex Color Reference in openzone, empty in wld, VertexColorsReference in lantern
type WldFragDmRGBTrack struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
	NameRef  int32
	TrackRef int32
	Flags    uint32
}

func (e *WldFragDmRGBTrack) FragCode() int {
	return FragCodeDmRGBTrack
}

func (e *WldFragDmRGBTrack) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.Int32(e.NameRef)
	enc.Int32(e.TrackRef)
	enc.Uint32(e.Flags)

	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func (e *WldFragDmRGBTrack) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	e.NameRef = dec.Int32()
	e.TrackRef = dec.Int32()
	e.Flags = dec.Uint32()

	if dec.Error() != nil {
		return dec.Error()
	}
	return nil
}

func (e *WldFragDmRGBTrack) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragDmRGBTrack) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragDmRGBTrack) Tag() string {
	return e.tag
}

func (e *WldFragDmRGBTrack) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragDmRGBTrack) FragID() int {
	return e.fragID
}

func (e *WldFragDmRGBTrack) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragDmRGBTrack) FragType() string {
	return "RGTI"
}

func (e *WldFragDmRGBTrack) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
