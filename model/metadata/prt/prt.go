// prt contains particle rendering settings
package prt

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

//https://github.com/Zaela/EQGWeaponModelImporter/blob/master/src/prt.cpp

// PRT contains particle rendering settings
type PRT struct {
	name      string
	archive   archive.Reader
	particles []*geo.ParticleRender
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.Reader) (*PRT, error) {
	return &PRT{
		name:    name,
		archive: pfs,
	}, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*PRT, error) {
	e := &PRT{
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

// ParticleRenders returns a list of particle renders
func (e *PRT) ParticleRenders() []*geo.ParticleRender {
	return e.particles
}

// Name returns the name of the prt
func (e *PRT) Name() string {
	return e.name
}
