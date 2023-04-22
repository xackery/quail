package geo

import (
	"fmt"
	"io"
)

// Layer is a layer
type Layer struct {
	Name   string
	Entry0 string
	Entry1 string
}

// String returns a string representation of the layer
func (e *Layer) String() string {
	return fmt.Sprintf("%s|%s|%s", e.Name, e.Entry0, e.Entry1)
}

// WriteHeader writes the header to a file
func (e *Layer) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("name|entry0|entry1\n")
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	return nil
}

// Write writes the layer to a file
func (e *Layer) Write(w io.StringWriter) error {
	_, err := w.WriteString(e.String() + "\n")
	if err != nil {
		return fmt.Errorf("write layer: %w", err)
	}
	return nil
}
