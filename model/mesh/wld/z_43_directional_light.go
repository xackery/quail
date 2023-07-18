package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type directionalLight struct {
}

func (e *WLD) directionalLightRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &directionalLight{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("directionalLightRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *directionalLight) build(e *WLD) error {
	return nil
}

func (e *WLD) directionalLightWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
