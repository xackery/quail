package mds

import (
	"fmt"
	"strings"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// GLTFDecode imports a GLTF document
func (e *MDS) GLTFDecode(doc *gltf.Document) error {
	var err error
	for _, m := range doc.Materials {

		fmt.Println(m.Name)
		err = e.MaterialAdd(m.Name, "Opaque_MaxCB1.fx")
		if err != nil {
			return fmt.Errorf("add material %s: %w", m.Name, err)
		}

		if m.PBRMetallicRoughness.BaseColorTexture != nil {

			image := doc.Images[int(m.PBRMetallicRoughness.BaseColorTexture.Index)]
			if image == nil {
				return fmt.Errorf("expected image for '%s', but not found", m.Name)
			}

			bv := doc.BufferViews[int(*image.BufferView)]
			if bv == nil {
				return fmt.Errorf("texture '%s' expected buffer view %d, but not found", image.Name, *image.BufferView)
			}

			ext := ""
			if strings.ToLower(image.MimeType) == "image/png" {
				ext = ".png"
			}

			if ext == "" {
				return fmt.Errorf("unsupported mimetype in gltf image '%s'", image.Name)
			}

			imageName := strings.ToLower(image.Name)
			if !strings.HasSuffix(imageName, ext) {
				imageName += ext
			}

			data, err := modeler.ReadBufferView(doc, bv)
			if err != nil {
				return fmt.Errorf("readBufferView %d: %w", *image.BufferView, err)
			}

			err = e.archive.WriteFile(imageName, data)
			if err != nil {
				return fmt.Errorf("writeFile: %w", err)
			}

			err = e.MaterialPropertyAdd(m.Name, "e_TextureDiffuse0", 2, imageName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", imageName, err)
			}
		}

	}

	for _, n := range doc.Nodes {

		if n.Mesh == nil {
			continue
			//TODO: add skin mesh detection/bone node
			//return fmt.Errorf("no mesh on node '%s' found", n.Name)
		}
		m := doc.Meshes[*n.Mesh]
		if m == nil {
			return fmt.Errorf("accesing node '%s' mesh '%d' failed", n.Name, *n.Mesh)
		}
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

			bones := [][4]uint16{}
			jointIndex, ok := p.Attributes[gltf.JOINTS_0]
			if ok {
				bones, err = modeler.ReadJoints(doc, doc.Accessors[jointIndex], [][4]uint16{})
				if err != nil {
					return fmt.Errorf("readJoints: %w", err)
				}
			}

			weights := [][4]float32{}
			weightIndex, ok := p.Attributes[gltf.WEIGHTS_0]
			if ok {
				weights, err = modeler.ReadWeights(doc, doc.Accessors[weightIndex], [][4]float32{})
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

			tints, err := modeler.ReadColor(doc, doc.Accessors[p.Attributes[gltf.COLOR_0]], [][4]uint8{})
			if err != nil {
				return fmt.Errorf("readTint: %w", err)
			}

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
				vertex := &common.Vertex{}
				positions[i] = helper.ApplyQuaternion(positions[i], [4]float32{-0.5, 0.5, 0.5, -0.5})
				positions[i][0] *= n.Scale[0]
				positions[i][1] *= n.Scale[1]
				positions[i][2] *= n.Scale[2]
				vertex.Position = positions[i]

				vertex.Normal = normals[i]

				uvs[i][0] = uvs[i][0] * n.Scale[0]
				uvs[i][1] = uvs[i][1] * n.Scale[1]
				vertex.Uv = uvs[i]

				vertex.Tint = tints[i]

				if len(bones) > i {
					vertex.Joint = bones[i]
				} else {
					vertex.Joint = [4]uint16{}
				}

				if len(weights) > i {
					vertex.Weight = weights[i]
				} else {
					vertex.Weight = [4]float32{}
				}

				e.vertices = append(e.vertices, vertex)
			}
		}
	}

	//https://github.com/KhronosGroup/glTF-Tutorials/blob/master/gltfTutorial/gltfTutorial_007_Animations.md
	for _, a := range doc.Animations {

		fmt.Println("animation", a.Name)
	}
	return nil
}
