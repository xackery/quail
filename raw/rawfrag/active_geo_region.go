package rawfrag

import "io"

// WldFragActiveGeoRegion is empty in libeq, empty in openzone, ACTIVEGEOMETRYREGION in wld
type WldFragActiveGeoRegion struct {
	nameRef int32
}

func (e *WldFragActiveGeoRegion) FragCode() int {
	return FragCodeActiveGeoRegion
}

func (e *WldFragActiveGeoRegion) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragActiveGeoRegion) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}

func (e *WldFragActiveGeoRegion) NameRef() int32 {
	return e.nameRef
}
