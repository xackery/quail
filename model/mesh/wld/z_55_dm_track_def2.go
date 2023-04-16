package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ghostiam/binstruct"
	"github.com/xackery/quail/log"
)

// dmTrackDef2 0x37 55
type dmTrackDef2 struct {
}

func (e *WLD) dmTrackDef2Read(r io.ReadSeeker, fragmentOffset int) error {
	def := &dmTrackDef2{}

	dec := binstruct.NewDecoder(r, binary.LittleEndian)
	err := dec.Decode(def)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *dmTrackDef2) build(e *WLD) error {
	return nil
}
