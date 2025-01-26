package raw

import (
	"io"
)

type Tog struct {
	MetaFileName string
	Entries      []*TogEntry
}

// Identity returns the type of the struct
func (tog *Tog) Identity() string {
	return "Tog"
}

type TogEntry struct {
	Position [3]float32
	Rotation [3]float32
	Scale    float32
	Name     string
	FileType string
}

func (tog *Tog) Read(r io.ReadSeeker) error {
	return nil
}

// SetFileName sets the name of the file
func (tog *Tog) SetFileName(name string) {
	tog.MetaFileName = name
}

// FileName returns the name of the file
func (tog *Tog) FileName() string {
	return tog.MetaFileName
}
