package geo

// Mesh is a mesh, used by WLD
type Mesh struct {
	Name      string
	Vertices  []*Vertex
	Triangles []*Triangle
}
