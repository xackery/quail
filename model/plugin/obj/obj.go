// obj is a 3d model export format
package obj

import (
	"fmt"

	"github.com/xackery/quail/model/geo"
)

type ObjData struct {
	Name      string
	Materials []*geo.Material
	Vertices  []*geo.Vertex
	Triangles []*geo.Triangle
}

type ObjRequest struct {
	Data       *ObjData
	ObjPath    string
	MtlPath    string
	MattxtPath string
}

func (e *ObjData) String() string {
	return fmt.Sprintf("&{Materials (%d):[%+v]\n  Vertices (%d):[%+v]\n  Triangles (%d):[%+v]\n}", len(e.Materials), e.Materials, len(e.Vertices), e.Vertices, len(e.Triangles), e.Triangles)
}

// objCache contains temporary data needed to convert obj to eq mesh format
type objCache struct {
	vertices     []*geo.Vector3
	normals      []*geo.Vector3
	uvs          []*geo.Vector2
	vertexLookup map[string]int
}
