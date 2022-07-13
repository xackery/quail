package ter

import (
	"fmt"

	"github.com/g3n/engine/math32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

// GLTFImport takes a provided GLTF path and loads relative data as a mod
func (e *TER) GLTFImport(path string) error {
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
				return fmt.Errorf("primitive in mesh '%s' is mode %d, unsupported", m.Name, p.Mode)
			}

			materialName := ""
			if p.Material != nil {
				materialName = doc.Materials[*p.Material].Name
			}

			indices, err := modeler.ReadIndices(doc, doc.Accessors[*p.Indices], []uint32{})
			if err != nil {
				return fmt.Errorf("readIndices: %w", err)
			}

			for i := 0; i < len(indices); i += 3 {
				err = e.FaceAdd([3]uint32{uint32(indices[i]), uint32(indices[i+1]), uint32(indices[i+2])}, materialName, 0)
				if err != nil {
					return fmt.Errorf("addTriangle: %w", err)
				}
			}

			posIndex, ok := p.Attributes[gltf.POSITION]
			if !ok {
				return fmt.Errorf("primitive in mesh '%s' has no position", m.Name)
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
			} //return fmt.Errorf("primitive in mesh '%s' has no normal", m.Name)

			//fmt.Printf("normal: %+v\n", normal)

			posIndex, ok = p.Attributes[gltf.TEXCOORD_0]
			uv := [][2]float32{}
			if ok {
				uv, err = modeler.ReadTextureCoord(doc, doc.Accessors[posIndex], [][2]float32{})
				if err != nil {
					return fmt.Errorf("readTextureCoord: %w", err)
				}
			} //return fmt.Errorf("primitive in mesh '%s' has no texcoord", m.Name)
			//fmt.Printf("uv: %+v\n", uv)

			for i := 0; i < len(pos); i++ {
				posEntry := math32.NewVector3(pos[i][0], pos[i][1], pos[i][2])
				normalEntry := math32.NewVec3()
				if len(normal) > i {
					normalEntry.X = normal[i][0]
					normalEntry.Y = normal[i][1]
					normalEntry.Z = normal[i][2]
				}
				uvEntry := math32.NewVec2()
				if len(uv) > i {
					uvEntry.X = uv[i][0]
					uvEntry.Y = uv[i][1]
				}
				err = e.VertexAdd(posEntry, normalEntry, uvEntry)
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
