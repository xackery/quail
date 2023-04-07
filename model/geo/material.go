package geo

import (
	"fmt"
	"strings"
)

// Property contains data about a property
type Properties []*Property

// String returns a string representation of the properties
func (e Properties) String() string {
	value := ""
	for _, p := range e {
		value += fmt.Sprintf("[%s: %s], ", p.Name, p.Value)
	}
	if len(e) > 0 {
		value = value[0 : len(value)-2]
	}
	return value
}

// RGB represents R,G,B as uint8
type RGB struct {
	R uint8
	G uint8
	B uint8
}

// String returns a string representation of the RGB
func (e *RGB) String() string {
	return fmt.Sprintf("%d,%d,%d", e.R, e.G, e.B)
}

// RGBA represents R,G,B,A as uint8
type RGBA struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// String returns a string representation of the RGBA
func (e *RGBA) String() string {
	return fmt.Sprintf("%d,%d,%d,%d", e.R, e.G, e.B, e.A)
}

// Material contains data about a material
type Material struct {
	Name       string
	ShaderName string
	Flag       uint32
	Properties Properties
}

// String returns a string representation of the material
func (e *Material) String() string {
	return fmt.Sprintf("%s %s %d %s", e.Name, e.ShaderName, e.Flag, e.Properties)
}

// MaterialByName is a slice of Material
type MaterialByName []*Material

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
