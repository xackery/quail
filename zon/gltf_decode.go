package zon

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/qmuntal/gltf"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
)

// GLTFDecode imports a GLTF document
func (e *ZON) GLTFDecode(doc *gltf.Document) error {

	for _, n := range doc.Nodes {

		if n.Mesh == nil {
			return fmt.Errorf("no mesh was referred to on node %s", n.Name)
		}
		m := doc.Meshes[*n.Mesh]
		if m == nil {
			return fmt.Errorf("node %s refers to mesh %d and it was not found", n.Name, *n.Mesh)
		}

		nodeName := strings.ToLower(n.Name)
		meshName := strings.ToLower(m.Name)
		if strings.Contains(meshName, "ter_") ||
			strings.Contains(meshName, ".ter") ||
			meshName == e.name ||
			strings.Contains(nodeName, "ter_") ||
			strings.Contains(nodeName, ".ter") ||
			nodeName == e.name {
			if !strings.HasSuffix(nodeName, ".ter") {
				nodeName += ".ter"
			}
			ml, err := ter.New(nodeName, e.archive)
			if err != nil {
				return fmt.Errorf("ter new: %w", err)
			}
			err = ml.GLTFDecode(doc)
			if err != nil {
				return fmt.Errorf("ter import: %w", err)
			}

			buf := &bytes.Buffer{}
			err = ml.Encode(buf)
			if err != nil {
				return fmt.Errorf("ter encode: %w", err)
			}
			e.archive.WriteFile(nodeName, buf.Bytes())
			e.terrains = append(e.terrains, ml)
			e.models = append(e.models, &model{
				baseName: helper.BaseName(nodeName),
				name:     nodeName,
			})
			continue
		}

		if strings.HasSuffix(meshName, ".mod") {
			ml, err := mod.New(meshName, e.archive)
			if err != nil {
				return fmt.Errorf("mod new: %w", err)
			}
			err = ml.GLTFDecode(doc)
			if err != nil {
				return fmt.Errorf("mod import: %w", err)
			}
			e.mods = append(e.mods, ml)
			e.objects = append(e.objects, &object{
				modelName: meshName,
				name:      meshName,
			})
			continue
		}

		//if strings.HasSuffix(meshName, ".mds") {
		ml, err := mds.New(meshName, e.archive)
		if err != nil {
			return fmt.Errorf("mds new: %w", err)
		}
		err = ml.GLTFDecode(doc)
		if err != nil {
			return fmt.Errorf("mds import: %w", err)
		}
		e.mdses = append(e.mdses, ml)
		e.objects = append(e.objects, &object{
			modelName: meshName,
			name:      meshName,
		})
		//}
	}
	//https://github.com/KhronosGroup/glTF-Tutorials/blob/master/gltfTutorial/gltfTutorial_007_Animations.md
	for _, a := range doc.Animations {

		fmt.Println("animation", a.Name)
	}
	return nil
}
