package mod

import (
	"fmt"
	"os"
	"strings"

	"github.com/qmuntal/gltf"

	qgltf "github.com/xackery/quail/gltf"
)

// GLTFEncode exports a provided mod file to gltf format
func (e *MOD) GLTFEncode(doc *qgltf.GLTF) error {
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
			if e.archive != nil {
				diffuseData, err = e.archive.File(textureDiffuseName)
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
			if e.archive != nil {
				normalData, err = e.archive.File(textureNormalName)
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
	var skinIndex *uint32
	joints := []uint32{}
	for _, b := range e.bones {
		node := &gltf.Node{
			Name:        b.Name,
			Children:    []uint32{}, // TODO: traverse and get children
			Translation: b.Pivot,
			Rotation:    b.Rotation,
			Scale:       b.Scale,
		}

		nodeIndex := doc.NodeAdd(node)
		joints = append(joints, nodeIndex)
	}

	if len(e.bones) > 0 {
		tmp := doc.SkinAdd(&gltf.Skin{
			Name:   "ROOT",
			Joints: joints,
		})
		skinIndex = &tmp
	}

	// ******** PRIM GENERATION *****
	for _, o := range e.triangles {
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
				// TODO: fiddle fix
				// x90
				//v.Position = helper.ApplyQuaternion(v.Position, [4]float32{0.7071068, 0, 0, 0.7071068})

				prim.Positions = append(prim.Positions, v.Position)
				prim.Normals = append(prim.Normals, v.Normal)
				prim.Uvs = append(prim.Uvs, v.Uv)
				prim.Joints = append(prim.Joints, v.Joint)
				prim.Weights = append(prim.Weights, v.Weight)
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
		Skin: skinIndex,
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
