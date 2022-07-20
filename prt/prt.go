// prt contains particle rendering settings
package prt

import "github.com/xackery/quail/common"

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

func (e *PRT) ParticleRenders() []*common.ParticleRender {
	return e.particles
}
