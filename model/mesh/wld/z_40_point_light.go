package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type pointLight struct {
}

func (e *WLD) pointLightRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &pointLight{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("pointLightRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *pointLight) build(e *WLD) error {
	return nil
}

func (e *WLD) pointLightWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
