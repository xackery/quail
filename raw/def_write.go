package raw

import (
	"encoding/base64"
	"fmt"
	"io"
)

func (e *Def) Write(w io.Writer) error {
	data, err := base64.StdEncoding.DecodeString(e.Data)
	if err != nil {
		return fmt.Errorf("def decode: %w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("def write: %w", err)
	}
	return nil
}
