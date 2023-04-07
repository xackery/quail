package lit

import (
	"fmt"
	"os"

	"github.com/xackery/quail/dump"
)

// BlenderExport exports a LIT file to a directory for use in blender
func (e *LIT) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	if len(e.lights) > 0 {
		lw, err := os.Create(fmt.Sprintf("%s/light.txt", path))
		if err != nil {
			return fmt.Errorf("create light.txt: %w", err)
		}

		defer lw.Close()
		lw.WriteString("rgba\n")
		for _, light := range e.lights {
			lw.WriteString(dump.Str(light) + "\n")
		}
	}

	return nil
}
