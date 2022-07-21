package wld

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/common"
)

func (e *WLD) triangleAdd(meshName string, index [3]uint32, materialName string, flag uint32) error {

	var mesh *mesh
	for i := range e.meshes {
		if e.meshes[i].name == meshName {
			break
		}
	}
	if mesh == nil {
		return fmt.Errorf("mesh %s not found", meshName)
	}

	materialName = strings.ToLower(materialName)
	if materialName == "" || strings.HasPrefix(materialName, "empty_") {
		mesh.triangles = append(mesh.triangles, &common.Triangle{
			Index:        index,
			MaterialName: materialName,
			Flag:         flag,
		})
		return nil
	}

	for _, o := range e.materials {
		if o.Name != materialName {
			continue
		}

		mesh.triangles = append(mesh.triangles, &common.Triangle{
			Index:        index,
			MaterialName: materialName,
			Flag:         flag,
		})
		return nil
	}

	return fmt.Errorf("materialName not found: '%s' (%d)", materialName, len(e.materials))
}
