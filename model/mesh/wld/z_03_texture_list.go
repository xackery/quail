package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// textureList 0x03 03
type textureList struct {
	NameRef      int32
	TextureNames []string
}

func (e *WLD) textureListRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &textureList{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	textureCount := dec.Int32()

	for i := 0; i < int(textureCount+1); i++ {
		nameLength := dec.Uint16()
		def.TextureNames = append(def.TextureNames, decodeStringHash(dec.Bytes(int(nameLength)))) // TODO: this actually is encoded
	}
	if dec.Error() != nil {
		return fmt.Errorf("textureListRead: %s", dec.Error())
	}

	log.Debugf("textureList%+v\n", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *textureList) build(e *WLD) error {
	return nil
}

func (e *WLD) textureListWrite(w io.Writer, fragmentOffset int) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	def := e.fragments[fragmentOffset].(*textureList)
	enc.Int32(def.NameRef)
	enc.Int32(int32(len(def.TextureNames) + 1))
	for _, textureName := range def.TextureNames {
		enc.Uint16(uint16(len(textureName)))
		enc.String(textureName)
	}
	if enc.Error() != nil {
		return fmt.Errorf("textureListWrite: %s", enc.Error())
	}
	return nil
}
