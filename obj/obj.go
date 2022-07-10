package obj

import (
	"fmt"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

type ObjData struct {
	Name      string
	Materials []*common.Material
	Vertices  []*common.Vertex
	Triangles []*common.Triangle
}

type ObjRequest struct {
	Obj        *ObjData
	ObjPath    string
	MtlPath    string
	MattxtPath string
}

func (e *ObjData) String() string {
	return fmt.Sprintf("&{Materials (%d):[%+v]\n  Vertices (%d):[%+v]\n  Triangles (%d):[%+v]\n}", len(e.Materials), e.Materials, len(e.Vertices), e.Vertices, len(e.Triangles), e.Triangles)
}

// objCache contains temporary data needed to convert obj to eq mesh format
type objCache struct {
	vertices     []math32.Vector3
	normals      []math32.Vector3
	uvs          []math32.Vector2
	vertexLookup map[string]int
}
