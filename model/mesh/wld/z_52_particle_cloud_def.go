package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type particleCloudDef struct {
}

func (e *WLD) particleCloudDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &particleCloudDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("particleCloudDefRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *particleCloudDef) build(e *WLD) error {
	return nil
}

func (e *WLD) particleCloudDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
