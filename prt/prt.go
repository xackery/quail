// prt contains particle rendering settings
package prt

import "github.com/xackery/quail/common"

//https://github.com/Zaela/EQGWeaponModelImporter/blob/master/src/prt.cpp

type PRT struct {
	name      string
	archive   common.Archiver
	particles []*particleEntry
}

type particleEntry struct {
	id              uint32 //id is actorsemittersnew.edd
	id2             uint32
	name            string
	unknownA        [5]uint32 //Pretty sure last 3 have something to do with durations
	duration        uint32
	unknownB        uint32
	unknownFFFFFFFF int32
	unknownC        uint32
}

func New(name string, archive common.Archiver) (*PRT, error) {
	return &PRT{
		name:    name,
		archive: archive,
	}, nil
}
