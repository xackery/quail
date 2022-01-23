package mod

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

func (e *MOD) ImportObj(objReader io.Reader, mtlReader io.Reader) error {
	dec, err := obj.DecodeReader(objReader, mtlReader)
	if err != nil {
		return fmt.Errorf("decodeReader: %w", err)
	}
	for i, o := range dec.Materials {
		err = e.AddMaterial(o.Name, "Opaque_MaxCB1.fx")
		if err != nil {
			return fmt.Errorf("add material %s: %w", i, err)
		}
		if o.MapKd != "" {
			data, err := ioutil.ReadFile(o.MapKd)
			if err != nil {
				return fmt.Errorf("readfile %s: %w", o.MapKd, err)
			}
			fe, err := common.NewFileEntry(o.MapKd, data)
			if err != nil {
				return fmt.Errorf("new file %s: %w", o.MapKd, err)
			}
			e.files = append(e.files, fe)
			// TODO: add string based values?
			//e.AddMaterialProperty(o.Name, "e_TextureDiffuse0", 2, )
		}
	}

	for i, o := range dec.Objects {
		// this is returning 13, original is 24
		for j, f := range o.Faces {
			err = e.AddVertex(
				math32.Vector3{X: float32(f.Vertices[0]), Y: float32(f.Vertices[1]), Z: float32(f.Vertices[2])},
				math32.Vector3{X: float32(f.Normals[0]), Y: float32(f.Normals[1]), Z: float32(f.Normals[2])},
				math32.Vector2{X: float32(f.Uvs[0]), Y: float32(f.Uvs[1])},
			)
			if err != nil {
				return fmt.Errorf("add vertex object %d face %d: %w", i, j, err)
			}
		}
	}

	return nil
}
