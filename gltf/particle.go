package gltf

import "github.com/xackery/quail/common"

func (e *GLTF) PaticleAdd(particle *common.ParticleEntry) error {
	e.particles = append(e.particles, particle)
	return nil
}
