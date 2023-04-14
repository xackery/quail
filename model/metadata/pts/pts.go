// pts is an EverQuest file containing particle location and attachment data
package pts

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// https://github.com/Zaela/EQGWeaponModelImporter/blob/master/src/pts.cpp

// PTS contains particle location and attachment data
type PTS struct {
	name            string
	archive         archive.Reader
	particleManager *geo.ParticleManager
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.Reader) (*PTS, error) {
	return &PTS{
		name:            name,
		archive:         pfs,
		particleManager: &geo.ParticleManager{},
	}, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*PTS, error) {
	e := &PTS{
		name:    name,
		archive: pfs,
	}
	data, err := pfs.File(file)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", file, err)
	}
	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return e, nil
}

// Name returns the name of the file
func (e *PTS) Name() string {
	return e.name
}
