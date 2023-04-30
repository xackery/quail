package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type track struct {
	nameRef  int32
	trackRef int32
	flags    uint32
	sleep    uint32 // if 0x01 is set, this is the number of milliseconds to sleep before starting the animation
}

func (e *WLD) trackRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &track{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.trackRef = dec.Int32()
	def.flags = dec.Uint32()
	// TODO: find an example with a flag sleep enabled
	if def.flags&0x01 == 0x01 {
		log.Debugf("Found a sleep enabled track!")
		def.sleep = dec.Uint32()
	}

	if dec.Error() != nil {
		return fmt.Errorf("trackRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *track) build(e *WLD) error {
	return nil
}
