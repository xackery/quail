package export

import (
	"fmt"

	qgltf "github.com/xackery/quail/gltf"
)

// GLTFExport exports a provided mod file to gltf format
func (e *Export) GLTFExport(doc *qgltf.GLTF) error {
	if e.model == nil {
		return fmt.Errorf("no model loaded")
	}
	return e.model.GLTFExport(doc)
}
