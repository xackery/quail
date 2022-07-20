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

	for _, m := range doc.Meshes {
		meshName := strings.ToLower(m.Name)
		if strings.Contains(meshName, "ter_") || strings.Contains(meshName, ".ter") || meshName == e.name {
			if !strings.HasSuffix(meshName, ".ter") {
				meshName += ".ter"
			}
			ml, err := ter.New(meshName, e.archive)
			if err != nil {
				return fmt.Errorf("ter new: %w", err)
			}
			err = ml.GLTFImport(doc)
			if err != nil {
				return fmt.Errorf("ter import: %w", err)
			}

			buf := &bytes.Buffer{}
			err = ml.Encode(buf)
			if err != nil {
				return fmt.Errorf("ter encode: %w", err)
			}
			e.archive.WriteFile(meshName, buf.Bytes())
			e.terrains = append(e.terrains, ml)
			e.models = append(e.models, &Model{
				baseName: helper.BaseName(meshName),
				name:     meshName,
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
			e.objects = append(e.objects, &Object{
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
		e.objects = append(e.objects, &Object{
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
