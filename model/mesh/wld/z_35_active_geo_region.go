package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type activeGeoRegion struct {
}

func (e *WLD) activeGeoRegionRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &activeGeoRegion{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("activeGeoRegionRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *activeGeoRegion) build(e *WLD) error {
	return nil
}

func (e *WLD) activeGeoRegionWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
