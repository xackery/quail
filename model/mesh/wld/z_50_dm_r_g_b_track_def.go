package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type dmRGBTrackDef struct {
}

func (e *WLD) dmRGBTrackDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &dmRGBTrackDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("dmRGBTrackDefRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *dmRGBTrackDef) build(e *WLD) error {
	return nil
}

func (e *WLD) dmRGBTrackDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
