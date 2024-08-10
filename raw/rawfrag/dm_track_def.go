package rawfrag

import "io"

// WldFragDMTrack is DmTrackDef in libeq, empty in openzone, empty in wld
type WldFragDMTrackDef struct {
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

// WldFragDMTrack is DmTrack in libeq, Mesh Animated Vertices Reference in openzone, empty in wld, MeshAnimatedVerticesReference in lantern
type WldFragDMTrack struct {
}

func (e *WldFragDMTrack) FragCode() int {
	return FragCodeDMTrack
}

func (e *WldFragDMTrack) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragDMTrack) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}
