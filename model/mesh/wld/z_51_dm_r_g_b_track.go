package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type dmRGBTrack struct {
}

func (e *WLD) dmRGBTrackRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &dmRGBTrack{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("dmRGBTrackRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *dmRGBTrack) build(e *WLD) error {
	return nil
}

func (e *WLD) dmRGBTrackWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
