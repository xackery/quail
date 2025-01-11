package rawfrag

import "io"

// WldFragUserData is empty in libeq, empty in openzone, USERDATA in wld
type WldFragUserData struct {
}

func (e *WldFragUserData) FragCode() int {
	return FragCodeUserData
}

func (e *WldFragUserData) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragUserData) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}

func (e *WldFragUserData) NameRef() int32 {
	return 0
}
