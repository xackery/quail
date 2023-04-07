package pts

import (
	"fmt"
	"os"

	"github.com/xackery/quail/dump"
)

// BlenderExport exports a PTS file to a directory for use in blender
func (e *PTS) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	if len(e.particles) > 0 {
		pw, err := os.Create(fmt.Sprintf("%s/particle_point.txt", path))
		if err != nil {
			return fmt.Errorf("create particle_point.txt: %w", err)
		}

		defer pw.Close()
		pw.WriteString("name bone translation rotation scale\n")
		for _, p := range e.particles {
			pw.WriteString(dump.Str(p.Name) + " ")
			pw.WriteString(dump.Str(p.Bone) + " ")
			pw.WriteString(dump.Str(p.Translation) + " ")
			pw.WriteString(dump.Str(p.Rotation) + " ")
			pw.WriteString(dump.Str(p.Scale) + "\n")
		}
	}

	return nil
}
