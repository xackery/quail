package quail

import (
	"github.com/xackery/quail/quail/def"
)

type Quail struct {
	Meshes        []*def.Mesh
	Animations    []*def.Animation
	Zone          *def.Zone
	materialCache map[string]*def.Material
}

// New returns a new Quail instance
func New() *Quail {
	return &Quail{
		materialCache: make(map[string]*def.Material),
	}
}
