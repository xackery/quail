package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type skyRegion struct {
}

func (e *WLD) skyRegionRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &skyRegion{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("skyRegionRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *skyRegion) build(e *WLD) error {
	return nil
}

func (e *WLD) skyRegionWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
