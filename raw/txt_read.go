package raw

import (
	"io"
)

// Txt is a text file
type Txt struct {
	MetaFileName string
	Data         string
}

// Identity notes this is a txt file
func (txt *Txt) Identity() string {
	return "txt"
}

// Read takes data
func (txt *Txt) Read(r io.ReadSeeker) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	txt.Data = string(data)
	return nil

}

// SetFileName sets the name of the file
func (txt *Txt) SetFileName(name string) {
	txt.MetaFileName = name
}

// FileName returns the name of the file
func (txt *Txt) FileName() string {
	return txt.MetaFileName
}
