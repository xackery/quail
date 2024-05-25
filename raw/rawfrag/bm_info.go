package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

// WldFragBMInfo is BmInfo in libeq, Texture Bitmap Names in openzone, FRAME and BMINFO in wld, BitmapName in lantern
type WldFragBMInfo struct {
	NameRef      int32    `yaml:"name_ref"`
	TextureNames []string `yaml:"texture_names"`
}

func (e *WldFragBMInfo) FragCode() int {
	return FragCodeBMInfo
}

func (e *WldFragBMInfo) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	start := enc.Pos()

	enc.Int32(e.NameRef)
	enc.Int32(int32(len(e.TextureNames) - 1))
	enc.StringLenPrefixUint16(string(helper.WriteStringHash(strings.Join(e.TextureNames, ""))))

	diff := enc.Pos() - start
	paddingSize := (4 - diff%4) % 4
	enc.Bytes(make([]byte, paddingSize))
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragBMInfo) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	textureCount := dec.Int32()

	for i := 0; i < int(textureCount+1); i++ {
		nameLength := dec.Uint16()
		e.TextureNames = append(e.TextureNames, helper.ReadStringHash((dec.Bytes(int(nameLength)))))
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
