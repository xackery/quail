package rawfrag

import "io"

// DirectionalLigtOld is empty in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type WldFragDirectionalLightOld struct {
}

func (e *WldFragDirectionalLightOld) FragCode() int {
	return FragCodeDirectionalLightOld
}

func (e *WldFragDirectionalLightOld) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDirectionalLightOld) Read(r io.ReadSeeker) error {
	return nil
}
