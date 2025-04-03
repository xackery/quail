package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Encode encodes a v4 zone dat file
// https://github.com/EQEmu/zone-utilities/blob/master/src/common/eqg_v4_loader.cpp#L115
func (e *DatIw) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("encoder error: %w", err)
	}

	return nil
}
