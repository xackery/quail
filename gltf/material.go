package gltf

import (
	"bytes"
	"fmt"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/malashin/dds"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/xackery/quail/common"
)

func (e *GLTF) Material(name string) *uint32 {
	return e.materials[name]
}

func (e *GLTF) MaterialAdd(req *common.Material, diffuseData []byte, normalData []byte) (*uint32, error) {

	materialDataIndex := uint32(0)

	index := &materialDataIndex
	if e.doc == nil {
		return index, fmt.Errorf("gltf is not initialized")
	}

	index, ok := e.materials[req.Name]
	if ok {
		return index, nil
	}

	if req.Name == "" || strings.HasPrefix(req.Name, "empty_") {
		e.doc.Materials = append(e.doc.Materials, &gltf.Material{
			Name: req.Name,
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
				BaseColorFactor: &[4]float32{1.0, 1.0, 1.0, 1},
				MetallicFactor:  gltf.Float(0),
			},
		})

		index = gltf.Index(uint32(len(e.doc.Materials) - 1))
		e.materials[req.Name] = index
		return index, nil
	}

	textureDiffuseName := ""
	textureNormalName := ""
	for _, p := range req.Properties {
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
		e.doc.Materials = append(e.doc.Materials, &gltf.Material{
			Name: req.Name,
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
				BaseColorFactor: &[4]float32{1.0, 1.0, 1.0, 1},
				MetallicFactor:  gltf.Float(0),
			},
		})
		index = gltf.Index(uint32(len(e.doc.Materials) - 1))
		e.materials[req.Name] = index
		return index, nil
	}

	diffuseBuf := bytes.NewBuffer(diffuseData)
	normalBuf := bytes.NewBuffer(normalData)

	if diffuseBuf.Len() == 0 {
		return index, fmt.Errorf("texture '%s' not found", textureDiffuseName)
	}

	pngData, err := toPNG(diffuseBuf, textureDiffuseName)
	if err != nil {
		return index, fmt.Errorf("gltfToPNG diffuse: %w", err)
	}
	textureDiffuseName = strings.ReplaceAll(textureDiffuseName, ".dds", ".png")
	diffuseBuf = bytes.NewBuffer(pngData)

	if normalBuf.Len() > 0 {
		pngData, err = toPNG(normalBuf, textureNormalName)
		if err != nil {
			return index, fmt.Errorf("gltfToPNG normal: %w", err)
		}
		normalBuf = bytes.NewBuffer(pngData)

		textureNormalName = strings.ReplaceAll(textureNormalName, ".dds", ".png")
	}

	meshName := strings.TrimSuffix(textureDiffuseName, ".png")
	imageIdx, err := modeler.WriteImage(e.doc, textureDiffuseName, "image/png", diffuseBuf)
	if err != nil {
		return index, fmt.Errorf("writeImage to gtlf: %w", err)
	}
	e.doc.Textures = append(e.doc.Textures, &gltf.Texture{Source: gltf.Index(imageIdx)})
	diffuseTexture := &gltf.TextureInfo{
		Index: uint32(len(e.doc.Textures) - 1),
	}

	var normalTexture *gltf.NormalTexture

	if normalBuf.Len() > 0 {
		imageIdx, err = modeler.WriteImage(e.doc, textureNormalName, "image/png", normalBuf)
		if err != nil {
			return index, fmt.Errorf("writeImage to gtlf: %w", err)
		}
		e.doc.Textures = append(e.doc.Textures, &gltf.Texture{Source: gltf.Index(imageIdx)})
		normalTexture = &gltf.NormalTexture{
			Index: gltf.Index(uint32(len(e.doc.Textures) - 1)),
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

	e.doc.Materials = append(e.doc.Materials, newMaterial)

	/*doc.Materials = append(doc.Materials, &gltf.Material{
		Name: modelName,
		PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
			BaseColorFactor: &[4]float32{1.0, 1.0, 1.0, 1},
			MetallicFactor:  gltf.Float(0),
		},
	})*/

	index = gltf.Index(uint32(len(e.doc.Materials) - 1))
	e.materials[req.Name] = index
	return index, nil
}

func toPNG(buf *bytes.Buffer, name string) ([]byte, error) {
	switch filepath.Ext(name) {
	case ".dds":
		img, err := dds.Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("dds.Decode %s: %w", name, err)
		}

		buf = bytes.NewBuffer(nil)
		err = png.Encode(buf, img)
		if err != nil {
			return nil, fmt.Errorf("png.Encode %s: %w", name, err)
		}
		return buf.Bytes(), nil
	case ".png":
		return buf.Bytes(), nil
	case "":
	default:
	}
	return nil, fmt.Errorf("unsupported extension '%s'", filepath.Ext(name))
}
