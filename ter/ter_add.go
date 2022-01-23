package ter

import (
	"fmt"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

func (e *TER) AddMaterial(name string, shaderName string) error {
	e.materials = append(e.materials, &common.Material{
		Name:       name,
		ShaderName: shaderName,
		Properties: []*common.Property{},
	})
	return nil
}

func (e *TER) AddMaterialProperty(materialName string, propertyName string, typeValue uint32, floatValue float32, intValue uint32) error {
	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}
		o.Properties = append(o.Properties, &common.Property{
			Name:       propertyName,
			TypeValue:  typeValue,
			FloatValue: floatValue,
			IntValue:   intValue,
		})
		return nil
	}
	return fmt.Errorf("materialName not found: %s", materialName)
}

func (e *TER) AddVertex(position math32.Vector3, position2 math32.Vector3, uv math32.Vector2) error {
	e.vertices = append(e.vertices, &common.Vertex{
		Position: position,
		Normal:   position2,
		Uv:       uv,
	})
	return nil
}

func (e *TER) AddTriangle(index math32.Vector3, materialName string, flag uint32) error {
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
