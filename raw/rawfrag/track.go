package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragTrack is a bone in a skeleton. It is Track in libeq, Mob Skeleton Piece Track Reference in openzone, TRACKINSTANCE in wld, TrackDefFragment in lantern
type WldFragTrack struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
	NameRef  int32
	TrackRef int32
	Flags    uint32
	Sleep    uint32 // if 0x01 is set, this is the number of milliseconds to sleep before starting the animation
}

func (e *WldFragTrack) FragCode() int {
	return FragCodeTrack
}

func (e *WldFragTrack) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.TrackRef)
	enc.Uint32(e.Flags)
	if e.Flags&0x01 == 0x01 {
		enc.Uint32(e.Sleep)
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragTrack) Read(r io.ReadSeeker, isNewWorld bool) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.TrackRef = dec.Int32()
	e.Flags = dec.Uint32()
	if e.Flags&0x01 == 0x01 {
		e.Sleep = dec.Uint32()
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragTrack) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragTrack) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragTrack) Tag() string {
	return e.tag
}

func (e *WldFragTrack) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragTrack) FragID() int {
	return e.fragID
}

func (e *WldFragTrack) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragTrack) FragType() string {
	return "TRKI"
}

func (e *WldFragTrack) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
