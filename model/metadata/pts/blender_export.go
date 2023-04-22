package pts

import (
	"fmt"
	"os"
)

// BlenderExport exports a PTS file to a directory for use in blender
func (e *PTS) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	err = e.particleManager.BlenderExport(path)
	if err != nil {
		return fmt.Errorf("particleManager.WriteFile: %w", err)
	}

	return nil
}
