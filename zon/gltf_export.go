package zon

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/g3n/engine/math32"
	"github.com/qmuntal/gltf"

	qgltf "github.com/xackery/quail/gltf"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
)

func (e *ZON) GLTFExport(doc *qgltf.GLTF) error {
	var err error
	if doc == nil {
		doc, err = qgltf.New()
		if err != nil {
			return fmt.Errorf("new: %w", err)
		}
	}

	for _, model := range e.models {
		modelData, err := e.eqg.File(model.name)
		if err != nil {
			return fmt.Errorf("model file %s: %w", model.name, err)
		}

		switch filepath.Ext(model.name) {
		case ".ter":
			baseName := strings.TrimPrefix(baseName(model.name), "ter_")
			e, err := ter.NewEQG(baseName, e.eqg)
			if err != nil {
				return fmt.Errorf("ter.NewEQG: %w", err)
			}
			err = e.Load(bytes.NewReader(modelData))
			if err != nil {
				return fmt.Errorf("ter load %s: %w", baseName, err)
			}
			err = e.GLTF(doc)
			if err != nil {
				return fmt.Errorf("ter gltf %s: %w", baseName, err)
			}
		case ".mod":
			baseName := strings.TrimPrefix(baseName(model.name), "ter_")
			e, err := mod.NewEQG(baseName, e.eqg)
			if err != nil {
				return fmt.Errorf("mod.NewEQG: %w", err)
			}
			err = e.Load(bytes.NewReader(modelData))
			if err != nil {
				continue
				//return fmt.Errorf("mod load %s: %w", baseName, err)
			}
			err = e.GLTFExport(doc)
			if err != nil {
				return fmt.Errorf("mod gltf %s: %w", baseName, err)
			}
		case ".mds":
			baseName := strings.TrimPrefix(baseName(model.name), "ter_")
			e, err := mds.NewEQG(baseName, e.eqg)
			if err != nil {
				return fmt.Errorf("mds.NewEQG: %w", err)
			}
			err = e.Load(bytes.NewReader(modelData))
			if err != nil {
				return fmt.Errorf("mds load %s: %w", baseName, err)
			}
			err = e.GLTFExport(doc)
			if err != nil {
				return fmt.Errorf("mds gltf %s: %w", baseName, err)
			}

		default:
			return fmt.Errorf("unsupported model: %s", model.name)
		}

	}

	//math32.NewVec3().ApplyQuaternion(q *math32.Quaternion)

	for _, obj := range e.objects {
		if strings.HasPrefix(obj.name, "ter_") {
			//fmt.Println("skipping",obj.name)
			continue
		}

		baseName := baseName(obj.name)

		index, err := doc.Mesh(baseName)
		if err != nil {
			fmt.Println("failed", err)
			continue
			//TODO: fix
			//return fmt.Errorf("mesh: %w", err)
		}
		qRot := math32.NewQuaternion(0, 0, 0, 0).SetFromEuler(&math32.Vector3{X: obj.rotation[0], Y: obj.rotation[1], Z: obj.rotation[2]})
		doc.NodeAdd(&gltf.Node{
			Name:        obj.name,
			Mesh:        index,
			Translation: obj.translation,
			Rotation:    [4]float32{qRot.X, qRot.Y, qRot.Z, qRot.W},
			Scale:       [3]float32{obj.scale, obj.scale, obj.scale},
		})
	}

	return nil
}
