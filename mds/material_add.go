package mds

import (
	"fmt"
	"strings"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
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
	e.materials = append(e.materials, &common.Material{
		Name:       name,
		ShaderName: shaderName,
		Properties: []*common.Property{},
	})
	return nil
}

func (e *MDS) MaterialPropertyAdd(materialName string, propertyName string, category uint32, value string) error {
	materialName = strings.ToLower(materialName)
	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}
		o.Properties = append(o.Properties, &common.Property{
			Name:     propertyName,
			Category: category,
			Value:    strings.ToLower(value),
		})
		return nil
	}
	return fmt.Errorf("materialName not found: '%s' (%d)", materialName, len(e.materials))
}

func (e *MDS) VertexAdd(position *math32.Vector3, normal *math32.Vector3, tint *common.Tint, uv *math32.Vector2, uv2 *math32.Vector2) error {
	e.vertices = append(e.vertices, &common.Vertex{
		Position: position,
		Normal:   normal,
		Tint:     tint,
		Uv:       uv,
		Uv2:      uv2,
	})
	return nil
}

func (e *MDS) FaceAdd(index [3]uint32, materialName string, flag uint32) error {
	materialName = strings.ToLower(materialName)
	if materialName == "" || strings.HasPrefix(materialName, "empty_") {
		e.triangles = append(e.triangles, &common.Triangle{
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

		e.triangles = append(e.triangles, &common.Triangle{
			Index:        index,
			MaterialName: materialName,
			Flag:         flag,
		})
		return nil
	}

	return fmt.Errorf("materialName not found: '%s' (%d)", materialName, len(e.materials))
}
