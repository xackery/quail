package wld

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// GLTFDecode imports a GLTF document
func (e *WLD) GLTFDecode(doc *gltf.Document) error {
	var err error
	for _, m := range doc.Materials {
		materialName := m.Name
		err = e.MaterialAdd(materialName, "Opaque_MaxCB1.fx")
		if err != nil {
			return fmt.Errorf("add material %s: %w", materialName, err)
		}

		if m.PBRMetallicRoughness.BaseColorTexture != nil {

			image := doc.Images[int(m.PBRMetallicRoughness.BaseColorTexture.Index)]
			if image == nil {
				return fmt.Errorf("expected image for '%s', but not found", materialName)
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

			err = e.MaterialPropertyAdd(materialName, "e_TextureDiffuse0", 2, imageName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", imageName, err)
			}
		}

	}

	for _, n := range doc.Nodes {

		nodeName := strings.ToLower(n.Name)
		isTer := false
		if nodeName == e.name {
			isTer = true
		}
		if nodeName+".ter" == e.name {
			isTer = true
		}
		if nodeName == e.name+".ter" {
			isTer = true
		}

		if n.Mesh == nil {
			return fmt.Errorf("no mesh on node '%s' found", n.Name)
		}
		m := doc.Meshes[*n.Mesh]
		if m == nil {
			return fmt.Errorf("accesing node '%s' mesh '%d' failed", n.Name, *n.Mesh)
		}
		meshName := strings.ToLower(m.Name)
		if meshName == e.name {
			isTer = true
		}
		if meshName+".ter" == e.name {
			isTer = true
		}
		if meshName == e.name+".ter" {
			isTer = true
		}

		if !isTer {
			continue
		}

		meshName = helper.BaseName(e.name) + ".ter"
		mesh := &mesh{
			name: meshName,
		}
		e.meshes = append(e.meshes, mesh)

		for _, p := range m.Primitives {
			if p.Mode != gltf.PrimitiveTriangles {
				return fmt.Errorf("primitive in mesh '%s' is mode %d, unsupported", meshName, p.Mode)
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
				err = e.TriangleAdd(meshName, [3]uint32{uint32(indices[i]), uint32(indices[i+1]), uint32(indices[i+2])}, materialName, 0)
				if err != nil {
					return fmt.Errorf("triangleAdd: %w", err)
				}
			}

			posIndex, ok := p.Attributes[gltf.POSITION]
			if !ok {
				return fmt.Errorf("primitive in mesh '%s' has no position", meshName)
			}
			positions, err := modeler.ReadPosition(doc, doc.Accessors[posIndex], [][3]float32{})
			if err != nil {
				return fmt.Errorf("readPosition: %w", err)
			}

			/*
				in game with no adjustments:
				N: +X
				E: +Z
				Top: -Y

				in editor:
				N: +Y
				E: +X
				Top: +Z
			*/
			//default -y front, +x on right

			// fiddle locations
			for i := range positions {
				if i == 0 {

				}
				//positions[i] = helper.ApplyQuaternion(positions[i], helper.EulerToQuaternion([3]float32{0, 90, 180}))
				// x90 perfect except x/y
				//positions[i] = helper.ApplyQuaternion(positions[i], [4]float32{0.7071068, 0, 0, 0.7071068})
				// x90 y90
				//positions[i] = helper.ApplyQuaternion(positions[i], [4]float32{0.5, 0.5, 0.5, 0.5})
				// x90 y270
				positions[i] = helper.ApplyQuaternion(positions[i], [4]float32{-0.5, 0.5, 0.5, -0.5})

				/*tmpPos := [3]float32{positions[i][0], positions[i][1], positions[i][2]}
				positions[i][0] = tmpPos[2]
				positions[i][1] = tmpPos[0]
				positions[i][2] = -tmpPos[1]*/
			}

			//fmt.Printf("pos: %+v\n", pos)
			normals := [][3]float32{}
			normalIndex, ok := p.Attributes[gltf.NORMAL]
			if ok {
				normals, err = modeler.ReadNormal(doc, doc.Accessors[normalIndex], [][3]float32{})
				if err != nil {
					return fmt.Errorf("readNormal: %w", err)
				}
			} //return fmt.Errorf("primitive in mesh '%s' has no normal", meshName)

			/*for i := range normals {
				tmpPos := [3]float32{normals[i][0], normals[i][1], normals[i][2]}
				normals[i][0] = tmpPos[2]
				normals[i][1] = tmpPos[0]
				normals[i][2] = -tmpPos[1]
			}*/

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
			} //return fmt.Errorf("primitive in mesh '%s' has no normal", meshName)

			//fmt.Printf("normal: %+v\n", normal)

			uvIndex, ok := p.Attributes[gltf.TEXCOORD_0]
			uvs := [][2]float32{}
			if ok {
				uvs, err = modeler.ReadTextureCoord(doc, doc.Accessors[uvIndex], [][2]float32{})
				if err != nil {
					return fmt.Errorf("readTextureCoord: %w", err)
				}
			}
			//return fmt.Errorf("primitive in mesh '%s' has no texcoord", meshName)
			//fmt.Printf("uv: %+v\n", uv)

			for i := 0; i < len(positions); i++ {
				positions[i][0] *= n.Scale[0] * 2
				positions[i][1] *= n.Scale[1] * 2
				positions[i][2] *= n.Scale[2] * 2
				normalEntry := [3]float32{}
				if len(normals) > i {
					normalEntry[0] = normals[i][0]
					normalEntry[1] = normals[i][1]
					normalEntry[2] = normals[i][2]
				}
				uvEntry := [2]float32{}
				if len(uvs) > i {
					uvEntry[0] = uvs[i][0] * n.Scale[0]
					uvEntry[1] = uvs[i][1] * n.Scale[1]
				}
				tint := &common.Tint{R: 128, G: 128, B: 128}
				//fmt.Printf("%d pos: %0.0f %0.0f %0.0f, normal: %+v, uv: %+v\n", i, posEntry[0], posEntry[1], posEntry[2], normalEntry, uvEntry)

				mesh.vertices = append(mesh.vertices, &common.Vertex{
					Position: positions[i],
					Normal:   normalEntry,
					Tint:     tint,
					Uv:       uvEntry,
					Uv2:      uvEntry,
				})
			}
		}
	}

	//https://github.com/KhronosGroup/glTF-Tutorials/blob/master/gltfTutorial/gltfTutorial_007_Animations.md
	for _, a := range doc.Animations {

		fmt.Println("animation", a.Name)
	}
	return nil
}
