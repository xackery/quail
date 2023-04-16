package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ghostiam/binstruct"
	"github.com/xackery/quail/log"
)

type skeletonTrack struct {
	NameRef          int16
	SkeletonTrackRef int16
	Flags            uint32
}

func (e *WLD) skeletonTrackRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &skeletonTrack{}

	dec := binstruct.NewDecoder(r, binary.LittleEndian)
	err := dec.Decode(def)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *skeletonTrack) build(e *WLD) error {
	return nil
}
