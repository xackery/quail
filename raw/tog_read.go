package raw

import (
	"io"
)

type Tog struct {
	Entries []*TogEntry
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
