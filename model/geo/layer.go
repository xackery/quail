package geo

import "fmt"

// Layer is a layer
type Layer struct {
	Name   string
	Entry0 string
	Entry1 string
}

// String returns a string representation of the layer
func (e *Layer) String() string {
	return fmt.Sprintf("%s %s %s", e.Name, e.Entry0, e.Entry1)
}
