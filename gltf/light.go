package gltf

import (
	"fmt"

	"github.com/qmuntal/gltf"
)

type lightEntry struct {
	index        *uint32
	Name         string     `json:"name"`
	Color        [3]float32 `json:"color"`
	Intensity    float32    `json:"intensity"`
	EmissionType string     `json:"type"`
	Range        float32    `json:"range,omitempty"`
}

// LightAdd adds a light, emissionType includes directional, point, spot
func (e *GLTF) LightAdd(name string, color [3]float32, intensity float32, emissionType string, distance float32) *uint32 {
	entry, ok := e.lights[name]
	if ok {
		return entry.index
	}

	le := &lightEntry{
		Name:         name,
		Color:        color,
		Intensity:    intensity,
		EmissionType: emissionType,
		Range:        distance,
	}

	type lightNode struct {
		Light *uint32 `json:"light"`
	}

	index := gltf.Index(uint32(len(e.lights)) - 1)
	le.index = index

	ext := gltf.Extensions{}
	ext["KHR_lights_punctual"] = &lightNode{Light: index}

	e.doc.Nodes = append(e.doc.Nodes, &gltf.Node{
		Name:       name,
		Extensions: ext,
	})

	e.lights[name] = le
	return index
}

func (e *GLTF) Light(name string) (*uint32, error) {
	entry, ok := e.lights[name]
	if ok {
		return entry.index, nil
	}
	return nil, fmt.Errorf("'%s' not found", name)
}
