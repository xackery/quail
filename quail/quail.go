package quail

import "github.com/xackery/quail/common"

type Quail struct {
	Models        []*common.Model
	Animations    []*common.Animation
	Zone          *common.Zone
	materialCache map[string]*common.Material
}

// New returns a new Quail instance
func New() *Quail {
	return &Quail{
		materialCache: make(map[string]*common.Material),
	}
}
