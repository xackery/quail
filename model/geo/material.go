package geo

import (
	"fmt"
	"io"
	"strings"
)

// Material contains data about a material
type Material struct {
	ID         int32
	Name       string
	ShaderName string
	Flag       uint32
	Properties MaterialProperties
}

// WriteHeader writes the header for a material file
func (e Material) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("id|material_name|flag|category\n")
	return err
}

// Write writes a material to a file
func (e Material) Write(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("%d|%s|%d|%s\n", e.ID, e.Name, e.Flag, e.ShaderName))
	return err
}

// MaterialByName is a slice of Material
type MaterialByName []Material

// Len returns the length of the slice
func (s MaterialByName) Len() int {
	return len(s)
}

// Swap swaps the elements with indexes i and j.
func (s MaterialByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the element with
func (s MaterialByName) Less(i, j int) bool {
	return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name)
}
