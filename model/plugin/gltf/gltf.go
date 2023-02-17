// gltf handles GLTF file format related operations
package gltf

import (
	"github.com/qmuntal/gltf"
	"github.com/xackery/quail/common"
)

type GLTF struct {
	doc             *gltf.Document
	meshes          map[string]*meshEntry
	materials       map[string]*uint32
	lights          map[string]*lightEntry
	particleRenders []*common.ParticleRender
	particlePoints  []*common.ParticlePoint
	//	gltfBoneBuffer map[int]uint32
}

func New() (*GLTF, error) {
	e := &GLTF{
		meshes:    make(map[string]*meshEntry),
		lights:    make(map[string]*lightEntry),
		materials: make(map[string]*uint32),
		doc:       gltf.NewDocument(),
	}
	return e, nil
}
