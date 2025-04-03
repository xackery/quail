package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

func (e *DatFe) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("encoder error: %w", err)
	}

	return nil
}
