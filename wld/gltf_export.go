package wld

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"

	"github.com/malashin/dds"
	"github.com/xackery/quail/common"
)

//Ref: https://github.com/LanternEQ/LanternExtractor/blob/development/0.2.0/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs

// GLTFExport exports a provided mod file to gltf format
func (e *WLD) GLTFExport(w io.Writer) error {
	var err error
	doc := gltf.NewDocument()
	e.gltfMaterialBuffer = make(map[string]*uint32)
	e.gltfBoneBuffer = make(map[int]uint32)

	modelName := strings.TrimSuffix(e.name, ".ter")

	mesh := &gltf.Mesh{
		Name: modelName,
	}

	type primCache struct {
		materialIndex *uint32
		positions     [][3]float32
		normals       [][3]float32
		uvs           [][2]float32
		//joints        [][4]uint16
		//weights       [][4]uint16
		indices       []uint16
		uniqueIndices map[uint32]uint16
	}
	prims := make(map[*uint32]*primCache)

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
		matIndex, err := e.gltfAddCacheMaterial(doc, o.MaterialName)
		if err != nil {
			return fmt.Errorf("addMaterial: %w", err)
		}

		prim, ok := prims[matIndex]
		if !ok {
			prim = &primCache{
				materialIndex: matIndex,
				uniqueIndices: make(map[uint32]uint16),
			}
			prims[matIndex] = prim
		}

		for i := 0; i < 3; i++ {
			index, ok := prim.uniqueIndices[o.Index[i]]
			if !ok {
				v := e.vertices[int(o.Index[i])]
				prim.positions = append(prim.positions, [3]float32{v.Position.X, v.Position.Y, v.Position.Z})
				prim.normals = append(prim.normals, [3]float32{v.Normal.X, v.Normal.Y, v.Normal.Z})
				prim.uvs = append(prim.uvs, [2]float32{v.Uv.X, v.Uv.Y})
				prim.uniqueIndices[o.Index[i]] = uint16(len(prim.positions) - 1)
				index = uint16(len(prim.positions) - 1)
			}
			prim.indices = append(prim.indices, index)
		}

		/*primitive.Attributes, err = modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
			Position:       prim.positions,
			Normal:         prim.normals,
			TextureCoord_0: prim.uvs,
		})*/

	}

	for _, prim := range prims {
		primitive := &gltf.Primitive{
			Mode:     gltf.PrimitiveTriangles,
			Material: prim.materialIndex,
		}

		/*for _, pos := range prim.positions {
			x, y, z := pos[0], pos[1], pos[2]
			for _, b := range e.bones {
				if b.pivot.X != x {
					continue
				}
				if b.pivot.Y != y {
					continue
				}
				if b.pivot.Z != z {
					continue
				}
				prim.joints = append(prim.joints, [4]uint16{uint16(b.pivot.X), uint16(b.pivot.Y), uint16(b.pivot.Z)})
				prim.weights = append(prim.weights, [4]uint16{1, 1, 1, 1})
			}
		}*/

		primitive.Attributes = map[string]uint32{
			gltf.POSITION:   modeler.WritePosition(doc, prim.positions),
			gltf.NORMAL:     modeler.WriteNormal(doc, prim.normals),
			gltf.TEXCOORD_0: modeler.WriteTextureCoord(doc, prim.uvs),
			//gltf.JOINTS_0:   modeler.WriteJoints(doc, prim.joints),
			//gltf.WEIGHTS_0:  modeler.WriteWeights(doc, prim.weights),
		}

		primitive.Indices = gltf.Index(modeler.WriteIndices(doc, prim.indices))
		mesh.Primitives = append(mesh.Primitives, primitive)
	}

	//fmt.Println("last indices:", *primitive.Indices, "total buffers:", len(doc.BufferViews))
	doc.Meshes = append(doc.Meshes, mesh)
	doc.Nodes = append(doc.Nodes, &gltf.Node{
		Name: modelName,
		Mesh: gltf.Index(uint32(len(doc.Meshes) - 1)),
		//Skin: skinIndex,
	})

	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, uint32(len(doc.Nodes)-1))
	for _, buff := range doc.Buffers {
		buff.EmbeddedResource()
	}

	enc := gltf.NewEncoder(w)
	enc.AsBinary = false
	enc.SetJSONIndent("", "\t")
	err = enc.Encode(doc)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}

