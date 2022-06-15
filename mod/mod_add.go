package mod

import (
	"fmt"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

func (e *MOD) AddMaterial(name string, shaderName string) error {
	e.materials = append(e.materials, &common.Material{
		Name:       name,
		ShaderName: shaderName,
		Properties: []*common.Property{},
	})
	return nil
}

func (e *MOD) AddMaterialProperty(materialName string, propertyName string, category uint32, value string) error {
	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}
		o.Properties = append(o.Properties, &common.Property{
			Name:     propertyName,
			Category: category,
			Value:    value,
		})
		return nil
	}
	return fmt.Errorf("materialName not found: %s", materialName)
}

func (e *MOD) AddVertex(position math32.Vector3, position2 math32.Vector3, uv math32.Vector2) error {
	e.vertices = append(e.vertices, &common.Vertex{
		Position: position,
		Normal:   position2,
		Uv:       uv,
	})
	return nil
}

func (e *MOD) AddTriangle(index math32.Vector3, materialName string, flag uint32) error {
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

	return fmt.Errorf("materialName not found: %s", materialName)
}

func (e *MOD) AddBone(name string, unknown [13]float32) error {
	e.bones = append(e.bones, &bone{
		name:    name,
		unknown: unknown,
	})
	return nil
}

func (e *MOD) AddBoneAssignment(unknown [9]uint32) error {
	e.boneAssignments = append(e.boneAssignments, &boneAssignment{
		unknown: unknown,
	})
	return nil
}
