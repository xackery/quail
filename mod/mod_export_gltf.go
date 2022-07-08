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
)

// ExportGLTF exports a provided mod file to gltf format
func (e *MOD) ExportGLTF(w io.Writer) error {
	var err error
	doc := gltf.NewDocument()

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

		if buf.Len() == 0 {
			return fmt.Errorf("%s not found", textureDiffuseName)
		}

		if filepath.Ext(textureDiffuseName) == ".dds" {
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
		}

		imageIdx, err := modeler.WriteImage(doc, textureDiffuseName, "image/png", buf)
		if err != nil {
			return fmt.Errorf("writeImage to gtlf: %w", err)
		}
		doc.Textures = append(doc.Textures, &gltf.Texture{Source: gltf.Index(imageIdx)})

		doc.Materials = append(doc.Materials, &gltf.Material{
			Name: mat.Name,
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
				BaseColorTexture: &gltf.TextureInfo{
					Index: uint32(len(doc.Textures) - 1),
				},
				MetallicFactor: gltf.Float(0),
			},
		})
	}

	mesh := &gltf.Mesh{
		Name: e.name,
	}

	prim := &gltf.Primitive{
		Mode: gltf.PrimitiveTriangles,
	}
	mesh.Primitives = append(mesh.Primitives, prim)

	positions := [][3]float32{}
	normals := [][3]float32{}
	uvs := [][2]float32{}
	indices := []uint16{}

	for _, vert := range e.vertices {
		positions = append(positions, [3]float32{vert.Position.X, vert.Position.Y, vert.Position.Z})
		normals = append(normals, [3]float32{vert.Normal.X, vert.Normal.Y, vert.Normal.Z})
		uvs = append(uvs, [2]float32{vert.Uv.X, vert.Uv.Y})
	}
	for _, o := range e.triangles {
		indices = append(indices, uint16(o.Index.X))
		indices = append(indices, uint16(o.Index.Y))
		indices = append(indices, uint16(o.Index.Z))
	}

	prim.Attributes, err = modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
		Position:       positions,
		Normal:         normals,
		TextureCoord_0: uvs,
	})
	if err != nil {
		return fmt.Errorf("writeAttributes: %w", err)
	}
	prim.Indices = gltf.Index(modeler.WriteIndices(doc, indices))
	doc.Meshes = append(doc.Meshes, mesh)

	for _, buff := range doc.Buffers {
		buff.EmbeddedResource()
	}

	enc := gltf.NewEncoder(w)
	enc.AsBinary = false
	err = enc.Encode(doc)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
