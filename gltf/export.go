package gltf

import (
	"fmt"
	"io"

	"github.com/qmuntal/gltf"
	"github.com/xackery/quail/common"
)

func (e *GLTF) Export(w io.Writer) error {
	for i := uint32(0); i < uint32(len(e.doc.Nodes)); i++ {
		e.doc.Scenes[0].Nodes = append(e.doc.Scenes[0].Nodes, i)
	}

	isEverQuestExtension := false

	if len(e.particles) > 0 {
		if !isEverQuestExtension {
			e.doc.ExtensionsUsed = append(e.doc.ExtensionsUsed, "everquest")
			isEverQuestExtension = true
		}

		type particleEntries struct {
			Particles []*common.ParticleEntry `json:"particles"`
		}
		ext := particleEntries{}
		e.doc.Extensions = gltf.Extensions{}

		ext.Particles = e.particles

		e.doc.Extensions["everquest"] = ext
	}
	if len(e.lights) > 0 {
		e.doc.ExtensionsUsed = append(e.doc.ExtensionsUsed, "KHR_lights_punctual")

		type lightEntries struct {
			Lights []*lightEntry `json:"lights"`
		}
		ext := lightEntries{}
		e.doc.Extensions = gltf.Extensions{}

		for _, light := range e.lights {
			ext.Lights = append(ext.Lights, light)
		}
		e.doc.Extensions["KHR_lights_punctual"] = ext
	}
	for _, buff := range e.doc.Buffers {
		buff.EmbeddedResource()
	}

	enc := gltf.NewEncoder(w)
	enc.AsBinary = false
	enc.SetJSONIndent("", "\t")
	err := enc.Encode(e.doc)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
