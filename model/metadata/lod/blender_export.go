package lod

import (
	"fmt"
	"os"

	"github.com/xackery/quail/dump"
)

// BlenderExport exports a LIT file to a directory for use in blender
func (e *LOD) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	if len(e.lods) > 0 {
		lw, err := os.Create(fmt.Sprintf("%s/lod.txt", path))
		if err != nil {
			return fmt.Errorf("create lod.txt: %w", err)
		}
		defer lw.Close()
		lw.WriteString("category model_name distance\n")
		for _, le := range e.lods {
			lw.WriteString(dump.Str(le.Category) + " ")
			lw.WriteString(dump.Str(le.ObjectName) + " ")
			lw.WriteString(dump.Str(le.Distance) + "\n")
		}
	}

	return nil
}
