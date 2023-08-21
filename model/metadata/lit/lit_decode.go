package lit

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// Decode will decode a lit
func Decode(lits []*common.RGBA, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	lightCount := dec.Uint32()
	for i := 0; i < int(lightCount); i++ {
		lits = append(lits, &common.RGBA{
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
