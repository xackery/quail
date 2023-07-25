package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type zone struct {
}

func (e *WLD) zoneRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &zone{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("zoneRead: %v", dec.Error())
	}

	//log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *zone) build(e *WLD) error {
	return nil
}

func (e *WLD) zoneWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
