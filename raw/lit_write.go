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
		enc.Uint8(entry[0])
		enc.Uint8(entry[1])
		enc.Uint8(entry[2])
		enc.Uint8(entry[3])
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
