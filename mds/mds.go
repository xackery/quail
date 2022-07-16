package mds

// https://github.com/Zaela/EQGModelImporter/blob/master/src/mds.cpp

import (
	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

// MDS is a zon file struct
type MDS struct {
	// name is used as an identifier
	name string
	// path is used for relative paths when looking for flat file texture references
	path string
	// eqg is used as an alternative to path when loading data from a eqg file
	eqg                common.Archiver
	materials          []*common.Material
	vertices           []*common.Vertex
	faces              []*common.Face
	bones              []*bone
	files              []common.Filer
	gltfMaterialBuffer map[string]*uint32
	gltfBoneBuffer     map[int]uint32
}

type bone struct {
	name          string
	next          int32
	childrenCount uint32
	childIndex    int32
	pivot         *math32.Vector3
	rot           *math32.Vector4
	scale         *math32.Vector3
}

func New(name string, path string) (*MDS, error) {
	e := &MDS{
		name: name,
		path: path,
	}
	return e, nil
}

func NewEQG(name string, eqg common.Archiver) (*MDS, error) {
	e := &MDS{
		name: name,
		eqg:  eqg,
	}
	return e, nil
}

func (e *MDS) SetName(value string) {
	e.name = value
}

func (e *MDS) SetPath(value string) {
	e.path = value
}
