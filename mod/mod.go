package mod

import (
	"github.com/xackery/quail/common"
)

// MOD is a zon file struct
type MOD struct {
	// name is used as an identifier
	name string
	// path is used for relative paths when looking for flat file texture references
	path string
	// eqg is used as an alternative to path when loading data from a eqg file
	eqg             common.Archiver
	materials       []*common.Material
	vertices        []*common.Vertex
	faces           []*common.Face
	bones           []*Bone
	boneAssignments []*boneAssignment
	files           []common.Filer
}

type Bone struct {
	Delay       int32
	Translation [3]float32
	Rotation    [4]float32
	Scale       [3]float32
}

type boneAssignment struct {
	unknown [9]uint32
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
