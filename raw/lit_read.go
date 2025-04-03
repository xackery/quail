package raw

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type Lit struct {
	MetaFileName string
	Entries      [][4]uint8
}

// Identity returns the type of the struct
func (lit *Lit) Identity() string {
	return "lit"
}

// Decode will read a lit
func (lit *Lit) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	lightCount := dec.Uint32()
	for i := 0; i < int(lightCount); i++ {
		lit.Entries = append(lit.Entries, [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()})
	}

	pos := dec.Pos()
	endPos, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}
	if pos < endPos {
		remaining := dec.Bytes(int(endPos - pos))
		if !bytes.Equal(remaining, []byte{0x0, 0x0, 0x0, 0x0}) {
			fmt.Printf("remaining bytes: %s\n", hex.Dump(remaining))
			return fmt.Errorf("%d bytes remaining (%d total)", endPos-pos, endPos)
		}
	}
	if pos > endPos {
		return fmt.Errorf("read past end of file")
	}
	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	return nil
}

// SetFileName sets the name of the file
func (lit *Lit) SetFileName(name string) {
	lit.MetaFileName = name
}

// FileName returns the name of the file
func (lit *Lit) FileName() string {
	return lit.MetaFileName
}
