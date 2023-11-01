package raw

import (
	"encoding/base64"
	"fmt"
	"io"
)

func (png *Png) Write(w io.Writer) error {
	data, err := base64.StdEncoding.DecodeString(png.Data)
	if err != nil {
		return fmt.Errorf("png decode: %w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("png write: %w", err)
	}

	return nil
}
