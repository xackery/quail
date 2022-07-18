package mod

import (
	"fmt"
	"os"
	"strings"

	"github.com/qmuntal/gltf"

	qgltf "github.com/xackery/quail/gltf"
)

// GLTFExport exports a provided mod file to gltf format
func (e *MOD) GLTFExport(doc *qgltf.GLTF) error {
	var err error
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}

	modelName := strings.TrimSuffix(e.name, ".mds")

	mesh := &gltf.Mesh{
		Name: modelName,
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
			if e.eqg != nil {
				diffuseData, err = e.eqg.File(textureDiffuseName)
				if err != nil {
					return fmt.Errorf("file %s: %w", textureDiffuseName, err)
				}
			}
			if len(diffuseData) == 0 && e.path != "" {
				diffuseData, err = os.ReadFile(fmt.Sprintf("%s/%s", e.path, textureDiffuseName))
				if err != nil {
					return fmt.Errorf("file %s: %w", textureDiffuseName, err)
				}
			}
		}

		var normalData []byte
		if len(textureNormalName) > 0 {
			if e.eqg != nil {
				normalData, err = e.eqg.File(textureNormalName)
				if err != nil {
					return fmt.Errorf("file %s: %w", textureNormalName, err)
				}
			}
			if len(normalData) == 0 && e.path != "" {
				diffuseData, err = os.ReadFile(fmt.Sprintf("%s/%s", e.path, textureDiffuseName))
				if err != nil {
					return fmt.Errorf("file %s: %w", textureDiffuseName, err)
				}
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
			//Translation: [3]float32{b.pivot.X, b.pivot.Y, b.pivot.Z},
			Rotation: [4]float32{b.rot.X, b.rot.Y, b.rot.Z, b.rot.W},
			Scale:    [3]float32{b.scale.X, b.scale.Y, b.scale.Z},
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
	for _, o := range e.faces {
		matName := o.MaterialName
		if strings.HasPrefix(matName, e.name+"_") {
			matName = fmt.Sprintf("c_%s_s02_m01", e.name)
		}

		matIndex := doc.Material(matName)
		if matIndex == nil {
			val := uint32(0)
			matIndex = &val
		}

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
				prim.Positions = append(prim.Positions, [3]float32{v.Position.X, v.Position.Y, v.Position.Z})
				prim.Normals = append(prim.Normals, [3]float32{v.Normal.X, v.Normal.Y, v.Normal.Z})
				prim.Uvs = append(prim.Uvs, [2]float32{v.Uv.X, v.Uv.Y})
				prim.UniqueIndices[o.Index[i]] = uint16(len(prim.Positions) - 1)
				index = uint16(len(prim.Positions) - 1)
			}
			prim.Indices = append(prim.Indices, index)
		}

	}

	meshIndex := doc.MeshAdd(mesh)

	for _, prim := range prims {
		err = doc.PrimitiveAdd(e.name, prim)
		if err != nil {
			return fmt.Errorf("primitiveAdd: %w", err)
		}
	}

	doc.NodeAdd(&gltf.Node{
		Name: modelName,
		Mesh: meshIndex,
		//Skin: skinIndex,
	})

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
