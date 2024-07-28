package rawfrag

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDmRGBTrack is DmRGBTrack in libeq, Vertex Color Reference in openzone, empty in wld, VertexColorsReference in lantern
type WldFragDmRGBTrack struct {
	NameRef  int32
	TrackRef int32
	Flags    uint32
}

func (e *WldFragDmRGBTrack) FragCode() int {
	return FragCodeDmRGBTrack
}

func (e *WldFragDmRGBTrack) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.Int32(e.NameRef)
	enc.Int32(e.TrackRef)
	enc.Uint32(e.Flags)

	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func (e *WldFragDmRGBTrack) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	e.NameRef = dec.Int32()
	e.TrackRef = dec.Int32()
	e.Flags = dec.Uint32()

	if dec.Error() != nil {
		return dec.Error()
	}
	return nil
}
