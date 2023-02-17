package ter

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/common"
)

func (e *TER) triangleAdd(index [3]uint32, materialName string, flag uint32) error {
	materialName = strings.ToLower(materialName)
	if materialName == "" || strings.HasPrefix(materialName, "empty_") {
		e.triangles = append(e.triangles, &common.Triangle{
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

		e.triangles = append(e.triangles, &common.Triangle{
			Index:        index,
			MaterialName: materialName,
			Flag:         flag,
		})
		return nil
	}

	return fmt.Errorf("materialName not found: '%s' (%d)", materialName, len(e.materials))
}
