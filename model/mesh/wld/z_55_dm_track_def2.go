package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// dmTrackDef2 0x37 55
type dmTrackDef2 struct {
}

func (e *WLD) dmTrackDef2Read(r io.ReadSeeker, fragmentOffset int) error {
	def := &dmTrackDef2{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("dmTrackDef2Read: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *dmTrackDef2) build(e *WLD) error {
	return nil
}
