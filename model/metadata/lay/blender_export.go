package lay

import (
	"fmt"
	"os"
)

// BlenderExport exports a LAY file to a directory for use in blender
func (e *LAY) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	err = e.layerManager.BlenderExport(path)
	if err != nil {
		return fmt.Errorf("layerManager.WriteFile: %w", err)
	}

	return nil
}
