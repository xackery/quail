package rawfrag

import "io"

// WldFragActiveGeoRegion is empty in libeq, empty in openzone, ACTIVEGEOMETRYREGION in wld
type WldFragActiveGeoRegion struct {
}

func (e *WldFragActiveGeoRegion) FragCode() int {
	return FragCodeActiveGeoRegion
}

func (e *WldFragActiveGeoRegion) Write(w io.Writer) error {
	return nil
}

func (e *WldFragActiveGeoRegion) Read(r io.ReadSeeker) error {
	return nil
}
