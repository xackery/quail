package rawfrag

import "io"

// WldFragSkyRegion is empty in libeq, empty in openzone, SKYREGION in wld
type WldFragSkyRegion struct {
}

func (e *WldFragSkyRegion) FragCode() int {
	return FragCodeSkyRegion
}

func (e *WldFragSkyRegion) Write(w io.Writer) error {
	return nil
}

func (e *WldFragSkyRegion) Read(r io.ReadSeeker) error {
	return nil
}
