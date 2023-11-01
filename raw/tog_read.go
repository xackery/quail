package raw

import (
	"io"
)

type Tog struct {
	MetaFileName string      `yaml:"file_name"`
	Entries      []*TogEntry `yaml:"entries"`
}

type TogEntry struct {
	Position [3]float32 `yaml:"position"`
	Rotation [3]float32 `yaml:"rotation"`
	Scale    float32    `yaml:"scale"`
	Name     string     `yaml:"name"`
	FileType string     `yaml:"file_type"`
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
