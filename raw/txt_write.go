package raw

import (
	"fmt"
	"io"
)

// Write will write a txt file
func (txt *Txt) Write(w io.Writer) error {
	_, err := w.Write([]byte(txt.Data))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
