package geo

import (
	"fmt"
	"io"
)

// Vertex stores information related to a mesh
type Vertex struct {
	Position Vector3
	Normal   Vector3
	Tint     RGBA
	Uv       Vector2
	Uv2      Vector2
}

// NewVertex returns a new vertex
func NewVertex() Vertex {
	return Vertex{
		Position: Vector3{},
		Normal:   Vector3{},
		Tint:     RGBA{},
		Uv:       Vector2{},
		Uv2:      Vector2{},
	}
}

// WriteHeader writes the header for a vertex
func (e *Vertex) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("position|normal|uv|uv2|tint\n")
	return err
}

// WriteString writes a vertex to a string writer
func (e *Vertex) Write(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("%s|%s|%s|%s|%s\n",
		&Vector3{X: e.Position.Y, Y: -e.Position.X, Z: e.Position.Z},
		&Vector3{X: e.Normal.Y, Y: -e.Normal.X, Z: e.Normal.Z},
		e.Uv,
		e.Uv2,
		e.Tint))
	return err
}
