package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// TODO: nameref is always 0, textureref is always 0, flags seem like textureref
type textureRef struct {
	NameRef    int32
	TextureRef int16
	Flags      uint32
}

func (e *WLD) textureRefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &textureRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.TextureRef = dec.Int16()
	def.Flags = dec.Uint32()
	if dec.Error() != nil {
		return fmt.Errorf("textureRefRead: %s", dec.Error())
	}
	log.Debugf("texture: %+v\n", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *textureRef) build(e *WLD) error {
	return nil
}

func (e *WLD) textureRefWrite(w io.Writer, fragmentOffset int) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	def := e.fragments[fragmentOffset].(*textureRef)
	enc.Int32(def.NameRef)
	enc.Int16(def.TextureRef)
	enc.Uint32(def.Flags)
	if enc.Error() != nil {
		return fmt.Errorf("textureRefWrite: %s", enc.Error())
	}
	return nil
}
