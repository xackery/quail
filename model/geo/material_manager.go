package geo

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/log"
)

// MaterialManager manages materials
type MaterialManager struct {
	materials map[int32]*Material
}

// NewMaterialManager creates a new material manager
func NewMaterialManager() *MaterialManager {
	return &MaterialManager{
		materials: make(map[int32]*Material),
	}
}

// BlenderExport writes material.txt and material_property.txt to path
func (e *MaterialManager) BlenderExport(path string) error {
	if len(e.materials) == 0 {
		return nil
	}

	propertyPath := fmt.Sprintf("%s/material_property.txt", path)
	pw, err := os.Create(propertyPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", propertyPath, err)
	}
	defer pw.Close()
	property := MaterialProperty{}
	err = property.WriteHeader(pw)
	if err != nil {
		return fmt.Errorf("write property header: %w", err)
	}

	materialPath := fmt.Sprintf("%s/material.txt", path)
	mw, err := os.Create(materialPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", materialPath, err)
	}
	defer mw.Close()
	material := Material{}
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

// BlenderImport reads a material file
func (e *MaterialManager) BlenderImport(path string) error {
	materialPath := fmt.Sprintf("%s/material.txt", path)
	r, err := os.Open(materialPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", materialPath, err)
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
		material := Material{
			ID:         helper.AtoI32(parts[0]),
			Name:       parts[1],
			Flag:       helper.AtoU32(parts[2]),
			ShaderName: parts[3],
		}
		e.materials[material.ID] = &material
	}

	propertyPath := fmt.Sprintf("%s/material_property.txt", path)
	r, err = os.Open(propertyPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", propertyPath, err)
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
			material.Properties = append(material.Properties, MaterialProperty{
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
func (e *MaterialManager) Add(material Material) error {
	material.Name = strings.ToLower(material.Name)
	if material.ShaderName == "" {
		material.ShaderName = "Opaque_MaxCB1.fx"
	}

	e.materials[material.ID] = &material
	return nil
}

// PropertyAdd adds a property to a material
func (e *MaterialManager) PropertyAdd(materialName string, property MaterialProperty) error {
	materialName = strings.ToLower(materialName)
	for i, _ := range e.materials {
		material := e.materials[i]
		if material.Name != materialName {
			continue
		}

		material.Properties = append(e.materials[i].Properties, property)
		log.Debugf("Added property %v to material %s", property, materialName)
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
func (e *MaterialManager) ByID(id int32) (Material, bool) {
	material, ok := e.materials[id]
	return *material, ok
}

// Materials returns all materials
func (e *MaterialManager) Materials() map[int32]*Material {
	return e.materials
}
