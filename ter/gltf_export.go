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
)

// GLTFExport exports a provided ter file to gltf format
func (e *TER) GLTFExport(w io.Writer) error {
	var err error
	doc := gltf.NewDocument()

	modelName := strings.TrimSuffix(e.name, ".ter")

	mesh := &gltf.Mesh{
		Name: modelName,
	}

	for _, mat := range e.materials {
		textureDiffuseName := ""
		for _, p := range mat.Properties {
			if p.Category == 2 && strings.ToLower(p.Name) == "e_texturediffuse0" {
				textureDiffuseName = p.Value
			}
		}

		buf := &bytes.Buffer{}
		if len(textureDiffuseName) > 0 {
			for _, fe := range e.files {
				if fe.Name() != textureDiffuseName {
					continue
				}
				buf = bytes.NewBuffer(fe.Data())
				break
			}
		}

		if buf.Len() == 0 && len(textureDiffuseName) > 0 {
			return fmt.Errorf("texture '%s' not found", textureDiffuseName)
		}

		switch filepath.Ext(textureDiffuseName) {
		case ".dds":
			img, err := dds.Decode(buf)
			if err != nil {
				return fmt.Errorf("dds.Decode %s: %w", textureDiffuseName, err)
			}
			buf = bytes.NewBuffer(nil)
			err = png.Encode(buf, img)
			if err != nil {
				return fmt.Errorf("png.Encode %s: %w", textureDiffuseName, err)
			}
			textureDiffuseName = strings.ReplaceAll(textureDiffuseName, ".dds", ".png")
		case ".png":
		case "":
		default:
			return fmt.Errorf("material %s has a texture of %s which is unsupported", e.name, textureDiffuseName)
		}

		if len(textureDiffuseName) > 0 {
			meshName := strings.TrimSuffix(textureDiffuseName, ".png")
			imageIdx, err := modeler.WriteImage(doc, textureDiffuseName, "image/png", buf)
			if err != nil {
				return fmt.Errorf("writeImage to gtlf: %w", err)
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
		} else {
			doc.Materials = append(doc.Materials, &gltf.Material{
				Name: modelName,
				PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
					BaseColorFactor: &[4]float32{0, 0, 0, 1},
					MetallicFactor:  gltf.Float(0),
				},
			})
		}

		prim := &gltf.Primitive{
			Mode:     gltf.PrimitiveTriangles,
			Material: gltf.Index(uint32(len(doc.Materials) - 1)),
		}
		mesh.Primitives = append(mesh.Primitives, prim)

		positions := [][3]float32{}
		normals := [][3]float32{}
		uvs := [][2]float32{}
		indices := []uint16{}

		for i, o := range e.triangles {
			if o.MaterialName != mat.Name {
				continue
			}
			positions = append(positions, [3]float32{e.vertices[i].Position.X, e.vertices[i].Position.Y, e.vertices[i].Position.Z})
			normals = append(normals, [3]float32{e.vertices[i].Normal.X, e.vertices[i].Normal.Y, e.vertices[i].Normal.Z})
			uvs = append(uvs, [2]float32{e.vertices[i].Uv.X, e.vertices[i].Uv.Y})

			indices = append(indices, uint16(o.Index[0]))
			indices = append(indices, uint16(o.Index[1]))
			indices = append(indices, uint16(o.Index[2]))
			if i == 0 {
				fmt.Println(positions, normals, uvs, indices)
			}
		}

		prim.Attributes, err = modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
			Position:       positions,
			Normal:         normals,
			TextureCoord_0: uvs,
		})
		if err != nil {
			return fmt.Errorf("writeAttributes: %w", err)
		}
		fmt.Println(modelName, "has", len(positions), "positions and", len(indices), "indices based on", len(e.triangles), "total triangles")
		prim.Indices = gltf.Index(modeler.WriteIndices(doc, indices))
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
