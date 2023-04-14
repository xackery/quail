package zon

import (
	"fmt"
	"os"
)

func (e *ZON) BlenderExport(dir string) error {
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

	err = e.objectManager.WriteFile(fmt.Sprintf("%s/object.txt", path))
	if err != nil {
		return fmt.Errorf("objectManager.WriteFile: %w", err)
	}

	curPath := fmt.Sprintf("%s/model.txt", path)
	if len(e.models) > 0 {
		ow, err := os.Create(curPath)
		if err != nil {
			return fmt.Errorf("create file %s: %w", curPath, err)
		}
		defer ow.Close()

		ow.WriteString("name|base_name\n")

		for _, o := range e.models {
			ow.WriteString(fmt.Sprintf("%s|%s\n", o.name, o.baseName))
		}
	}

	return nil
}
