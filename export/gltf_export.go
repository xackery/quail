package export

import (
	"fmt"

	qgltf "github.com/xackery/quail/gltf"
)

// GLTFEncode exports a provided mod file to gltf format
func (e *Export) GLTFEncode(doc *qgltf.GLTF) error {
	if e.model == nil {
		return fmt.Errorf("no model loaded")
	}
	return e.model.GLTFEncode(doc)
}
