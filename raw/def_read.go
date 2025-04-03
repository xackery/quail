package raw

import (
	"encoding/base64"
	"io"
)

type Def struct {
	MetaFileName string
	Data         string
}

// Identity returns the type of the struct
func (e *Def) Identity() string {
	return "def"
}

func (e *Def) Read(r io.ReadSeeker) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	e.Data = base64.StdEncoding.EncodeToString(data)
	return nil
}

// SetFileName sets the name of the file
func (e *Def) SetFileName(name string) {
	e.MetaFileName = name
}

// FileName returns the name of the file
func (e *Def) FileName() string {
	return e.MetaFileName
}
