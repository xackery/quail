package raw

import (
	"io"
)

type DatFe struct {
	MetaFileName string
	Version      uint32
}

func (e *DatFe) Identity() string {
	return "datfe"
}

func (e *DatFe) Read(r io.ReadSeeker) error {

	return nil
}

// SetName sets the name of the file
func (e *DatFe) SetFileName(name string) {
	e.MetaFileName = name
}

func (e *DatFe) FileName() string {
	return e.MetaFileName
}
