// tog contains terrain details
package tog

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
)

type TOG struct {
	name    string
	archive common.ArchiveReadWriter
	objects []*Object
}

type Object struct {
	Name     string
	Position [3]float32
	Rotation [3]float32
	Scale    float32
	FileType string
	FileName string
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, archive common.ArchiveReadWriter) (*TOG, error) {
	e := &TOG{
		name:    name,
		archive: archive,
	}
	return e, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, archive common.ArchiveReadWriter, file string) (*TOG, error) {
	e := &TOG{
		name:    name,
		archive: archive,
	}
	data, err := archive.File(file)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", file, err)
	}
	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return e, nil
}
