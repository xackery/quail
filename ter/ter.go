package ter

import (
	"github.com/xackery/quail/common"
)

// TER is a zon file struct
type TER struct {
	materials []*common.Material
	vertices  []*common.Vertex
	triangles []*common.Triangle
	files     []*common.FileEntry
}
