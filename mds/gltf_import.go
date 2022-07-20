package mds

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/g3n/engine/math32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/xackery/quail/common"
)

// GLTFImport imports a GLTF document
func (e *MDS) GLTFImport(doc *gltf.Document) error {
	var err error
	for _, m := range doc.Materials {
		//TODO: add _mat.txt parsing
		err = e.MaterialAdd(m.Name, "Opaque_MaxCB1.fx")
		if err != nil {
			return fmt.Errorf("add material %s: %w", m.Name, err)
		}
		err = e.MaterialPropertyAdd(m.Name, "e_TextureDiffuse0", 2, strings.ToLower(m.Name))
		if err != nil {
			return fmt.Errorf("materialPropertyAdd %s: %w", m.Name, err)
		}
	}

	/*for _, m := range doc.Skins {

	}*/

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
					return fmt.Errorf("faceAdd: %w", err)
				}
			}

			posIndex, ok := p.Attributes[gltf.POSITION]
			if !ok {
				return fmt.Errorf("primitive in mesh '%s' has no position", m.Name)
			}
			positions, err := modeler.ReadPosition(doc, doc.Accessors[posIndex], [][3]float32{})
			if err != nil {
				return fmt.Errorf("readPosition: %w", err)
			}

			/*for i, _ := range positions {
				tmp := positions[i][1]
				positions[i][1] = -positions[i][2]
				positions[i][2] = tmp

			}*/

			jointIndex, ok := p.Attributes[gltf.JOINTS_0]
			if ok {
				e.joints, err = modeler.ReadJoints(doc, doc.Accessors[jointIndex], [][4]uint16{})
				if err != nil {
					return fmt.Errorf("readJoints: %w", err)
				}
			}

			weightIndex, ok := p.Attributes[gltf.WEIGHTS_0]
			if ok {
				e.weights, err = modeler.ReadWeights(doc, doc.Accessors[weightIndex], [][4]float32{})
				if err != nil {
					return fmt.Errorf("readJoints: %w", err)
				}
			}

			//fmt.Printf("pos: %+v\n", pos)
			normals := [][3]float32{}
			normalIndex, ok := p.Attributes[gltf.NORMAL]
			if ok {
				normals, err = modeler.ReadNormal(doc, doc.Accessors[normalIndex], [][3]float32{})
				if err != nil {
					return fmt.Errorf("readNormal: %w", err)
				}
			} //return fmt.Errorf("primitive in mesh '%s' has no normal", m.Name)

			tints := &color.RGBA{255, 255, 255, 255}
			tintIndex, ok := p.Attributes[gltf.COLOR_0]
			if ok {
				tintRaw, err := modeler.ReadColor(doc, doc.Accessors[tintIndex], [][4]uint8{})
				if err != nil {
					return fmt.Errorf("readTint: %w", err)
				}
				tints.R = tintRaw[0][0]
				tints.G = tintRaw[0][1]
				tints.B = tintRaw[0][2]
				tints.A = tintRaw[0][3]
			} //return fmt.Errorf("primitive in mesh '%s' has no normal", m.Name)

			//fmt.Printf("normal: %+v\n", normal)

			uvIndex, ok := p.Attributes[gltf.TEXCOORD_0]
			uvs := [][2]float32{}
			if ok {
				uvs, err = modeler.ReadTextureCoord(doc, doc.Accessors[uvIndex], [][2]float32{})
				if err != nil {
					return fmt.Errorf("readTextureCoord: %w", err)
				}
			}
			//return fmt.Errorf("primitive in mesh '%s' has no texcoord", m.Name)
			//fmt.Printf("uv: %+v\n", uv)

			for i := 0; i < len(positions); i++ {
				posEntry := math32.NewVector3(positions[i][0], positions[i][1], positions[i][2])
				normalEntry := math32.NewVec3()
				if len(normals) > i {
					normalEntry.X = normals[i][0]
					normalEntry.Y = normals[i][1]
					normalEntry.Z = normals[i][2]
				}
				uvEntry := math32.NewVec2()
				if len(uvs) > i {
					uvEntry.X = uvs[i][0]
					uvEntry.Y = uvs[i][1]
				}
				tint := &common.Tint{R: 128, G: 128, B: 128}
				//fmt.Printf("%d pos: %0.0f %0.0f %0.0f, normal: %+v, uv: %+v\n", i, posEntry.X, posEntry.Y, posEntry.Z, normalEntry, uvEntry)
				err = e.VertexAdd(posEntry, normalEntry, tint, uvEntry, uvEntry)
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
