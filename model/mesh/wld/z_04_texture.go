package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// texture is typically an animated sprite sheet
type texture struct {
	NameRef        int32
	Flags          uint32
	TextureCount   uint32
	TextureCurrent uint32
	Sleep          uint32
	TextureRefs    []uint32
}

func (e *WLD) textureRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &texture{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.Flags = dec.Uint32()
	def.TextureCount = dec.Uint32()
	if def.Flags&0x20 != 0 {
		def.TextureCurrent = dec.Uint32()
	}
	if def.Flags&0x08 != 0 && def.Flags&0x10 != 0 {
		def.Sleep = dec.Uint32()
	}
	for i := 0; i < int(def.TextureCount); i++ {
		def.TextureRefs = append(def.TextureRefs, dec.Uint32())
	}
	if dec.Error() != nil {
		return fmt.Errorf("textureRead: %s", dec.Error())
	}

	log.Debugf("texture: %+v\n", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *texture) build(e *WLD) error {
	return nil
}

func (e *WLD) textureWrite(w io.Writer, fragmentOffset int) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	def := e.fragments[fragmentOffset].(*texture)
	enc.Int32(def.NameRef)
	enc.Uint32(def.Flags)
	enc.Uint32(def.TextureCount)
	if def.Flags&0x20 != 0 {
		enc.Uint32(def.TextureCurrent)
	}
	if def.Flags&0x08 != 0 && def.Flags&0x10 != 0 {
		enc.Uint32(def.Sleep)
	}
	for _, textureRef := range def.TextureRefs {
		enc.Uint32(textureRef)
	}
	if enc.Error() != nil {
		return fmt.Errorf("textureWrite: %s", enc.Error())
	}
	return nil
}
