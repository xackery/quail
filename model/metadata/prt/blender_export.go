package prt

import (
	"fmt"
	"os"

	"github.com/xackery/quail/dump"
)

// BlenderExport exports a PRT file to a directory for use in blender
func (e *PRT) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	if len(e.particles) > 0 {
		pw, err := os.Create(fmt.Sprintf("%s/particle_render.txt", path))
		if err != nil {
			return fmt.Errorf("create particle_render.txt: %w", err)
		}
		defer pw.Close()
		pw.WriteString("duration id id2 particle_point unknowna unknownb unknownc unknownffffffff\n")
		for _, p := range e.particles {
			pw.WriteString(dump.Str(p.Duration) + " ")
			pw.WriteString(dump.Str(p.ID) + " ")
			pw.WriteString(dump.Str(p.ID2) + " ")
			pw.WriteString(dump.Str(p.ParticlePoint) + " ")
			pw.WriteString(dump.Str(p.UnknownA) + " ")
			pw.WriteString(dump.Str(p.UnknownB) + " ")
			pw.WriteString(dump.Str(p.UnknownC) + " ")
			pw.WriteString(dump.Str(p.UnknownFFFFFFFF) + "\n")
		}
	}

	return nil
}
