package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

type Lit struct {
	Entries []*common.RGBA
}

// Decode will decode a lit
func (lit *Lit) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	lightCount := dec.Uint32()
	for i := 0; i < int(lightCount); i++ {
		lit.Entries = append(lit.Entries, &common.RGBA{
			R: dec.Uint8(),
			G: dec.Uint8(),
			B: dec.Uint8(),
			A: dec.Uint8(),
		})
	}
	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	return nil
}
