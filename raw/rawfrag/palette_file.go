package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDefaultPaletteFile is DefaultPaletteFile in libeq, empty in openzone, DEFAULTPALETTEFILE in wld
type WldFragDefaultPaletteFile struct {
	nameRef    int32  `yaml:"name_ref"`
	NameLength uint16 `yaml:"name_length"`
	FileName   string `yaml:"file_name"`
}

func (e *WldFragDefaultPaletteFile) FragCode() int {
	return FragCodeDefaultPaletteFile
}

func (e *WldFragDefaultPaletteFile) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint16(e.NameLength)
	enc.String(e.FileName)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDefaultPaletteFile) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.NameLength = dec.Uint16()
	e.FileName = dec.StringFixed(int(e.NameLength))
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragDefaultPaletteFile) NameRef() int32 {
	return e.nameRef
}
