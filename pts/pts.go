// pts is an EverQuest file containing particle location and attachment data
package pts

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
)

// https://github.com/Zaela/EQGWeaponModelImporter/blob/master/src/pts.cpp

type PTS struct {
	name      string
	archive   common.ArchiveReader
	particles []*common.ParticlePoint
}

func New(name string, archive common.ArchiveReader) (*PTS, error) {
	return &PTS{
		name:    name,
		archive: archive,
	}, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, archive common.ArchiveReadWriter, file string) (*PTS, error) {
	e := &PTS{
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

func (e *PTS) ParticlePoints() []*common.ParticlePoint {
	return e.particles
}
