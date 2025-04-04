package raw

import (
	"encoding/base64"
	"io"
)

type Obg struct {
	MetaFileName string
	Data         string
}

// Identity returns the type of the struct
func (e *Obg) Identity() string {
	return "obg"
}

func (e *Obg) Read(r io.ReadSeeker) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	e.Data = base64.StdEncoding.EncodeToString(data)
	return nil
}

// SetFileName sets the name of the file
func (e *Obg) SetFileName(name string) {
	e.MetaFileName = name
}

// FileName returns the name of the file
func (e *Obg) FileName() string {
	return e.MetaFileName
}
