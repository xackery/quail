package mod

import (
	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

// MOD is a zon file struct
type MOD struct {
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

func New(name string, path string) (*MOD, error) {
	e := &MOD{
		name: name,
		path: path,
	}
	return e, nil
}

func NewEQG(name string, eqg common.Archiver) (*MOD, error) {
	e := &MOD{
		name: name,
		eqg:  eqg,
	}
	return e, nil
}

func (e *MOD) SetName(value string) {
	e.name = value
}

func (e *MOD) SetPath(value string) {
	e.path = value
}
