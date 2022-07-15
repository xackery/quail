package mod

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/spate/glimage/dds"
	"github.com/xackery/quail/common"
)

// GLTFExport exports a provided mod file to gltf format
func (e *MOD) GLTFExport(w io.Writer) error {
	var err error
	doc := gltf.NewDocument()
	e.gltfMaterialBuffer = make(map[string]*uint32)

	modelName := strings.TrimSuffix(e.name, ".ter")

	mesh := &gltf.Mesh{
		Name: modelName,
	}

	type primCache struct {
		materialIndex *uint32
		positions     [][3]float32
		normals       [][3]float32
		uvs           [][2]float32
		indices       []uint16
		uniqueIndices map[uint32]uint16
	}
	prims := make(map[*uint32]*primCache)

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

		primitive.Attributes = map[string]uint32{
			gltf.POSITION:   modeler.WritePosition(doc, prim.positions),
			gltf.NORMAL:     modeler.WriteNormal(doc, prim.normals),
			gltf.TEXCOORD_0: modeler.WriteTextureCoord(doc, prim.uvs),
		}

		primitive.Indices = gltf.Index(modeler.WriteIndices(doc, prim.indices))
		mesh.Primitives = append(mesh.Primitives, primitive)
	}
	//fmt.Println("last indices:", *primitive.Indices, "total buffers:", len(doc.BufferViews))
	doc.Meshes = append(doc.Meshes, mesh)
	doc.Nodes = append(doc.Nodes, &gltf.Node{Name: modelName, Mesh: gltf.Index(uint32(len(doc.Meshes) - 1))})
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

func (e *MOD) gltfAddCacheMaterial(doc *gltf.Document, name string) (*uint32, error) {
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
	for _, p := range mat.Properties {
		if p.Category == 2 && strings.ToLower(p.Name) == "e_texturediffuse0" {
			textureDiffuseName = p.Value
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

	buf := &bytes.Buffer{}
	for _, fe := range e.files {
		if fe.Name() != textureDiffuseName {
			continue
		}
		buf = bytes.NewBuffer(fe.Data())
		break
	}

	if buf.Len() == 0 {
		return material, fmt.Errorf("texture '%s' not found", textureDiffuseName)
	}

	switch filepath.Ext(textureDiffuseName) {
	case ".dds":
		img, err := dds.Decode(buf)
		if err != nil {
			return material, fmt.Errorf("dds.Decode %s: %w", textureDiffuseName, err)
		}
		buf = bytes.NewBuffer(nil)
		err = png.Encode(buf, img)
		if err != nil {
			return material, fmt.Errorf("png.Encode %s: %w", textureDiffuseName, err)
		}
		textureDiffuseName = strings.ReplaceAll(textureDiffuseName, ".dds", ".png")
	case ".png":
	case "":
	default:
		return material, fmt.Errorf("material %s has a texture of %s which is unsupported", e.name, textureDiffuseName)
	}

	meshName := strings.TrimSuffix(textureDiffuseName, ".png")
	imageIdx, err := modeler.WriteImage(doc, textureDiffuseName, "image/png", buf)
	if err != nil {
		return material, fmt.Errorf("writeImage to gtlf: %w", err)
	}
	doc.Textures = append(doc.Textures, &gltf.Texture{Source: gltf.Index(imageIdx)})

	doc.Materials = append(doc.Materials, &gltf.Material{
		Name: meshName,
		PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
			BaseColorTexture: &gltf.TextureInfo{
				Index: uint32(len(doc.Textures) - 1),
			},
			MetallicFactor: gltf.Float(0),
		},
	})
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
