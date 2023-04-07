package mds

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/model/geo"
)

func (e *MDS) MaterialAdd(name string, shaderName string) error {
	name = strings.ToLower(name)

	if shaderName == "" {
		shaderName = "Opaque_MaxCB1.fx"
	}
	name = strings.ToLower(name)
	for _, mat := range e.materials {
		if mat.Name != name {
			continue
		}

		return nil
	}
	e.materials = append(e.materials, &geo.Material{
		Name:       name,
		ShaderName: shaderName,
		Properties: []*geo.Property{},
	})
	return nil
}

func (e *MDS) MaterialPropertyAdd(materialName string, propertyName string, category uint32, value string) error {
	materialName = strings.ToLower(materialName)
	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}
		o.Properties = append(o.Properties, &geo.Property{
			Name:     propertyName,
			Category: category,
			Value:    strings.ToLower(value),
		})
		return nil
	}
	return fmt.Errorf("materialName not found: '%s' (%d)", materialName, len(e.materials))
}

func (e *MDS) FaceAdd(index *geo.UIndex3, materialName string, flag uint32) error {
	materialName = strings.ToLower(materialName)
	if materialName == "" || strings.HasPrefix(materialName, "empty_") {
		e.triangles = append(e.triangles, &geo.Triangle{
			Index:        index,
			MaterialName: materialName,
			Flag:         flag,
		})
		return nil
	}

	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}

		e.triangles = append(e.triangles, &geo.Triangle{
			Index:        index,
			MaterialName: materialName,
			Flag:         flag,
		})
		return nil
	}

	return fmt.Errorf("materialName not found: '%s' (%d)", materialName, len(e.materials))
}
