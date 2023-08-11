package quail

import (
	"github.com/xackery/quail/quail/def"
)

type Quail struct {
	Meshes     []*def.Mesh
	Animations []*def.Animation
	Zone       *def.Zone
}

// New returns a new Quail instance
func New() *Quail {
	return &Quail{}
}
