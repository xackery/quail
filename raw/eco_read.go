package raw

import (
	"io"
)

type Eco struct {
	MetaFileName string
}

// Identity returns the type of the struct
func (e *Eco) Identity() string {
	return "eco"
}

func (e *Eco) Read(r io.ReadSeeker) error {

	return nil
}

// SetFileName sets the name of the file
func (e *Eco) SetFileName(name string) {
	e.MetaFileName = name
}

// FileName returns the name of the file
func (e *Eco) FileName() string {
	return e.MetaFileName
}
