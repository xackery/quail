package mds

import (
	"fmt"
	"strings"

	"github.com/qmuntal/gltf"

	"github.com/xackery/quail/common"
	qgltf "github.com/xackery/quail/gltf"
)

// GLTFEncode exports a provided mod file to gltf format
func (e *MDS) GLTFEncode(doc *qgltf.GLTF) error {
	var err error
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}

	modelName := strings.TrimSuffix(e.name, ".mds")

	meshCount := 1

	prims := make(map[*uint32]*qgltf.Primitive)

	lastDiffuseName := ""
	for _, material := range e.materials {
		matName := material.Name

		modelIndex := common.NumberEnding(matName)
		if modelIndex+1 > meshCount {
			meshCount = modelIndex + 1
		}

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
		//fmt.Printf("%+v %s|%s\n", material, textureDiffuseName, textureNormalName)

		var diffuseData []byte
		if len(textureDiffuseName) == 0 {
			textureDiffuseName = lastDiffuseName
		}
		if len(textureDiffuseName) > 0 {
			lastDiffuseName = textureDiffuseName

			diffuseData, err = e.archive.File(textureDiffuseName)
			if err != nil {
				return fmt.Errorf("diffuse file %s: %w", textureDiffuseName, err)
			}

		}

		var normalData []byte
		if len(textureNormalName) > 0 {
			normalData, err = e.archive.File(textureNormalName)
			if err != nil {
				return fmt.Errorf("normal file %s: %w", textureNormalName, err)
			}
		}
		_, err = doc.MaterialAdd(material, diffuseData, normalData)
		if err != nil {
			return fmt.Errorf("MaterialAdd %s: %w", material.Name, err)
		}
	}

	for i := 0; i < meshCount; i++ {
		meshName := fmt.Sprintf("%s_%02d", modelName, i)
		fmt.Println("adding mesh", meshName)
		mesh := &gltf.Mesh{Name: meshName}
		meshIndex := doc.MeshAdd(mesh)
		doc.NodeAdd(&gltf.Node{
			Name: meshName,
			Mesh: meshIndex,
		})
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

	tmpCache := make(map[string]bool)
	fmt.Println(len(e.faces), "faces")
	// ******** PRIM GENERATION *****
	for _, o := range e.faces {
		matName := o.MaterialName
		/*if strings.HasPrefix(matName, e.name+"_") {
			matName = fmt.Sprintf("c_%s_s02_m01", e.name)
		}*/

		matIndex := doc.Material(matName)

		tmpCache[matName] = true
		if matIndex == nil {
			val := uint32(0)
			matIndex = &val
		}
		meshName := modelName + "_00"

		prim, ok := prims[matIndex]
		if !ok {
			prim = qgltf.NewPrimitive()
			prim.MaterialIndex = matIndex

			prim.MeshName = meshName
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
			prim.Joints = e.joints
			prim.Weights = e.weights
		}
	}

	for matName := range tmpCache {
		fmt.Println("matCache", matName)
	}

	for _, prim := range prims {
		meshName := fmt.Sprintf("%s_00", modelName)
		err = doc.PrimitiveAdd(meshName, prim)
		if err != nil {
			return fmt.Errorf("primitiveAdd: %w", err)
		}

	}

	for i := 0; i < meshCount; i++ {
		baseMeshName := fmt.Sprintf("%s_00", modelName)
		meshName := fmt.Sprintf("%s_%02d", modelName, i)
		err = doc.PrimitiveClone(baseMeshName, meshName, i)
		if err != nil {
			return fmt.Errorf("primitive clone %d: %w", i, err)
		}
	}

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
