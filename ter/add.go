package ter

import "fmt"

func (e *TER) AddMaterial(name string, shaderName string) error {
	e.materials = append(e.materials, &material{
		name:       name,
		shaderName: shaderName,
	})
	return nil
}

func (e *TER) AddMaterialProperty(materialName string, propertyName string, typeValue uint32, floatValue float32, intValue uint32) error {
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
