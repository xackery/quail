package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type first struct {
	NameRef int32
}

func (e *WLD) firstRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &first{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	if dec.Error() != nil {
		return fmt.Errorf("firstRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *first) build(e *WLD) error {
	return nil
}

func (e *WLD) firstWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
