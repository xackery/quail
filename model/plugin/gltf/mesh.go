package gltf

import (
	"fmt"

	"github.com/qmuntal/gltf"
)

type meshEntry struct {
	index *uint32
	mesh  *gltf.Mesh
}

func (e *GLTF) MeshAdd(mesh *gltf.Mesh) *uint32 {
	entry, ok := e.meshes[mesh.Name]
	if ok {
		return entry.index
	}
	e.doc.Meshes = append(e.doc.Meshes, mesh)
	index := gltf.Index(uint32(len(e.doc.Meshes) - 1))
	e.meshes[mesh.Name] = &meshEntry{
		index: index,
		mesh:  mesh,
	}
	//fmt.Println("added mesh", mesh.Name)
	return index
}

func (e *GLTF) MeshIndex(name string) (*uint32, error) {
	entry, ok := e.meshes[name]
	if ok {
		return entry.index, nil
	}
	return nil, fmt.Errorf("'%s' not found", name)
}
