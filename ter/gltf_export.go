package ter

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

// GLTFExport exports a provided ter file to gltf format
func (e *TER) GLTFExport(w io.Writer) error {
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
		uniqueIndeces map[uint16]bool
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
				uniqueIndeces: make(map[uint16]bool),
			}
			prims[matIndex] = prim
		}

		prim.indices = append(prim.indices, uint16(o.Index[0]))
		prim.indices = append(prim.indices, uint16(o.Index[1]))
		prim.indices = append(prim.indices, uint16(o.Index[2]))
		prim.uniqueIndeces[uint16(o.Index[0])] = true
		prim.uniqueIndeces[uint16(o.Index[1])] = true
		prim.uniqueIndeces[uint16(o.Index[2])] = true
	}

	for _, prim := range prims {
		for index := range prim.uniqueIndeces {
			o := e.vertices[int(index)]
			prim.positions = append(prim.positions, [3]float32{o.Position.X, o.Position.Y, o.Position.Z})
			prim.normals = append(prim.normals, [3]float32{o.Normal.X, o.Normal.Y, o.Normal.Z})
			prim.uvs = append(prim.uvs, [2]float32{o.Uv.X, o.Uv.Y})
		}
		primitive := &gltf.Primitive{
			Mode:     gltf.PrimitiveTriangles,
			Material: prim.materialIndex,
		}
		primitive.Attributes, err = modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
			Position:       prim.positions,
			Normal:         prim.normals,
			TextureCoord_0: prim.uvs,
		})
		if err != nil {
			return fmt.Errorf("writeAttributes: %w", err)
		}
		primitive.Indices = gltf.Index(modeler.WriteIndices(doc, prim.indices))
		mesh.Primitives = append(mesh.Primitives, primitive)
	}

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

func (e *TER) gltfAddCacheMaterial(doc *gltf.Document, name string) (*uint32, error) {
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
				BaseColorFactor: &[4]float32{0, 0, 0, 1},
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
			BaseColorFactor: &[4]float32{0, 0, 0, 1},
			MetallicFactor:  gltf.Float(0),
		},
	})*/

	material = gltf.Index(uint32(len(doc.Materials) - 1))
	e.gltfMaterialBuffer[name] = material
	return material, nil
}
