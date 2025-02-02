package raw

import (
	"bytes"
	"io"
)

// Unk is a fallback type for raw data we can't parse
type Unk struct {
	MetaFileName string
	Data         []byte
}

// Identity returns the type of the struct
func (unk *Unk) Identity() string {
	return "unk"
}

func (unk *Unk) Read(r io.ReadSeeker) error {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}
	unk.Data = buf.Bytes()

	return nil
}

// SetFileName sets the name of the file
func (unk *Unk) SetFileName(name string) {
	unk.MetaFileName = name
}

// FileName returns the name of the file
func (unk *Unk) FileName() string {
	return unk.MetaFileName
}