func (e *WLD) gltfAddCacheMaterial(doc *gltf.Document, name string) (*uint32, error) {
	materialIndex := uint32(0)
	material := &materialIndex

	material, ok := e.gltfMaterialBuffer[name]
	if ok {
		return material, nil
	}

	var mat *common.Material
	for _, o := range e.materials {
		if o.Name == name {
			mat = o
			break
		}
	}
	if mat == nil {
		if name == "" {
			doc.Materials = append(doc.Materials, &gltf.Material{
				Name: name,
				PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
					BaseColorFactor: &[4]float32{1.0, 1.0, 1.0, 1},
					MetallicFactor:  gltf.Float(0),
				},
			})

			material = gltf.Index(uint32(len(doc.Materials) - 1))
			e.gltfMaterialBuffer[name] = material
			return material, nil
		}
		return material, fmt.Errorf("material '%s' not found", name)
	}

	textureDiffuseName := ""
	textureNormalName := ""
	for _, p := range mat.Properties {
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
	if len(textureDiffuseName) == 0 {
		//return material, fmt.Errorf("material '%s' has no texturediffuse value", name)
		doc.Materials = append(doc.Materials, &gltf.Material{
			Name: name,
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
				BaseColorFactor: &[4]float32{1.0, 1.0, 1.0, 1},
				MetallicFactor:  gltf.Float(0),
			},
		})
		material = gltf.Index(uint32(len(doc.Materials) - 1))
		e.gltfMaterialBuffer[name] = material
		return material, nil
	}

	diffuseBuf := &bytes.Buffer{}
	normalBuf := &bytes.Buffer{}
	for _, fe := range e.files {
		if fe.Name() == textureDiffuseName {
			diffuseBuf = bytes.NewBuffer(fe.Data())
		}
		if fe.Name() == textureNormalName {
			normalBuf = bytes.NewBuffer(fe.Data())
		}
		if diffuseBuf.Len() > 0 && normalBuf.Len() > 0 {
			break
		}
	}

	if diffuseBuf.Len() == 0 {
		return material, fmt.Errorf("texture '%s' not found", textureDiffuseName)
	}

	err := gltfToPNG(diffuseBuf, textureDiffuseName)
	if err != nil {
		return material, fmt.Errorf("gltfToPNG diffuse: %w", err)
	}
	textureDiffuseName = strings.ReplaceAll(textureDiffuseName, ".dds", ".png")

	if normalBuf.Len() > 0 {
		fmt.Println("normal", textureNormalName, "is", normalBuf.Len())
		err = gltfToPNG(normalBuf, textureNormalName)
		if err != nil {
			return material, fmt.Errorf("gltfToPNG normal: %w", err)
		}

		textureNormalName = strings.ReplaceAll(textureNormalName, ".dds", ".png")
		fmt.Println("normal", textureNormalName, "is", normalBuf.Len())
	}

	meshName := strings.TrimSuffix(textureDiffuseName, ".png")
	imageIdx, err := modeler.WriteImage(doc, textureDiffuseName, "image/png", diffuseBuf)
	if err != nil {
		return material, fmt.Errorf("writeImage to gtlf: %w", err)
	}
	doc.Textures = append(doc.Textures, &gltf.Texture{Source: gltf.Index(imageIdx)})
	diffuseTexture := &gltf.TextureInfo{
		Index: uint32(len(doc.Textures) - 1),
	}

	var normalTexture *gltf.NormalTexture
	if normalBuf.Len() > 0 {
		imageIdx, err = modeler.WriteImage(doc, textureNormalName, "image/png", normalBuf)
		if err != nil {
			return material, fmt.Errorf("writeImage to gtlf: %w", err)
		}
		doc.Textures = append(doc.Textures, &gltf.Texture{Source: gltf.Index(imageIdx)})
		normalTexture = &gltf.NormalTexture{
			Index: gltf.Index(uint32(len(doc.Textures) - 1)),
		}
	}

	newMaterial := &gltf.Material{
		Name: meshName,

		PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
			BaseColorTexture: diffuseTexture,
			MetallicFactor:   gltf.Float(0),
		},
	}
	if normalTexture != nil {
		newMaterial.NormalTexture = normalTexture
	}

	doc.Materials = append(doc.Materials, newMaterial)

	/*doc.Materials = append(doc.Materials, &gltf.Material{
		Name: modelName,
		PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
			BaseColorFactor: &[4]float32{1.0, 1.0, 1.0, 1},
			MetallicFactor:  gltf.Float(0),
		},
	})*/

	material = gltf.Index(uint32(len(doc.Materials) - 1))
	e.gltfMaterialBuffer[name] = material
	return material, nil
}

/*
func (e *WLD) gltfBoneChildren(doc *gltf.Document, children *[]uint32, boneIndex int) error {

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

func gltfToPNG(buf *bytes.Buffer, name string) error {
	switch filepath.Ext(name) {
	case ".dds":
		img, err := dds.Decode(buf)
		if err != nil {
			return fmt.Errorf("dds.Decode %s: %w", name, err)
		}

		buf = bytes.NewBuffer(nil)
		err = png.Encode(buf, img)
		if err != nil {
			return fmt.Errorf("png.Encode %s: %w", name, err)
		}
	case ".png":
	case "":
	default:
		return fmt.Errorf("unsupported extension: %s", filepath.Ext(name))
	}
	return nil
}
