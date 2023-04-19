package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type ambientLight struct {
}

func (e *WLD) ambientLightRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &ambientLight{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("ambientLightRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *ambientLight) build(e *WLD) error {
	return nil
}
