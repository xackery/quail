package zon

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/g3n/engine/math32"
	"github.com/qmuntal/gltf"

	qgltf "github.com/xackery/quail/gltf"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
)

func (e *ZON) GLTFEncode(doc *qgltf.GLTF) error {
	var err error
	if doc == nil {
		doc, err = qgltf.New()
		if err != nil {
			return fmt.Errorf("new: %w", err)
		}
	}

	for _, model := range e.models {
		modelData, err := e.archive.File(model.name)
		if err != nil {
			return fmt.Errorf("model file %s: %w", model.name, err)
		}

		switch filepath.Ext(model.name) {
		case ".ter":
			baseName := strings.TrimPrefix(helper.BaseName(model.name), "ter_")
			e, err := ter.New(baseName, e.archive)
			if err != nil {
				return fmt.Errorf("ter.NewEQG: %w", err)
			}
			err = e.Decode(bytes.NewReader(modelData))
			if err != nil {
				return fmt.Errorf("ter decode %s: %w", baseName, err)
			}
			err = e.GLTFEncode(doc)
			if err != nil {
				return fmt.Errorf("ter gltf %s: %w", baseName, err)
			}
		case ".mod":
			baseName := strings.TrimPrefix(helper.BaseName(model.name), "ter_")
			e, err := mod.New(baseName, e.archive)
			if err != nil {
				return fmt.Errorf("mod new: %w", err)
			}
			err = e.Decode(bytes.NewReader(modelData))
			if err != nil {
				continue
				//return fmt.Errorf("mod decode %s: %w", baseName, err)
			}
			err = e.GLTFEncode(doc)
			if err != nil {
				return fmt.Errorf("mod gltf %s: %w", baseName, err)
			}
		case ".mds":
			baseName := strings.TrimPrefix(helper.BaseName(model.name), "ter_")
			e, err := mds.New(baseName, e.archive)
			if err != nil {
				return fmt.Errorf("mds new: %w", err)
			}
			err = e.Decode(bytes.NewReader(modelData))
			if err != nil {
				return fmt.Errorf("mds decode %s: %w", baseName, err)
			}
			err = e.GLTFEncode(doc)
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

		baseName := helper.BaseName(obj.name)

		index, err := doc.MeshIndex(baseName)
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

	for _, light := range e.lights {
		doc.LightAdd(light.name, light.color, light.radius, "directional", 0)
	}
	return nil
}
