package wld

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/model/geo"
)

func (e *WLD) MaterialAdd(name string, shaderName string) error {
	name = strings.ToLower(name)
	if shaderName == "" {
		shaderName = "Opaque_MaxCB1.fx"
	}
	e.materials = append(e.materials, &geo.Material{
		Name:       name,
		ShaderName: shaderName,
		Properties: []*geo.Property{},
	})
	return nil
}

func (e *WLD) MaterialPropertyAdd(materialName string, propertyName string, category uint32, value string) error {
	materialName = strings.ToLower(materialName)
	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}
		o.Properties = append(o.Properties, &geo.Property{
			Name:     propertyName,
			Category: category,
			Value:    value,
		})
		return nil
	}
	return fmt.Errorf("materialName not found: '%s' (%d)", materialName, len(e.materials))
}

func (e *WLD) TriangleAdd(meshName string, index *geo.UIndex3, materialName string, flag uint32) error {

	materialName = strings.ToLower(materialName)

	var mesh *mesh
	for i := 0; i < len(e.meshes); i++ {
		if e.meshes[i].name == meshName {
			mesh = e.meshes[i]
			break
		}
	}
	if mesh == nil {
		return fmt.Errorf("mesh %s not found", meshName)
	}
	if materialName == "" || strings.HasPrefix(materialName, "empty_") {
		mesh.triangles = append(mesh.triangles, &geo.Triangle{
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

		mesh.triangles = append(mesh.triangles, &geo.Triangle{
			Index:        index,
			MaterialName: materialName,
			Flag:         flag,
		})
		return nil
	}

	return fmt.Errorf("materialName not found: '%s' (%d)", materialName, len(e.materials))
}
