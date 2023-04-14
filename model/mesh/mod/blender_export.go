package mod

import (
	"fmt"
	"os"
)

// BlenderExport exports the MOD to a directory for use in Blender.
func (e *MOD) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	vw, err := os.Create(fmt.Sprintf("%s/info.txt", path))
	if err != nil {
		return fmt.Errorf("create info.txt: %w", err)
	}
	defer vw.Close()
	vw.WriteString(fmt.Sprintf("version=%d\n", e.version))

	if e.MaterialManager.Count() > 0 {
		err = e.MaterialManager.WriteFile(fmt.Sprintf("%s/material.txt", path), fmt.Sprintf("%s/material_property.txt", path))
		if err != nil {
			return fmt.Errorf("materialManager.WriteFile: %w", err)
		}
	}

	if e.particleManager.PointCount() > 0 {
		err = e.particleManager.WriteFile(fmt.Sprintf("%s/particle_point.txt", path), fmt.Sprintf("%s/particle_render.txt", path))
		if err != nil {
			return fmt.Errorf("particleManager.WriteFile: %w", err)
		}
	}

	/*for _, file := range e.files {
		fw, err := os.Create(fmt.Sprintf("%s/%s", path, file.Name()))
		if err != nil {
			return fmt.Errorf("create %s: %w", file.Name(), err)
		}
		defer fw.Close()
		fw.Write(file.Data())
	}*/

	err = e.meshManager.WriteFile(path)
	if err != nil {
		return fmt.Errorf("write meshManager: %w", err)
	}

	return nil
}
