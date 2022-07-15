package mod

import (
	"fmt"

	"github.com/g3n/engine/math32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/xackery/quail/common"
)

// GLTFImport takes a provided GLTF path and loads relative data as a mod
func (e *MOD) GLTFImport(path string) error {
	doc, err := gltf.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	for _, m := range doc.Materials {
		//TODO: add _mat.txt parsing
		fmt.Println("add material", m.Name)
		err = e.MaterialAdd(m.Name, "Opaque_MaxCB1.fx")
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

			//fmt.Printf("pos: %+v\n", pos)

			normal := [][3]float32{}
			posIndex, ok = p.Attributes[gltf.NORMAL]
			if ok {
				normal, err = modeler.ReadNormal(doc, doc.Accessors[posIndex], [][3]float32{})
				if err != nil {
					return fmt.Errorf("readNormal: %w", err)
				}
			} else {
				for i := 0; i < len(pos); i++ {
					normal = append(normal, [3]float32{0, 0, 0})
				}
			}

			//fmt.Printf("normal: %+v\n", normal)

			uv := [][2]float32{}
			posIndex, ok = p.Attributes[gltf.TEXCOORD_0]
			if ok {
				uv, err = modeler.ReadTextureCoord(doc, doc.Accessors[posIndex], [][2]float32{})
				if err != nil {
					return fmt.Errorf("readTextureCoord: %w", err)
				}
			} else {
				for i := 0; i < len(pos); i++ {
					uv = append(uv, [2]float32{0, 0})
				}
			}

			//fmt.Printf("uv: %+v\n", uv)

			for i := 0; i < len(pos); i++ {
				tint := &common.Tint{R: 128, G: 128, B: 128}
				err = e.VertexAdd(math32.NewVector3(pos[i][0], pos[i][1], pos[i][2]),
					math32.NewVector3(normal[i][0], normal[i][1], normal[i][2]),
					tint,
					math32.NewVector2(uv[i][0], uv[i][1]),
					math32.NewVector2(uv[i][0], uv[i][1]))
				if err != nil {
					return fmt.Errorf("add vertex: %w", err)
				}
			}
		}
	}

	//https://github.com/KhronosGroup/glTF-Tutorials/blob/master/gltfTutorial/gltfTutorial_007_Animations.md
	for _, a := range doc.Animations {

		fmt.Println("animation", a.Name)
	}
	return nil
}
