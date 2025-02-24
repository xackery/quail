package raw

import (
	"encoding/base64"
	"fmt"
	"io"
)

func (tga *Tga) Write(w io.Writer) error {
	data, err := base64.StdEncoding.DecodeString(tga.Data)
	if err != nil {
		return fmt.Errorf("tga decode: %w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("tga write: %w", err)
	}
	return nil
}
