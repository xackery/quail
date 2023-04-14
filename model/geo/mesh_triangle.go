package geo

import (
	"fmt"
	"io"
)

// Triangle refers to the index of 3 vertices and maps it to a flag and material
type Triangle struct {
	Index        *UIndex3
	MaterialName string
	Flag         uint32
}

// NewTriangle returns a new triangle
func NewTriangle() *Triangle {
	return &Triangle{
		Index: &UIndex3{},
	}
}

// WriteHeader writes the header for a triangle
func (e *Triangle) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("index|flag|material_name\n")
	return err
}

// WriteString writes a triangle to a string writer
func (e *Triangle) Write(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("%s|%d|%s\n",
		e.Index,
		e.Flag,
		e.MaterialName))
	return err
}
