package raw

import (
	"encoding/base64"
	"fmt"
	"io"
)

// Png takes a raw PNG type and converts it to an image.Image friendly format
type Png struct {
	MetaFileName string
	Data         string
}

// Identity returns the type of the struct
func (png *Png) Identity() string {
	return "png"
}

func (png *Png) Read(r io.ReadSeeker) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("png readall: %w", err)
	}
	png.Data = base64.StdEncoding.EncodeToString(data)

	return nil
}

// SetFileName sets the name of the file
func (png *Png) SetFileName(name string) {
	png.MetaFileName = name
}

// FileName returns the name of the file
func (png *Png) FileName() string {
	return png.MetaFileName
}
