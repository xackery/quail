package vwld

import (
	"fmt"

	"github.com/xackery/quail/raw"
)

func (wld *VWld) Read(src *raw.Wld) error {
	wld.FileName = src.FileName()
	wld.Version = src.Version
	for i := 1; i < len(src.Fragments); i++ {
		fragment := src.Fragments[i]
		//log.Println("Fragment: ", raw.FragName(fragment.FragCode()), i)

		switch fragment.FragCode() {
		case raw.FragCodeBMInfo:
			fragData, ok := fragment.(*raw.WldFragBMInfo)
			if !ok {
				return fmt.Errorf("invalid bminfo fragment at offset %d", i)
			}
			name := raw.Name(fragData.NameRef)
			if len(name) == 0 {
				name = fmt.Sprintf("%d_BMINFO", i)
			}

			bitmap := &Bitmap{
				fragID:   uint32(i),
				Name:     name,
				Textures: fragData.TextureNames,
			}
			wld.Bitmaps = append(wld.Bitmaps, bitmap)
		case raw.FragCodeSimpleSpriteDef:
			fragData, ok := fragment.(*raw.WldFragSimpleSpriteDef)
			if !ok {
				return fmt.Errorf("invalid simplespritedef fragmentat offset %d", i)
			}

			name := raw.Name(fragData.NameRef)
			if len(name) == 0 {
				name = fmt.Sprintf("%d_SPRITE", i)
			}

			sprite := Sprite{
				fragID:       uint32(i),
				Name:         name,
				Flags:        fragData.Flags,
				CurrentFrame: fragData.CurrentFrame,
				Sleep:        fragData.Sleep,
			}
			for _, bitmapRef := range fragData.BitmapRefs {
				bitmap := wld.bitmapByFragID(bitmapRef)
				if bitmap == nil {
					return fmt.Errorf("simple sprite found without matching bminfo at offset %d", i)
				}

				sprite.Bitmaps = append(sprite.Bitmaps, bitmap.Name)
			}
			wld.Sprites = append(wld.Sprites, &sprite)
		case raw.FragCodeSimpleSprite:
			fragData, ok := fragment.(*raw.WldFragSimpleSprite)
			if !ok {
				return fmt.Errorf("invalid simplesprite fragment at offset %d", i)
			}

			sprite := wld.spriteByFragID(uint32(fragData.SpriteRef))
			if sprite == nil {
				return fmt.Errorf("simplesprite found without matching simplespritedef at offset %d, ref: %d", i, fragData.SpriteRef)
			}

			name := raw.Name(fragData.NameRef)
			if len(name) == 0 {
				name = sprite.Name + "_INST"
			}
			spriteInstance := SpriteInstance{
				fragID: uint32(i),
				Name:   name,
				Flags:  fragData.Flags,
				Sprite: sprite.Name,
			}
			wld.SpriteInstances = append(wld.SpriteInstances, &spriteInstance)
		case raw.FragCodeBlitSpriteDef:
			fragData, ok := fragment.(*raw.WldFragBlitSpriteDef)
			if !ok {
				return fmt.Errorf("invalid blitspritedef fragment at offset %d", i)
			}

			name := raw.Name(fragData.NameRef)
			if len(name) == 0 {
				name = fmt.Sprintf("%d_SPB", i)
			}

			particle := &Particle{
				fragID:     uint32(i),
				Name:       name,
				Flags:      fragData.Flags,
				SpriteName: raw.Name(int32(fragData.BlitSpriteRef)),
				Unknown:    fragData.Unknown,
			}
			wld.Particles = append(wld.Particles, particle)
		case raw.FragCodeParticleCloudDef:
			fragData, ok := fragment.(*raw.WldFragParticleCloudDef)
			if !ok {
				return fmt.Errorf("invalid particleclouddef fragment at offset %d", i)
			}

			particle := wld.particleByFragID(uint32(fragData.ParticleRef))
			if particle == nil {
				return fmt.Errorf("particleclouddef found without matching blitspritedef at offset %d: %d", i, fragData.ParticleRef)
			}

			name := raw.Name(fragData.NameRef)
			if len(name) == 0 {
				name = particle.Name + "_PCD"
			}

			particleInstance := ParticleInstance{
				fragID:                uint32(i),
				Name:                  name,
				Unk1:                  fragData.Unk1,
				Unk2:                  fragData.Unk2,
				ParticleMovement:      fragData.ParticleMovement,
				Flags:                 fragData.Flags,
				SimultaneousParticles: fragData.SimultaneousParticles,
				Unk6:                  fragData.Unk6,
				Unk7:                  fragData.Unk7,
				Unk8:                  fragData.Unk8,
				Unk9:                  fragData.Unk9,
				Unk10:                 fragData.Unk10,
				SpawnRadius:           fragData.SpawnRadius,
				SpawnAngle:            fragData.SpawnAngle,
				SpawnLifespan:         fragData.SpawnLifespan,
				SpawnVelocity:         fragData.SpawnVelocity,
				SpawnNormalZ:          fragData.SpawnNormalZ,
				SpawnNormalX:          fragData.SpawnNormalX,
				SpawnNormalY:          fragData.SpawnNormalY,
				SpawnRate:             fragData.SpawnRate,
				SpawnScale:            fragData.SpawnScale,
				Color: RGBA{
					R: fragData.Color.R,
					G: fragData.Color.G,
					B: fragData.Color.B,
					A: fragData.Color.A,
				},
				Particle: particle.Name,
			}

			wld.ParticleInstances = append(wld.ParticleInstances, &particleInstance)
		case raw.FragCodeMaterialDef:
			fragData, ok := fragment.(*raw.WldFragMaterialDef)
			if !ok {
				return fmt.Errorf("invalid materialdef fragment at offset %d", i)
			}

			sprite := wld.spriteInstanceByFragID(fragData.SpriteInstanceRef)
			if sprite == nil {
				return fmt.Errorf("materialdef found without matching sprite at offset %d", i)
			}

			name := raw.Name(fragData.NameRef)
			if len(name) == 0 {
				name = fmt.Sprintf("%d_MDF", i)
			}

			material := &Material{
				fragID:        uint32(i),
				Name:          name,
				Flags:         fragData.Flags,
				RenderMethod:  fragData.RenderMethod,
				RGBPen:        fragData.RGBPen,
				Brightness:    fragData.Brightness,
				ScaledAmbient: fragData.ScaledAmbient,
				Texture:       sprite.Name,
				Pairs:         fragData.Pairs,
			}
			wld.Materials = append(wld.Materials, material)
		case raw.FragCodeMaterialPalette:
			fragData, ok := fragment.(*raw.WldFragMaterialPalette)
			if !ok {
				return fmt.Errorf("invalid materialpalette fragment at offset %d", i)
			}

			materials := []string{}
			for _, materialRef := range fragData.MaterialRefs {
				material := wld.materialByFragID(materialRef)
				if material == nil {
					return fmt.Errorf("materialpalette found without matching materialdef at offset %d", i)
				}
				materials = append(materials, material.Name)
			}

			materialInstance := MaterialInstance{
				fragID:    uint32(i),
				Name:      raw.Name(fragData.NameRef),
				Flags:     fragData.Flags,
				Materials: materials,
			}

			wld.MaterialInstances = append(wld.MaterialInstances, &materialInstance)
		case raw.FragCodeDmSpriteDef2:
			fragData, ok := fragment.(*raw.WldFragDmSpriteDef2)
			if !ok {
				return fmt.Errorf("invalid dmspritedef2 fragment at offset %d", i)
			}

			materialInstance := wld.materialInstanceByFragID(uint32(fragData.MaterialPaletteRef))
			if materialInstance == nil {
				return fmt.Errorf("dmspritedef2 found without matching materialpalette at offset %d", i)
			}

			// animationInstance := wld.animationInstanceByFragID(uint32(fragData.AnimationRef))
			// if animationInstance == nil {
			// 	return fmt.Errorf("dmspritedef2 found without matching animationinstance at offset %d", i)
			// }
			mesh := &Mesh{
				fragID:            uint32(i),
				Name:              raw.Name(fragData.NameRef),
				Flags:             fragData.Flags,
				MaterialInstance:  materialInstance.Name,
				AnimationInstance: "", //animationInstance.Name,
				Fragment3Ref:      fragData.Fragment3Ref,
				Fragment4Ref:      fragData.Fragment4Ref,
				Center: Vector3{
					X: fragData.Center.X,
					Y: fragData.Center.Y,
					Z: fragData.Center.Z,
				},
				Params2: UIndex3{
					X: fragData.Params2.X,
					Y: fragData.Params2.Y,
					Z: fragData.Params2.Z,
				},
				MaxDistance: fragData.MaxDistance,
				Min: Vector3{
					X: fragData.Min.X,
					Y: fragData.Min.Y,
					Z: fragData.Min.Z,
				},
				Max: Vector3{
					X: fragData.Max.X,
					Y: fragData.Max.Y,
					Z: fragData.Max.Z,
				},
				RawScale:    fragData.RawScale,
				MeshopCount: fragData.MeshopCount,
				Scale:       fragData.Scale,
				Vertices:    fragData.Vertices,
				UVs:         fragData.UVs,
				Normals:     fragData.Normals,
			}

			wld.Meshes = append(wld.Meshes, mesh)
		case raw.FragCodeTrackDef:
			fragData, ok := fragment.(*raw.WldFragTrackDef)
			if !ok {
				return fmt.Errorf("invalid trackdef fragment at offset %d", i)
			}

			name := raw.Name(fragData.NameRef)
			if len(name) == 0 {
				name = fmt.Sprintf("%d_TRACKDEF", i)
			}

			animation := &Animation{
				fragID: uint32(i),
				Name:   name,
				Flags:  fragData.Flags,
			}
			for _, transform := range fragData.BoneTransforms {
				animation.Transforms = append(animation.Transforms, &AnimationTransform{
					RotateDenominator: transform.RotateDenominator,
					RotateX:           transform.RotateX,
					RotateY:           transform.RotateY,
					RotateZ:           transform.RotateZ,
					ShiftX:            transform.ShiftX,
					ShiftY:            transform.ShiftY,
					ShiftZ:            transform.ShiftZ,
				})
			}
			wld.Animations = append(wld.Animations, animation)
		case raw.FragCodeTrack:
			fragData, ok := fragment.(*raw.WldFragTrack)
			if !ok {
				return fmt.Errorf("invalid track fragment at offset %d", i)
			}

			animation := wld.animationByFragID(uint32(fragData.Track))
			if animation == nil {
				return fmt.Errorf("track found without matching trackdef at offset %d", i)
			}

			name := raw.Name(fragData.NameRef)
			if len(name) == 0 {
				name = animation.Name + "_TRACK"
			}

			animationInstance := AnimationInstance{
				fragID:    uint32(i),
				Name:      name,
				Flags:     fragData.Flags,
				Animation: animation.Name,
				Sleep:     fragData.Sleep,
			}
			wld.AnimationInstances = append(wld.AnimationInstances, &animationInstance)
		case raw.FragCodeDMSprite:
			fragData, ok := fragment.(*raw.WldFragDMSprite)
			if !ok {
				return fmt.Errorf("invalid dmsprite fragment at offset %d", i)
			}

			mesh := wld.meshByFragID(uint32(fragData.DMSpriteRef))
			if mesh == nil {
				return fmt.Errorf("dmsprite found without matching dmspritedef2 at offset %d", i)
			}

			name := raw.Name(fragData.NameRef)
			if len(name) == 0 {
				name = mesh.Name + "_INST"
			}

			meshInstance := MeshInstance{
				fragID: uint32(i),
				Name:   name,
				Mesh:   mesh.Name,
				Params: fragData.Params,
			}
			wld.MeshInstances = append(wld.MeshInstances, &meshInstance)
		}

	}

	return nil
}
