package zon

import (
	"fmt"

	"github.com/g3n/engine/math32"
)

func (e *ZON) AddObject(modelName string, name string, position math32.Vector3, rotation math32.Vector3, scale float32) error {
	isModelFound := false
	for _, m := range e.models {
		if m.name == modelName {
			isModelFound = true
			break
		}
	}
	if !isModelFound {
		return fmt.Errorf("modelName %s not found", modelName)
	}
	e.objects = append(e.objects, &object{
		name:      name,
		modelName: modelName,
		position:  position,
		rotation:  rotation,
		scale:     scale,
	})
	return nil
}

func (e *ZON) AddModel(name string) error {
	e.models = append(e.models, &model{name: name})
	return nil
}

func (e *ZON) AddRegion(name string, center math32.Vector3, unknown math32.Vector3, extent math32.Vector3) error {
	e.regions = append(e.regions, &region{
		name:    name,
		center:  center,
		unknown: unknown,
		extent:  extent,
	})
	return nil
}

func (e *ZON) AddLight(name string, position math32.Vector3, color math32.Color, radius float32) error {
	e.lights = append(e.lights, &light{
		name:     name,
		position: position,
		color:    color,
		radius:   radius,
	})
	return nil
}
