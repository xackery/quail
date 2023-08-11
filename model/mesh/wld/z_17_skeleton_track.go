package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// 0x11 skeletonTrack
type skeletonTrack struct {
	nameRef          int16
	skeletonTrackRef int16
	flags            uint32
}

func (e *WLD) skeletonTrackRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &skeletonTrack{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int16()
	def.flags = dec.Uint32()
	def.skeletonTrackRef = dec.Int16()

	if dec.Error() != nil {
		return fmt.Errorf("skeletonTrackRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *skeletonTrack) build(e *WLD) error {
	return nil
}

func (e *WLD) skeletonTrackWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
