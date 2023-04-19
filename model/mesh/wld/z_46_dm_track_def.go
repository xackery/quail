package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type dmTrackDef struct {
}

func (e *WLD) dmTrackDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &dmTrackDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("dmTrackDefRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *dmTrackDef) build(e *WLD) error {
	return nil
}
