package raw

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type DatIw struct {
	MetaFileName string
	Version      uint32
}

func (e *DatIw) Identity() string {
	return "datiw"
}

func (e *DatIw) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

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
	return nil
}

// SetName sets the name of the file
func (e *DatIw) SetFileName(name string) {
	e.MetaFileName = name
}

func (e *DatIw) FileName() string {
	return e.MetaFileName
}
