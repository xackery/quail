package quail

import (
	"fmt"
	"io"
)

func (e *Quail) WldMarshal(w io.Writer) error {
	/*
		world := &common.Wld{}
		fragIndex := 1

		// TODO: material fragments

			for _, model := range e.Models {
				mesh := &wld.Mesh{}
				name := model.Header.Name

				mesh.NameRef = int32(pos)
				mesh.Flags = 0x00014003
				mesh.MaterialListRef = uint32(0) // TODO: add proper refs
				mesh.AnimationRef = int32(0)     // TODO: add proper refs
				mesh.Vertices = model.Vertices
				world.Fragments[fragIndex] = mesh
				fragIndex++
			}
	*/
	return fmt.Errorf("not yet implemented")
}
