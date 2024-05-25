package rawfrag

import (
	"io"
)

// WldFragDmRGBTrack is DmRGBTrack in libeq, Vertex Color Reference in openzone, empty in wld, VertexColorsReference in lantern
type WldFragDmRGBTrack struct {
}

func (e *WldFragDmRGBTrack) FragCode() int {
	return FragCodeDmRGBTrack
}

func (e *WldFragDmRGBTrack) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDmRGBTrack) Read(r io.ReadSeeker) error {
	return nil
}
