package gltf

import "github.com/xackery/quail/common"

func (e *GLTF) ParticleRenderAdd(particle *common.ParticleRender) error {
	e.particleRenders = append(e.particleRenders, particle)
	return nil
}

func (e *GLTF) ParticlePointAdd(particle *common.ParticlePoint) error {
	e.particlePoints = append(e.particlePoints, particle)
	return nil
}
