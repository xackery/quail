package export

import (
	"fmt"
)

// BlenderExport exports a provided exporter to a dir
func (e *Export) BlenderExport(dir string) error {
	if e.model == nil {
		return fmt.Errorf("no model loaded")
	}
	return e.model.BlenderExport(dir)
}
