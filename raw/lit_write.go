package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

func (lit *Lit) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Uint32(uint32(len(lit.Entries)))
	for _, entry := range lit.Entries {
		enc.Uint8(entry.R)
		enc.Uint8(entry.G)
		enc.Uint8(entry.B)
		enc.Uint8(entry.A)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
