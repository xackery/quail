package rawfrag

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDmRGBTrack is DmRGBTrack in libeq, Vertex Color Reference in openzone, empty in wld, VertexColorsReference in lantern
type WldFragDmRGBTrack struct {
	nameRef  int32
	TrackRef int32
	Flags    uint32
}

func (e *WldFragDmRGBTrack) FragCode() int {
	return FragCodeDmRGBTrack
}

func (e *WldFragDmRGBTrack) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.Int32(e.nameRef)
	enc.Int32(e.TrackRef)
	enc.Uint32(e.Flags)

	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func (e *WldFragDmRGBTrack) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	e.nameRef = dec.Int32()
	e.TrackRef = dec.Int32()
	e.Flags = dec.Uint32()

	if dec.Error() != nil {
		return dec.Error()
	}
	return nil
}

func (e *WldFragDmRGBTrack) NameRef() int32 {
	return e.nameRef
}
