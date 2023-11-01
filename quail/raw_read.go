package quail

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/raw"
)

// Read takes a raw type and converts it to a quail type
func (q *Quail) RawRead(in raw.Reader) error {
	if q == nil {
		return fmt.Errorf("quail is nil")
	}
	switch val := in.(type) {
	case *raw.Lay:
		return q.layRead(val)
	case *raw.Ani:
		return q.aniRead(val)
	case *raw.Wld:
		return q.wldRead(val)
	case *raw.Dds:
		return q.ddsRead(val)
	case *raw.Bmp:
		return q.bmpRead(val)
	case *raw.Png:
		return q.pngRead(val)
	case *raw.Mod:
		return q.modRead(val)
	default:
		return fmt.Errorf("unknown type %T", val)
	}
}

func RawRead(in raw.Reader, q *Quail) error {
	if q == nil {
		return fmt.Errorf("quail is nil")
	}
	return q.RawRead(in)
}

func (q *Quail) aniRead(in *raw.Ani) error {
	if q.Header == nil {
		q.Header = &common.Header{}
	}
	q.Header.Version = int(in.Version)
	q.Header.Name = "animation"

	return nil
}

func (q *Quail) layRead(in *raw.Lay) error {
	if q.Header == nil {
		q.Header = &common.Header{}
	}
	q.Header.Version = int(in.Version)
	q.Header.Name = "layer"
	for _, model := range q.Models {
		for _, entry := range in.Entries {
			lay := &common.Layer{
				Material: entry.Material,
				Diffuse:  entry.Diffuse,
				Normal:   entry.Normal,
			}

			model.Layers = append(model.Layers, lay)
		}
	}
	return nil
}

func (q *Quail) wldRead(wld *raw.Wld) error {
	if len(wld.Fragments) == 0 {
		return fmt.Errorf("no fragments found")
	}
	for i, frag := range wld.Fragments {
		model, err := q.wldConvertMesh(wld, frag)
		if err != nil {
			return fmt.Errorf("fragment %d: %w", i, err)
		}
		if model != nil {
			q.Models = append(q.Models, model)
		}
	}
	return nil
}

func (q *Quail) wldConvertMesh(world *raw.Wld, frag raw.FragmentReader) (*common.Model, error) {
	if frag.FragCode() != 0x36 {
		return nil, nil
	}

	d, ok := frag.(*raw.WldFragMesh)
	if !ok {
		return nil, fmt.Errorf("assertion failed, wanted *Mesh, got %T", frag)
	}

	name := raw.Name(d.NameRef)
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

	matList, ok := matListFrag.(*raw.WldFragMaterialList)
	if !ok {
		return nil, fmt.Errorf("assertion failed, wanted *MaterialList, got %T", matListFrag)
	}

	for _, matRef := range matList.MaterialRefs {
		dstMaterial := &common.Material{}

		materialFrag, ok := world.Fragments[int(matRef)]
		if !ok {
			return nil, fmt.Errorf("mat ref %d not found", matRef)
		}

		srcMaterial, ok := materialFrag.(*raw.WldFragMaterial)
		if !ok {
			return nil, fmt.Errorf("assertion failed, wanted *TextureRef, got %T", materialFrag)
		}

		srcTexture := &raw.WldFragTexture{}

		if srcMaterial.TextureRef > 0 {
			textureRefFrag, ok := world.Fragments[int(srcMaterial.TextureRef)]
			if !ok {
				return nil, fmt.Errorf("textureref ref %d not found", srcMaterial.TextureRef)
			}

			textureRef, ok := textureRefFrag.(*raw.WldFragTextureRef)
			if !ok {
				return nil, fmt.Errorf("assertion failed, wanted *TextureRef, got %T", textureRefFrag)
			}

			textureFrag, ok := world.Fragments[int(textureRef.TextureRef)]
			if !ok {
				return nil, fmt.Errorf("texture ref %d not found", textureRef.TextureRef)
			}

			srcTexture, ok = textureFrag.(*raw.WldFragTexture)
			if !ok {
				return nil, fmt.Errorf("assertion failed, wanted *Texture, got %T", textureFrag)
			}

			for _, tRef := range srcTexture.TextureRefs {
				textureRefFrag, ok = world.Fragments[int(tRef)]
				if !ok {
					return nil, fmt.Errorf("tref ref %d not found", tRef)
				}

				textureList, ok := textureRefFrag.(*raw.WldFragTextureList)
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
		model.Vertices[i].Tint = common.RGBA{R: color.R, G: color.G, B: color.B, A: color.A}
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

func (q *Quail) ddsRead(in *raw.Dds) error {
	if q.Textures == nil {
		q.Textures = make(map[string][]byte)
	}
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write dds: %w", err)
	}
	q.Textures[in.FileName()] = buf.Bytes()
	return nil
}

func (q *Quail) bmpRead(in *raw.Bmp) error {
	if q.Textures == nil {
		q.Textures = make(map[string][]byte)
	}
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write bmp: %w", err)
	}
	q.Textures[in.FileName()] = buf.Bytes()
	return nil
}

func (q *Quail) pngRead(in *raw.Png) error {
	if q.Textures == nil {
		q.Textures = make(map[string][]byte)
	}
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write png: %w", err)
	}
	q.Textures[in.FileName()] = buf.Bytes()
	return nil
}

func (q *Quail) modRead(in *raw.Mod) error {
	// FIXME: this is a stub
	return nil
}
