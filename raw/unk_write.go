package raw

import (
	"fmt"
	"io"
)

func (unk *Unk) Write(w io.Writer) error {
	_, err := w.Write(unk.Data)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
