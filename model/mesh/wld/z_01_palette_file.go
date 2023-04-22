package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// 0x01 1
type paletteFile struct {
	NameRef    int32
	NameLength uint16
	FileName   string
}

func (e *WLD) paletteFileRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &paletteFile{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.NameLength = dec.Uint16()
	def.FileName = dec.StringFixed(int(def.NameLength))
	if dec.Error() != nil {
		return fmt.Errorf("paletteFileRead: %s", dec.Error())
	}

	log.Debugf("paletteFile: %+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *paletteFile) build(e *WLD) error {
	return nil
}

func (e *WLD) paletteFileWrite(w io.Writer, fragmentOffset int) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	def := e.fragments[fragmentOffset].(*paletteFile)
	enc.Int32(def.NameRef)
	enc.Uint16(def.NameLength)
	enc.String(def.FileName)
	if enc.Error() != nil {
		return fmt.Errorf("paletteFileWrite: %s", enc.Error())
	}
	return nil
}
