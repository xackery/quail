package geo

import (
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/helper"
)

// MaterialProperties contains a list of material properties
type MaterialProperties []MaterialProperty

// MaterialProperty contains data about a material
type MaterialProperty struct {
	MaterialName string // parent material name
	Name         string
	Category     uint32
	Value        string
}

// WriteHeader writes a material property header to a file
func (e MaterialProperty) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("material_name|property_name|value|category\n")
	return err
}

// Write writes a material property to a file
func (e MaterialProperty) Write(w io.StringWriter) error {
	value := strings.ToLower(e.Value)
	if strings.ToLower(e.Name) == "e_fshininess0" {
		val := helper.AtoF32(e.Value)
		if val > 100 {
			val = 1.0
		} else {
			val /= 100
		}
		value = fmt.Sprintf("%f", val)
	}
	_, err := w.WriteString(fmt.Sprintf("%s|%s|%s|%d\n", e.MaterialName, e.Name, value, e.Category))
	return err
}
