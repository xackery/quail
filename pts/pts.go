// pts is an EverQuest file containing particle location and attachment data
package pts

import "github.com/xackery/quail/common"

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

func (e *PTS) ParticlePoints() []*common.ParticlePoint {
	return e.particles
}
