package raw

import (
	"encoding/base64"
	"fmt"
	"io"
)

func (dds *Dds) Write(w io.Writer) error {
	data, err := base64.StdEncoding.DecodeString(dds.Data)
	if err != nil {
		return fmt.Errorf("dds decode: %w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("dds write: %w", err)
	}
	return nil
}
