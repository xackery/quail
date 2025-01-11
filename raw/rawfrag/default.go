package rawfrag

import "io"

// WldFragDefault is empty in libeq, empty in openzone, DEFAULT?? in wld
type WldFragDefault struct {
}

func (e *WldFragDefault) FragCode() int {
	return FragCodeDefault
}

func (e *WldFragDefault) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragDefault) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}

func (e *WldFragDefault) NameRef() int32 {
	return 0
}
