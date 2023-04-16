package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ghostiam/binstruct"
	"github.com/xackery/quail/log"
)

type activeGeoRegion struct {
}

func (e *WLD) activeGeoRegionRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &activeGeoRegion{}

	dec := binstruct.NewDecoder(r, binary.LittleEndian)
	err := dec.Decode(def)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *activeGeoRegion) build(e *WLD) error {
	return nil
}
