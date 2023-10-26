package quail

import (
	"fmt"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/wld"
)

func (e *Quail) WldUnmarshal(world *common.Wld) error {
	if len(world.Fragments) == 0 {
		return fmt.Errorf("no fragments found")
	}
	e.Header = world.Header
	for i, frag := range world.Fragments {
		model, err := convertMesh(world, frag)
		if err != nil {
			return fmt.Errorf("fragment id %d 0x%x (%s): %w", i, frag.FragCode(), common.FragName(frag.FragCode()), err)
		}
		if model == nil {
			continue
		}
		e.Models = append(e.Models, model)
	}
	return nil
}

func convertMesh(world *common.Wld, frag common.FragmentReader) (*common.Model, error) {
	if frag.FragCode() != 0x36 {
		return nil, nil
	}

	d, ok := frag.(*wld.Mesh)
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

	scale := float32(1 / float32(int(1)<<int(d.RawScale))) // This allows vertex coordinates to be stored as integral values instead of floating-point values, without losing precision based on mesh size. Vertex values are multiplied by (1 shl `scale`) and stored in the vertex entries. FPSCALE is the internal name.

	matListFrag, ok := world.Fragments[int(d.MaterialListRef)]
	if !ok {
		return nil, fmt.Errorf("material list ref %d not found", d.MaterialListRef)
	}

	matList, ok := matListFrag.(*wld.MaterialList)
	if !ok {
		return nil, fmt.Errorf("assertion failed, wanted *MaterialList, got %T", matListFrag)
	}

	for _, matRef := range matList.MaterialRefs {
		dstMaterial := &common.Material{}

		materialFrag, ok := world.Fragments[int(matRef)]
		if !ok {
			return nil, fmt.Errorf("mat ref %d not found", matRef)
		}

		srcMaterial, ok := materialFrag.(*wld.Material)
		if !ok {
			return nil, fmt.Errorf("assertion failed, wanted *TextureRef, got %T", materialFrag)
		}

		srcTexture := &wld.Texture{}

		if srcMaterial.TextureRef > 0 {
			textureRefFrag, ok := world.Fragments[int(srcMaterial.TextureRef)]
			if !ok {
				return nil, fmt.Errorf("textureref ref %d not found", srcMaterial.TextureRef)
			}

			textureRef, ok := textureRefFrag.(*wld.TextureRef)
			if !ok {
				return nil, fmt.Errorf("assertion failed, wanted *TextureRef, got %T", textureRefFrag)
			}

			textureFrag, ok := world.Fragments[int(textureRef.TextureRef)]
			if !ok {
				return nil, fmt.Errorf("texture ref %d not found", textureRef.TextureRef)
			}

			srcTexture, ok = textureFrag.(*wld.Texture)
			if !ok {
				return nil, fmt.Errorf("assertion failed, wanted *Texture, got %T", textureFrag)
			}

			for _, tRef := range srcTexture.TextureRefs {
				textureRefFrag, ok = world.Fragments[int(tRef)]
				if !ok {
					return nil, fmt.Errorf("tref ref %d not found", tRef)
				}

				textureList, ok := textureRefFrag.(*wld.TextureList)
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

	for _, vertex := range d.Vertices {
		model.Vertices = append(model.Vertices, common.Vertex{
			Position: common.Vector3{
				X: float32(d.Center.X) + (float32(vertex[0]) * scale),
				Y: float32(d.Center.Y) + (float32(vertex[1]) * scale),
				Z: float32(d.Center.Z) + (float32(vertex[2]) * scale),
			},
		})
	}

	for i, normal := range d.Normals {
		if len(model.Triangles) <= i {
			for len(model.Triangles) <= i {
				model.Triangles = append(model.Triangles, common.Triangle{})
			}
		}
		model.Triangles[i].Index = common.UIndex3{
			X: uint32(normal[0]),
			Y: uint32(normal[1]),
			Z: uint32(normal[2]),
		}
	}

	for i, triangle := range d.Triangles {
		if len(model.Triangles) <= i {
			for len(model.Triangles) <= i {
				model.Triangles = append(model.Triangles, common.Triangle{})
			}
		}
		model.Triangles[i].Index = common.UIndex3{
			X: uint32(triangle.Index[0]),
			Y: uint32(triangle.Index[1]),
			Z: uint32(triangle.Index[2]),
		}
		model.Triangles[i].Flag = uint32(triangle.Flags)
	}

	for i, color := range d.Colors {
		if len(model.Vertices) <= i {
			for len(model.Vertices) <= i {
				model.Vertices = append(model.Vertices, common.Vertex{})
			}
		}
		model.Vertices[i].Tint = color
	}

	for i, uv := range d.UVs {
		if len(model.Vertices) <= i {
			for len(model.Vertices) <= i {
				model.Vertices = append(model.Vertices, common.Vertex{})
			}
		}

		model.Vertices[i].Uv = common.Vector2{
			X: float32(uv[0] / 256),
			Y: float32(uv[1] / 256),
		}
	}

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
