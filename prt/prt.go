// prt contains particle rendering settings
package prt

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
)

//https://github.com/Zaela/EQGWeaponModelImporter/blob/master/src/prt.cpp

type PRT struct {
	name      string
	archive   common.ArchiveReader
	particles []*common.ParticleRender
}

func New(name string, archive common.ArchiveReader) (*PRT, error) {
	return &PRT{
		name:    name,
		archive: archive,
	}, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, archive common.ArchiveReadWriter, file string) (*PRT, error) {
	e := &PRT{
		name:    name,
		archive: archive,
	}
	data, err := archive.File(file)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", file, err)
	}
	err = e.Load(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	return e, nil
}

func (e *PRT) ParticleRenders() []*common.ParticleRender {
	return e.particles
}
