package raw

import (
	"encoding/base64"
	"fmt"
	"io"
)

func (bmp *Bmp) Write(w io.Writer) error {
	data, err := base64.StdEncoding.DecodeString(bmp.Data)
	if err != nil {
		return fmt.Errorf("png decode: %w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("png write: %w", err)
	}
	return nil
}
