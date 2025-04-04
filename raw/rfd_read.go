package raw

import (
	"io"
)

type Rfd struct {
	MetaFileName string
}

// Identity returns the type of the struct
func (e *Rfd) Identity() string {
	return "rdf"
}

func (e *Rfd) Read(r io.ReadSeeker) error {

	return nil
}

// SetFileName sets the name of the file
func (e *Rfd) SetFileName(name string) {
	e.MetaFileName = name
}

// FileName returns the name of the file
func (e *Rfd) FileName() string {
	return e.MetaFileName
}
