package rawfrag

import "io"

// WldFragSkyRegion is empty in libeq, empty in openzone, SKYREGION in wld
type WldFragSkyRegion struct {
}

func (e *WldFragSkyRegion) FragCode() int {
	return FragCodeSkyRegion
}

func (e *WldFragSkyRegion) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragSkyRegion) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}
