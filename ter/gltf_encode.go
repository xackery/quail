package ter

import (
	"fmt"
	"strings"

	"github.com/qmuntal/gltf"

	qgltf "github.com/xackery/quail/gltf"
)

// GLTFEncode exports a provided mod file to gltf format
func (e *TER) GLTFEncode(doc *qgltf.GLTF) error {
	var err error
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}

	meshName := strings.TrimSuffix(e.name, ".mds")

	mesh := &gltf.Mesh{
		Name: meshName,
	}

	prims := make(map[*uint32]*qgltf.Primitive)

	lastDiffuseName := ""
	for _, material := range e.materials {

		textureDiffuseName := ""
		textureNormalName := ""
		for _, p := range material.Properties {
			if p.Category != 2 {
				continue
			}
			if strings.EqualFold(p.Name, "e_texturediffuse0") {
				textureDiffuseName = p.Value
				continue
			}

			if strings.EqualFold(p.Name, "e_texturenormal0") {
				textureNormalName = p.Value
				continue
			}
		}

		var diffuseData []byte
		if len(textureDiffuseName) == 0 {
			textureDiffuseName = lastDiffuseName
		}
		if len(textureDiffuseName) > 0 {
			lastDiffuseName = textureDiffuseName
			diffuseData, err = e.archive.File(textureDiffuseName)
			if err != nil {
				return fmt.Errorf("file %s: %w", textureDiffuseName, err)
			}
		}

		var normalData []byte
		if len(textureNormalName) > 0 {
			normalData, err = e.archive.File(textureNormalName)
			if err != nil {
				return fmt.Errorf("file %s: %w", textureNormalName, err)
			}
		}
		_, err = doc.MaterialAdd(material, diffuseData, normalData)
		if err != nil {
			return fmt.Errorf("MaterialAdd %s: %w", material.Name, err)
		}
	}

	// ******** MESH SKINNING *******
	/*var skinIndex *uint32
	for i, b := range e.bones {
		doc.Nodes = append(doc.Nodes, &gltf.Node{
			Name: b.name,
			//Translation: [3]float32{b.pivot[0], b.pivot[1], b.pivot[2]},
			Rotation: [4]float32{b.rot[0], b.rot[1], b.rot[2], b.rot[3]},
			Scale:    [3]float32{b.scale[0], b.scale[1], b.scale[2]},
		})
		//if strings.EqualFold(b.name, "ROOT_BONE") {
		//		rootNode = uint32(len(doc.Nodes) - 1)
		//}
		e.gltfBoneBuffer[i] = uint32(len(doc.Nodes) - 1)
	}

	for i, b := range e.bones {
		children := &[]uint32{}
		if b.childIndex > -1 {
			err = e.gltfBoneChildren(doc, children, int(b.childIndex))
			if err != nil {
				return fmt.Errorf("gltfBoneChildren: %w", err)
			}
		}

		fmt.Printf("%d %d %d %d children for %s: %d\n", i, b.next, b.childIndex, b.childrenCount, b.name, len(*children))
		if strings.EqualFold(b.name, "ROOT_BONE") {
			//*children = append(*children, rootNode)
			skin := &gltf.Skin{
				Name:   e.bones[0].name,
				Joints: *children,
			}
			doc.Skins = append(doc.Skins, skin)
			tmp := uint32(len(doc.Skins) - 1)
			skinIndex = &tmp
		} else {
			nodeIndex, ok := e.gltfBoneBuffer[i]
			if !ok {
				return fmt.Errorf("bone for %d not found", i)
			}
			node := doc.Nodes[int(nodeIndex)]
			node.Children = *children
		}
	}
	*/

	// ******** PRIM GENERATION *****
	for _, o := range e.triangles {
		matName := o.MaterialName
		if strings.HasPrefix(matName, e.name+"_") {
			matName = fmt.Sprintf("c_%s_s02_m01", e.name)
		}

		matIndex := doc.Material(matName)

		prim, ok := prims[matIndex]
		if !ok {
			prim = qgltf.NewPrimitive()
			prim.MaterialIndex = matIndex
			prims[matIndex] = prim
		}

		for i := 0; i < 3; i++ {
			index, ok := prim.UniqueIndices[o.Index[i]]
			if !ok {
				v := e.vertices[int(o.Index[i])]

				// x-90 y-270
				//v.Position = helper.ApplyQuaternion(v.Position, [4]float32{0.5, -0.5, 0.5, -0.5})
				// x90
				//v.Position = helper.ApplyQuaternion(v.Position, [4]float32{-0.7071068, 0, 0, 0.7071068})
				prim.Positions = append(prim.Positions, v.Position)
				prim.Normals = append(prim.Normals, [3]float32{v.Normal[0], v.Normal[1], v.Normal[2]})
				prim.Uvs = append(prim.Uvs, [2]float32{v.Uv[0], v.Uv[1]})
				prim.UniqueIndices[o.Index[i]] = uint16(len(prim.Positions) - 1)
				index = uint16(len(prim.Positions) - 1)
			}
			prim.Indices = append(prim.Indices, index)
		}

	}

	meshIndex := doc.MeshAdd(mesh)

	for _, prim := range prims {
		err = doc.PrimitiveAdd(meshName, prim)
		if err != nil {
			return fmt.Errorf("primitiveAdd: %w", err)
		}
	}

	doc.NodeAdd(&gltf.Node{
		Name: meshName,
		Mesh: meshIndex,
		//Skin: skinIndex,
	})

	for _, particle := range e.particleRenders {
		err = doc.ParticleRenderAdd(particle)
		if err != nil {
			return fmt.Errorf("ParticleRenderAdd: %w", err)
		}
	}

	for _, particle := range e.particlePoints {
		err = doc.ParticlePointAdd(particle)
		if err != nil {
			return fmt.Errorf("ParticlePointAdd: %w", err)
		}
	}

	return nil
}

/*
func (e *MOD) gltfBoneChildren(doc *gltf.Document, children *[]uint32, boneIndex int) error {

	nodeIndex, ok := e.gltfBoneBuffer[boneIndex]
	if !ok {
		return fmt.Errorf("bone %d node not found", boneIndex)
	}
	*children = append(*children, nodeIndex)

	bone := e.bones[boneIndex]
	if bone.next == -1 {
		return nil
	}

	return e.gltfBoneChildren(doc, children, int(bone.next))
}*/
