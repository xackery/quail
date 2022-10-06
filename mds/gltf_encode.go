package mds

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qmuntal/gltf"

	"github.com/xackery/quail/common"
	qgltf "github.com/xackery/quail/gltf"
	"github.com/xackery/quail/helper"
)

// GLTFEncode exports a provided mod file to gltf format
func (e *MDS) GLTFEncode(doc *qgltf.GLTF) error {
	var err error
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}

	prefix := ""

	modelName := strings.TrimSuffix(e.name, ".mds")
	modelName = strings.TrimSuffix(modelName, ".eqg")

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

	var skinIndex *uint32
	if e.skin != nil {
		matrixIndex := doc.WriteMatrix(e.skin.InverseBindMatrices)

		skin := &gltf.Skin{
			Name:                e.skin.Name,
			InverseBindMatrices: &matrixIndex,
		}

		nodes := make(map[int]uint32)

		for i, joint := range e.skin.Joints {
			nodeIndex, ok := nodes[i]
			if !ok {
				node := &gltf.Node{
					Name:        joint.Name,
					Translation: joint.Translation,
				}
				nodeIndex = doc.NodeAdd(node)
				nodes[i] = nodeIndex
			}

			if i == 0 {
				skin.Skeleton = &nodeIndex
			} else {
				skin.Joints = append(skin.Joints, nodeIndex)
			}

		}

		tmpSkinIndex := doc.SkinAdd(skin)
		skinIndex = &tmpSkinIndex
	}

	for i := 0; i < meshCount; i++ {
		meshName := fmt.Sprintf("%s_%02d", modelName, i)

		fmt.Println("adding mesh", meshName)
		if strings.Contains(meshName, "_") {
			prefix = meshName[0:strings.Index(meshName, "_")]
		}
		mesh := &gltf.Mesh{Name: meshName}
		meshIndex := doc.MeshAdd(mesh)
		node := &gltf.Node{
			Name: meshName,
			Mesh: meshIndex,
			Skin: skinIndex,
		}

		if e.isDecoded {
			//node.Rotation = [4]float32{0.0, 0.0, -0.7, 0.7}

			//node.Rotation = [4]float32{0.0, 0.0, 1.0, 0.0}
			//node.Rotation = [4]float32{-0.5, 0.5, -0.5, 0.5}
		}
		doc.NodeAdd(node)
	}

	fmt.Println("prefix", prefix)

	context := ""

	tmpCache := make(map[string]bool)
	fmt.Println(len(e.triangles), "faces")
	// ******** PRIM GENERATION *****
	for _, o := range e.triangles {
		matName := o.MaterialName
		/*if strings.HasPrefix(matName, e.name+"_") {
			matName = fmt.Sprintf("c_%s_s02_m01", e.name)
		}*/

		var suffix int
		context = matName
		parts := strings.Split(matName, "_")
		if len(parts) == 3 {
			prefix = parts[0]
			context = parts[1]
			suffix, err = strconv.Atoi(parts[2])
			if err == nil {
				suffix = 1

				matName = fmt.Sprintf("%s_%s_%02d", prefix, context, suffix)
			}
		}

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
				// x-90 y-270
				//v.Position = helper.ApplyQuaternion(v.Position, [4]float32{0.5, -0.5, 0.5, -0.5})
				// x90
				v.Position = helper.ApplyQuaternion(v.Position, [4]float32{0.7071068, 0, 0, 0.7071068})
				prim.Positions = append(prim.Positions, v.Position)
				prim.Normals = append(prim.Normals, [3]float32{v.Normal[0], v.Normal[1], v.Normal[2]})
				prim.Uvs = append(prim.Uvs, [2]float32{v.Uv[0], v.Uv[1]})
				prim.UniqueIndices[o.Index[i]] = uint16(len(prim.Positions) - 1)
				index = uint16(len(prim.Positions) - 1)
				prim.Joints = append(prim.Joints, v.Joint)
				prim.Weights = append(prim.Weights, v.Weight)
			}
			prim.Indices = append(prim.Indices, index)
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
