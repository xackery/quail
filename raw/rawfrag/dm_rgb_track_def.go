package rawfrag

import "io"

// WldFragDmRGBTrackDef is a list of colors, one per vertex, for baked lighting. It is DmRGBTrackDef in libeq, Vertex Color in openzone, empty in wld, VertexColors in lantern
type WldFragDmRGBTrackDef struct {
}

func (e *WldFragDmRGBTrackDef) FragCode() int {
	return FragCodeDmRGBTrackDef
}

func (e *WldFragDmRGBTrackDef) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDmRGBTrackDef) Read(r io.ReadSeeker) error {
	return nil
}
