package mod

import "github.com/xackery/quail/common"

// MOD is a zon file struct
type MOD struct {
	name            string
	materials       []*common.Material
	vertices        []*common.Vertex
	triangles       []*common.Triangle
	bones           []*bone
	boneAssignments []*boneAssignment
	files           []common.Filer
}

type bone struct {
	name    string
	unknown [13]float32
}

type boneAssignment struct {
	unknown [9]uint32
}

func New(name string) (*MOD, error) {
	e := &MOD{
		name: name,
	}
	return e, nil
}
