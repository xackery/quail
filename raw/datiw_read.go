package raw

import (
	"io"
)

type DatIw struct {
	MetaFileName string
	Version      uint32
}

func (e *DatIw) Identity() string {
	return "datiw"
}

func (e *DatIw) Read(r io.ReadSeeker) error {

	return nil
}

// SetName sets the name of the file
func (e *DatIw) SetFileName(name string) {
	e.MetaFileName = name
}

func (e *DatIw) FileName() string {
	return e.MetaFileName
}
