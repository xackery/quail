package quail

import (
	"fmt"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/raw"
)

func (q *Quail) wldRead(wld *raw.Wld) error {
	var err error
	if len(wld.Fragments) == 0 {
		return fmt.Errorf("no fragments found")
	}
	maxFragments := len(wld.Fragments)
	for i := 1; i < maxFragments; i++ {
		frag, ok := wld.Fragments[i]
		if !ok {
			return fmt.Errorf("fragment %d not found", i)
		}

		switch frag.FragCode() {
		case 0x10: // SkeletonTrack
			err = q.wldConvertSkeletonTrack(wld, frag)
			if err != nil {
				return fmt.Errorf("fragment %d: %w", i, err)
			}
		case 0x14: // Model
			d, ok := frag.(*raw.WldFragActorDef)
			if !ok {
				return fmt.Errorf("assertion failed, wanted *Mesh, got %T", frag)
			}

			name := raw.Name(d.NameRef)
			if name == "" {
				name = "unknown"
			}
			fmt.Println("added model", name)
		case 0x36: // Mesh
			model, err := q.wldConvertMesh(wld, frag)
			if err != nil {
				return fmt.Errorf("fragment %d: %w", i, err)
			}
			if model != nil {
				q.Models = append(q.Models, model)
			}
			fmt.Println("added mesh", model.Header.Name)
		default:
		}
	}
	return nil
}

func (q *Quail) wldConvertMesh(world *raw.Wld, frag raw.FragmentReadWriter) (*common.Model, error) {
	if frag.FragCode() != 0x36 {
		return nil, nil
	}

	d, ok := frag.(*raw.WldFragDmSpriteDef2)
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

	matList, ok := matListFrag.(*raw.WldFragMaterialPalette)
	if !ok {
		return nil, fmt.Errorf("assertion failed, wanted *MaterialList, got %T", matListFrag)
	}

	for _, matRef := range matList.MaterialRefs {
		dstMaterial := &common.Material{}

		materialFrag, ok := world.Fragments[int(matRef)]
		if !ok {
			return nil, fmt.Errorf("mat ref %d not found", matRef)
		}

		srcMaterial, ok := materialFrag.(*raw.WldFragMaterialDef)
		if !ok {
			return nil, fmt.Errorf("assertion failed, wanted *TextureRef, got %T", materialFrag)
		}

		srcTexture := &raw.WldFragSimpleSpriteDef{}

		if srcMaterial.TextureRef > 0 {
			textureRefFrag, ok := world.Fragments[int(srcMaterial.TextureRef)]
			if !ok {
				return nil, fmt.Errorf("textureref ref %d not found", srcMaterial.TextureRef)
			}

			textureRef, ok := textureRefFrag.(*raw.WldFragSimpleSprite)
			if !ok {
				return nil, fmt.Errorf("assertion failed, wanted *TextureRef, got %T", textureRefFrag)
			}

			textureFrag, ok := world.Fragments[int(textureRef.TextureRef)]
			if !ok {
				return nil, fmt.Errorf("texture ref %d not found", textureRef.TextureRef)
			}

			srcTexture, ok = textureFrag.(*raw.WldFragSimpleSpriteDef)
			if !ok {
				return nil, fmt.Errorf("assertion failed, wanted *Texture, got %T", textureFrag)
			}

			for _, tRef := range srcTexture.TextureRefs {
				textureRefFrag, ok = world.Fragments[int(tRef)]
				if !ok {
					return nil, fmt.Errorf("tref ref %d not found", tRef)
				}

				textureList, ok := textureRefFrag.(*raw.WldFragBMInfo)
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

func (q *Quail) wldConvertSkeletonTrack(world *raw.Wld, frag raw.FragmentReadWriter) error {
	if frag.FragCode() != 0x10 {
		return nil
	}

	d, ok := frag.(*raw.WldFragHierarchialSpriteDef)
	if !ok {
		return fmt.Errorf("assertion failed, wanted *SkeletonTrack, got %T", frag)
	}

	name := raw.Name(d.NameRef)
	if name == "" {
		return fmt.Errorf("skeleton track name ref %d not found", d.NameRef)
	}

	baseModel := q.ModelByName(name[0:3] + "_DMSPRITEDEF")
	if baseModel == nil {
		return fmt.Errorf("skeleton track %s base model not found", name[0:3]+"_DMSPRITEDEF")
	}

	for i, bone := range d.Bones {
		newBone := common.Bone{}
		boneName := raw.Name(bone.NameRef)
		if boneName == "" {
			return fmt.Errorf("bone %d on skeleton %s, name ref %d not found", i, name, bone.NameRef)
		}
		newBone.Name = boneName

		var boneModel *common.Model

		if bone.MeshOrSpriteRef > 0 {
			modelName := raw.Name(int32(bone.MeshOrSpriteRef))
			if modelName == "" {
				return fmt.Errorf("bone %d on skeleton %s mesh ref %d not found", i, name, bone.MeshOrSpriteRef)
			}

			boneModel = q.ModelByName(modelName)
			if boneModel == nil {
				return fmt.Errorf("bone %d on skeleton %s mesh %s not found", i, name, modelName)
			}
		}

		if bone.Track > 0 {
			trackFragRef, ok := world.Fragments[int(bone.Track)]
			if !ok {
				return fmt.Errorf("bone %d on skeleton %s track ref %d not found", i, name, bone.Track)
			}

			trackRef, ok := trackFragRef.(*raw.WldFragTrack)
			if !ok {
				return fmt.Errorf("assertion failed, wanted *TrackRef, got %T", trackFragRef)
			}

			trackFrag, ok := world.Fragments[int(trackRef.Track)]
			if !ok {
				return fmt.Errorf("bone %d on skeleton %s track %d not found", i, name, trackRef.Track)
			}

			track, ok := trackFrag.(*raw.WldFragTrackDef)
			if !ok {
				return fmt.Errorf("assertion failed, wanted *Track, got %T", trackFrag)
			}

			for _, boneTransform := range track.BoneTransforms {
				newBone.Pivot.X = float32(boneTransform.TranslationX)
				newBone.Pivot.Y = float32(boneTransform.TranslationY)
				newBone.Pivot.Z = float32(boneTransform.TranslationZ)
				newBone.Rotation.X = float32(boneTransform.RotationX)
				newBone.Rotation.Y = float32(boneTransform.RotationY)
				newBone.Rotation.Z = float32(boneTransform.RotationZ)
			}
			newBone.Flags = bone.Flags
		}

		newBone.ChildrenCount = uint32(len(bone.SubBones))
		for _, child := range bone.SubBones {
			newBone.Children = append(newBone.Children, int(child))
		}

		if boneModel != nil {
			boneModel.Bones = append(boneModel.Bones, newBone)
		}
		baseModel.Bones = append(baseModel.Bones, newBone)

	}

	fmt.Println("added skeleton track", name, "bones", len(d.Bones))
	return nil
}

func (q *Quail) ModelByName(name string) *common.Model {
	for _, model := range q.Models {
		if model.Header.Name == name {
			return model
		}
	}
	return nil
}
