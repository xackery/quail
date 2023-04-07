package ter

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/model/geo"
)

func (e *TER) MaterialByID(id int32) (string, error) {
	if id == -1 {
		return "", nil
	}
	if int(id) >= len(e.materials) {
		return "", fmt.Errorf("id '%d' is out of range (%d is max)", id, len(e.materials))
	}
	return e.materials[id].Name, nil
}

func (e *TER) MaterialAdd(name string, shaderName string) error {
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

func (e *TER) MaterialPropertyAdd(materialName string, propertyName string, category uint32, value string) error {
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
