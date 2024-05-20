package raw

import (
	"io"
)

// Edd contations particle definitions used by prt
// examples are in eq root, actoremittersnew.edd, environmentemittersnew.edd, spellsnew.edd
type Edd struct {
	MetaFileName string `yaml:"file_name"`
}

// Identity returns the type of the struct
func (edd *Edd) Identity() string {
	return "edd"
}

func (edd *Edd) Read(r io.ReadSeeker) error {
	return nil
}

// SetFileName sets the name of the file
func (edd *Edd) SetFileName(name string) {
	edd.MetaFileName = name
}

// FileName returns the name of the file
func (edd *Edd) FileName() string {
	return edd.MetaFileName
}
