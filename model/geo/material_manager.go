package geo

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/xackery/quail/helper"
)

// MaterialManager manages materials
type MaterialManager struct {
	materials []*Material
}

// WriteFile writes all materials to a file
func (e *MaterialManager) WriteFile(material_path string, property_path string) error {
	pw, err := os.Create(property_path)
	if err != nil {
		return fmt.Errorf("create file %s: %w", property_path, err)
	}
	defer pw.Close()
	property := &Property{}
	err = property.WriteHeader(pw)
	if err != nil {
		return fmt.Errorf("write property header: %w", err)
	}

	mw, err := os.Create(material_path)
	if err != nil {
		return fmt.Errorf("create file %s: %w", material_path, err)
	}
	defer mw.Close()
	material := &Material{}
	err = material.WriteHeader(mw)
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, m := range e.materials {
		err = m.Write(mw)
		if err != nil {
			return fmt.Errorf("write material: %w", err)
		}
		for _, p := range m.Properties {
			err = p.Write(pw)
			if err != nil {
				return fmt.Errorf("write property: %w", err)
			}
		}
	}
	return nil
}

// ReadFile reads a material file
func (e *MaterialManager) ReadFile(material_path string, property_path string) error {
	r, err := os.Open(material_path)
	if err != nil {
		return fmt.Errorf("open %s: %w", material_path, err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			return fmt.Errorf("invalid material.txt (expected 3 records) line %d: %s", lineNumber, line)
		}
		material := &Material{
			Name:       parts[0],
			Flag:       helper.AtoU32(parts[1]),
			ShaderName: parts[2],
		}
		e.materials = append(e.materials, material)
	}

	r, err = os.Open(property_path)
	if err != nil {
		return fmt.Errorf("open %s: %w", property_path, err)
	}
	scanner = bufio.NewScanner(r)
	lineNumber = 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			return fmt.Errorf("invalid material_property.txt (expected 4 records) line %d: %s", lineNumber, line)
		}
		isFound := false
		for _, material := range e.materials {
			if material.Name != parts[0] {
				continue
			}
			isFound = true
			material.Properties = append(material.Properties, &Property{
				Name:     parts[1],
				Value:    parts[2],
				Category: helper.AtoU32(parts[3]),
			})
		}
		if !isFound {
			return fmt.Errorf("material_property.txt material not found: %s", parts[0])
		}
	}
	r.Close()
	return nil
}

// Add adds a material to the list
func (e *MaterialManager) Add(name string, shaderName string) error {
	name = strings.ToLower(name)
	if shaderName == "" {
		shaderName = "Opaque_MaxCB1.fx"
	}

	e.materials = append(e.materials, &Material{
		Name:       name,
		ShaderName: shaderName,
		Properties: []*Property{},
	})
	return nil
}

// PropertyAdd adds a property to a material
func (e *MaterialManager) PropertyAdd(materialName string, propertyName string, category uint32, value string) error {
	materialName = strings.ToLower(materialName)
	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}
		o.Properties = append(o.Properties, &Property{
			MaterialName: materialName,
			Name:         propertyName,
			Category:     category,
			Value:        value,
		})
		return nil
	}
	return fmt.Errorf("materialName not found: '%s' (%d)", materialName, len(e.materials))
}

// Count returns the number of materials
func (e *MaterialManager) Count() int {
	return len(e.materials)
}

// Inspect dumps to stdout information about materials and properties
func (e *MaterialManager) Inspect() {
	fmt.Println(len(e.materials), "materials:")
	for i, material := range e.materials {
		fmt.Printf("  %d %s\n", i, material.Name)
	}
}

// ByID returns a material by id
func (e *MaterialManager) ByID(id int) (*Material, bool) {
	if id == -1 {
		return nil, false
	}
	if id >= len(e.materials) {
		fmt.Printf("id '%d' is out of range (%d is max)\n", id, len(e.materials))
		return nil, false
	}
	return e.materials[id], true
}

// SortByName sorts materials by name
func (e *MaterialManager) SortByName() {
	sort.Sort(MaterialByName(e.materials))
}

// Materials returns all materials
func (e *MaterialManager) Materials() []*Material {
	return e.materials
}
