package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

type userData struct {
}

func (e *WLD) userDataRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &userData{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return fmt.Errorf("paletteFileRead: %s", dec.Error())
	}

	log.Debugf("userData: %+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *userData) build(e *WLD) error {
	return nil
}

func (e *WLD) userDataWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
