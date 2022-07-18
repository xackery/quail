package zon

import (
	"fmt"
	"strings"

	"regexp"
)

var (
	numEndingRegex = regexp.MustCompile("[0-9]+$")
)

func (e *ZON) AddObject(modelName string, name string, position [3]float32, rotation [3]float32, scale float32) error {
	modelName = strings.ToLower(modelName)
	name = strings.ToLower(name)
	baseName := baseName(name)
	isModelFound := false
	for _, m := range e.models {
		if m.name == modelName {
			isModelFound = true
			break
		}
		if baseName == m.baseName {
			isModelFound = true
			break
		}
	}
	if !isModelFound {
		return fmt.Errorf("modelName %s not found", modelName)
	}
	e.objects = append(e.objects, &Object{
		name:        name,
		modelName:   modelName,
		translation: position,
		rotation:    rotation,
		scale:       scale,
	})
	return nil
}

func (e *ZON) AddModel(name string) error {
	name = strings.ToLower(name)
	for _, m := range e.models {
		if m.name == name {
			return nil
		}
	}
	e.models = append(e.models, &model{
		name:     name,
		baseName: baseName(name),
	})
	return nil
}

func (e *ZON) AddRegion(name string, center [3]float32, unknown [3]float32, extent [3]float32) error {
	name = strings.ToLower(name)
	e.regions = append(e.regions, &Region{
		name:    name,
		center:  center,
		unknown: unknown,
		extent:  extent,
	})
	return nil
}

func (e *ZON) AddLight(name string, position [3]float32, color [3]float32, radius float32) error {
	name = strings.ToLower(name)
	e.lights = append(e.lights, &Light{
		name:     name,
		position: position,
		color:    color,
		radius:   radius,
	})
	return nil
}

func baseName(in string) string {
	if strings.Contains(in, ".") {
		in = in[0:strings.Index(in, ".")]
	}

	in = numEndingRegex.ReplaceAllString(in, "")
	in = strings.TrimSuffix(in, "_")
	return in
}
