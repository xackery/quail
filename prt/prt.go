// prt contains particle rendering settings
package prt

import "github.com/xackery/quail/common"

//https://github.com/Zaela/EQGWeaponModelImporter/blob/master/src/prt.cpp

type PRT struct {
	name      string
	archive   common.Archiver
	particles []*common.ParticleEntry
}

func New(name string, archive common.Archiver) (*PRT, error) {
	return &PRT{
		name:    name,
		archive: archive,
	}, nil
}

func (e *PRT) Particles() []*common.ParticleEntry {
	return e.particles
}
