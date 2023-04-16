package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ghostiam/binstruct"
	"github.com/xackery/quail/log"
)

type directionalLightOld struct {
}

func (e *WLD) directionalLightOldRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &directionalLightOld{}

	dec := binstruct.NewDecoder(r, binary.LittleEndian)
	err := dec.Decode(def)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *directionalLightOld) build(e *WLD) error {
	return nil
}
