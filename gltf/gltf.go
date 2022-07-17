package gltf

import (
	"github.com/qmuntal/gltf"
)

type GLTF struct {
	doc            *gltf.Document
	meshes         map[string]*meshEntry
	materials      map[string]*uint32
	gltfBoneBuffer map[int]uint32
}

func New() (*GLTF, error) {
	e := &GLTF{
		meshes:    make(map[string]*meshEntry),
		materials: make(map[string]*uint32),
		doc:       gltf.NewDocument(),
	}
	return e, nil
}
