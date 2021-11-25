package s3d

import (
	"fmt"
	"io"

	"github.com/qmuntal/gltf"
)

// ImportGltf imports a gltf model to a wld fragment type
func (wld *Wld) ImportGltf(r io.ReadSeeker) error {
	doc := gltf.Document{}
	dec := gltf.NewDecoder(r)
	err := dec.Decode(&doc)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
