package wld

import (
	"fmt"

	"github.com/xackery/quail/common"
)

func Convert(world *common.Wld) ([]*common.Model, error) {
	models := []*common.Model{}
	for i, frag := range world.Fragments {
		model, err := convertMesh(world, frag)
		if err != nil {
			return nil, fmt.Errorf("fragment id %d 0x%x (%s): %w", i, frag.FragCode(), common.FragName(frag.FragCode()), err)
		}
		if model == nil {
			continue
		}
		models = append(models, model)
	}

	return models, nil
}

func convertMesh(world *common.Wld, frag common.FragmentReader) (*common.Model, error) {
	if frag.FragCode() != 0x36 {
		return nil, nil
	}

	d, ok := frag.(*Mesh)
	if !ok {
		return nil, fmt.Errorf("assertion failed, wanted *Mesh, got %T", frag)
	}
	name := world.Name(d.NameRef)
	if name == "" {
		name = "unknown"
	}
	model := common.NewModel(name)

	model.FileType = "mod"
	if d.Flags == 0x00018003 {
		model.FileType = "ter"
	}
	if d.Flags == 0x00014003 {
		model.FileType = "mod"
	}

	matListFrag, ok := world.Fragments[int(d.MaterialListRef)]
	if !ok {
		return nil, fmt.Errorf("material list ref %d not found", d.MaterialListRef)
	}

	matList, ok := matListFrag.(*MaterialList)
	if !ok {
		return nil, fmt.Errorf("assertion failed, wanted *MaterialList, got %T", matListFrag)
	}

	for _, matRef := range matList.MaterialRefs {
		dstMaterial := &common.Material{}

		materialFrag, ok := world.Fragments[int(matRef)]
		if !ok {
			return nil, fmt.Errorf("mat ref %d not found", matRef)
		}

		srcMaterial, ok := materialFrag.(*Material)
		if !ok {
			return nil, fmt.Errorf("assertion failed, wanted *TextureRef, got %T", materialFrag)
		}

		srcTexture := &Texture{}

		if srcMaterial.TextureRef > 0 {

			textureRefFrag, ok := world.Fragments[int(srcMaterial.TextureRef)]
			if !ok {
				return nil, fmt.Errorf("textureref ref %d not found", srcMaterial.TextureRef)
			}

			textureRef, ok := textureRefFrag.(*TextureRef)
			if !ok {
				return nil, fmt.Errorf("assertion failed, wanted *TextureRef, got %T", textureRefFrag)
			}

			textureFrag, ok := world.Fragments[int(textureRef.TextureRef)]
			if !ok {
				return nil, fmt.Errorf("texture ref %d not found", textureRef.TextureRef)
			}

			srcTexture, ok = textureFrag.(*Texture)
			if !ok {
				return nil, fmt.Errorf("assertion failed, wanted *Texture, got %T", textureFrag)
			}

			for _, tRef := range srcTexture.TextureRefs {
				textureRefFrag, ok = world.Fragments[int(tRef)]
				if !ok {
					return nil, fmt.Errorf("tref ref %d not found", tRef)
				}

				textureList, ok := textureRefFrag.(*TextureList)
				if !ok {
					return nil, fmt.Errorf("tref assertion failed, wanted *TextureRef, got %T", textureRefFrag)
				}

				for i, textureName := range textureList.TextureNames {
					property := &common.MaterialProperty{}
					property.Category = 2
					property.Name = fmt.Sprintf("e_TextureDiffuse%d", i)
					property.Value = textureName
					if dstMaterial.Name == "" {
						dstMaterial.Name = textureName
					}
					dstMaterial.Properties = append(dstMaterial.Properties, property)
				}
			}
		}

		dstMaterial.Flag = srcMaterial.Flags
		dstMaterial.ShaderName = "Opaque_MaxC1.fx"

		model.Materials = append(model.Materials, dstMaterial)
	}

	model.Vertices = d.Vertices
	model.Triangles = d.Triangles
	triIndex := 0
	for i, triangleMat := range d.TriangleMaterials {
		if int(triangleMat.MaterialID) >= len(model.Materials) {
			return nil, fmt.Errorf("triangle material %d is out of bounds %d", i, triangleMat.MaterialID)
		}
		name := model.Materials[int(triangleMat.MaterialID)].Name
		for j := 0; j < int(triangleMat.Count); j++ {
			model.Triangles[triIndex].MaterialName = name
			triIndex++
		}
	}

	return model, nil
}
