package mod

import (
	"fmt"

	"github.com/g3n/engine/math32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func (e *MOD) ImportGLTF(path string) error {
	doc, err := gltf.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	for _, m := range doc.Materials {
		//TODO: add _mat.txt parsing
		err = e.AddMaterial(m.Name, "Opaque_MaxCB1.fx")
		if err != nil {
			return fmt.Errorf("add material %s: %w", m.Name, err)
		}
	}
	for _, m := range doc.Meshes {
		for _, p := range m.Primitives {
			if p.Mode != gltf.PrimitiveTriangles {
				return fmt.Errorf("primitive in mesh %s is mode %d, unsupported", m.Name, p.Mode)
			}

			posIndex, ok := p.Attributes[gltf.POSITION]
			if !ok {
				return fmt.Errorf("primitive in mesh %s has no position", m.Name)
			}
			pos, err := modeler.ReadPosition(doc, doc.Accessors[posIndex], [][3]float32{})
			if err != nil {
				return fmt.Errorf("readPosition: %w", err)
			}

			fmt.Printf("pos: %+v\n", pos)

			posIndex, ok = p.Attributes[gltf.NORMAL]
			if !ok {
				return fmt.Errorf("primitive in mesh %s has no normal", m.Name)
			}
			normal, err := modeler.ReadNormal(doc, doc.Accessors[posIndex], [][3]float32{})
			if err != nil {
				return fmt.Errorf("readNormal: %w", err)
			}

			fmt.Printf("normal: %+v\n", normal)

			posIndex, ok = p.Attributes[gltf.TEXCOORD_0]
			if !ok {
				return fmt.Errorf("primitive in mesh %s has no texcoord", m.Name)
			}
			uv, err := modeler.ReadTextureCoord(doc, doc.Accessors[posIndex], [][2]float32{})
			if err != nil {
				return fmt.Errorf("readTextureCoord: %w", err)
			}
			fmt.Printf("uv: %+v\n", uv)

			for i := 0; i < len(pos); i++ {
				err = e.AddVertex(math32.Vector3{X: pos[i][0], Y: pos[i][1], Z: pos[i][2]},
					math32.Vector3{X: normal[i][0], Y: normal[i][1], Z: normal[i][2]},
					math32.Vector2{X: uv[i][0], Y: uv[i][1]})
				if err != nil {
					return fmt.Errorf("add vertex: %w", err)
				}
			}
		}
	}
	return nil
}
