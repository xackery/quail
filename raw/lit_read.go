package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

type Lit struct {
	MetaFileName string        `yaml:"file_name"`
	Entries      []*model.RGBA `yaml:"entries"`
}

// Identity returns the type of the struct
func (lit *Lit) Identity() string {
	return "lit"
}

// Decode will read a lit
func (lit *Lit) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	lightCount := dec.Uint32()
	for i := 0; i < int(lightCount); i++ {
		lit.Entries = append(lit.Entries, &model.RGBA{
			R: dec.Uint8(),
			G: dec.Uint8(),
			B: dec.Uint8(),
			A: dec.Uint8(),
		})
	}
	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	return nil
}

// SetFileName sets the name of the file
func (lit *Lit) SetFileName(name string) {
	lit.MetaFileName = name
}

// FileName returns the name of the file
func (lit *Lit) FileName() string {
	return lit.MetaFileName
}
