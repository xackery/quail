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
