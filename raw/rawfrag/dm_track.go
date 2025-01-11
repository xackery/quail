package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDMTrack is DmTrack in libeq, Mesh Animated Vertices Reference in openzone, empty in wld, MeshAnimatedVerticesReference in lantern
type WldFragDMTrack struct {
	nameRef  int32
	TrackRef int32
	Flags    uint32
}

func (e *WldFragDMTrack) FragCode() int {
	return FragCodeDMTrack
}

func (e *WldFragDMTrack) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
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
	e.nameRef = dec.Int32()
	e.TrackRef = dec.Int32()
	e.Flags = dec.Uint32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDMTrack) NameRef() int32 {
	return e.nameRef
}
