package mod

import (
	"fmt"

	"github.com/g3n/engine/math32"
)

func (e *MOD) AddMaterial(name string, shaderName string) error {
	e.materials = append(e.materials, &material{
		name:       name,
		shaderName: shaderName,
	})
	return nil
}

func (e *MOD) AddMaterialProperty(materialName string, propertyName string, typeValue uint32, floatValue float32, intValue uint32) error {
	for _, o := range e.materials {
		if o.name != materialName {
			continue
		}
		o.properties = append(o.properties, &property{
			name:       propertyName,
			typeValue:  typeValue,
			floatValue: floatValue,
			intValue:   intValue,
		})
		return nil
	}
	return fmt.Errorf("materialName not found: %s", materialName)
}

func (e *MOD) AddVertex(position math32.Vector3, position2 math32.Vector3, uv math32.Vector2) error {
	e.vertices = append(e.vertices, &vertex{
		position:  position,
		position2: position2,
		uv:        uv,
	})
	return nil
}

func (e *MOD) AddTriangle(index math32.Vector3, materialName string, flag uint32) error {
	for _, o := range e.materials {
		if o.name != materialName {
			continue
		}

		e.triangles = append(e.triangles, &triangle{
			index:        index,
			materialName: materialName,
			flag:         flag,
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
