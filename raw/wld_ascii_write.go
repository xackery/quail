package raw

import (
	"fmt"
	"io"
)

// Write writes data to a wld_ascii file
func (wld *WldAscii) Write(w io.Writer) error {
	_, err := w.Write([]byte(wld.Data))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
