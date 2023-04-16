package geo

import (
	"fmt"
	"io"
	"strings"
)

// Mesh is a mesh, used by WLD
type Mesh struct {
	Name       string
	Vertices   []Vertex
	Triangles  []Triangle
	Bones      []Bone
	Animations []BoneAnimation
}

// WriteHeader writes the header for a mesh file
func (e *Mesh) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("name\n")
	return err
}

// Write writes a mesh to a file
func (e *Mesh) Write(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("%s\n", strings.TrimSuffix(e.Name, "_DMSPRITEDEF")))
	return err
}
