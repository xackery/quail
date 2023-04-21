package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type dmTrack struct {
}

func (e *WLD) dmTrackRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &dmTrack{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("dmTrackRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *dmTrack) build(e *WLD) error {
	return nil
}