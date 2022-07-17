package gltf

import (
	"fmt"
	"io"

	"github.com/qmuntal/gltf"
)

func (e *GLTF) Export(w io.Writer) error {
	for i := uint32(0); i < uint32(len(e.doc.Nodes)); i++ {
		e.doc.Scenes[0].Nodes = append(e.doc.Scenes[0].Nodes, i)
	}

	for _, buff := range e.doc.Buffers {
		buff.EmbeddedResource()
	}

	enc := gltf.NewEncoder(w)
	enc.AsBinary = false
	enc.SetJSONIndent("", "\t")
	err := enc.Encode(e.doc)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
