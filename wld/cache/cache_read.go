package cache

import (
	"fmt"

	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

func (cm *CacheManager) Load(src *raw.Wld) error {
	cm.FileName = src.FileName()
	cm.Version = src.Version
	for i := 1; i < len(src.Fragments)+1; i++ {
		fragment := src.Fragments[i-1]
		//log.Println("Fragment: ", raw.FragName(fragment.FragCode()), i)

		switch fragment.FragCode() {
		case rawfrag.FragCodeGlobalAmbientLightDef: // turns to globalambientlight
			fragData, ok := fragment.(*rawfrag.WldFragGlobalAmbientLightDef)
			if !ok {
				return fmt.Errorf("invalid globalambientlightdef fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_GALD", i)
			}

			if cm.GlobalAmbientLight != "" {
				return fmt.Errorf("multiple globalambientlight found at offset %d", i)
			}
			cm.GlobalAmbientLight = tag
		case rawfrag.FragCodeBMInfo: // turns to bitmap
			fragData, ok := fragment.(*rawfrag.WldFragBMInfo)
			if !ok {
				return fmt.Errorf("invalid bminfo fragment at offset %d", i)
			}
			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_BMINFO", i)
			}

			bitmap := &Bitmap{
				fragID:   uint32(i),
				Tag:      tag,
				Textures: fragData.TextureNames,
			}
			cm.Bitmaps = append(cm.Bitmaps, bitmap)
		case rawfrag.FragCodeSimpleSpriteDef: // turns to sprite
			fragData, ok := fragment.(*rawfrag.WldFragSimpleSpriteDef)
			if !ok {
				return fmt.Errorf("invalid simplespritedef fragmentat offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_SPRITE", i)
			}

			sprite := SimpleSpriteDef{
				fragID:       uint32(i),
				Tag:          tag,
				Flags:        fragData.Flags,
				CurrentFrame: fragData.CurrentFrame,
				Sleep:        fragData.Sleep,
			}
			for _, bitmapRef := range fragData.BitmapRefs {
				bitmap := cm.bitmapByFragID(bitmapRef)
				if bitmap == nil {
					return fmt.Errorf("simple sprite found without matching bminfo at offset %d ref %d", i, bitmapRef)
				}

				sprite.BMInfos = append(sprite.BMInfos, [2]string{bitmap.Tag, bitmap.Textures[0]})
			}
			cm.SimpleSpriteDefs = append(cm.SimpleSpriteDefs, &sprite)
		case rawfrag.FragCodeSimpleSprite: // turns to spriteinstance
			fragData, ok := fragment.(*rawfrag.WldFragSimpleSprite)
			if !ok {
				return fmt.Errorf("invalid simplesprite fragment at offset %d", i)
			}

			sprite := cm.spriteByFragID(uint32(fragData.SpriteRef))
			if sprite == nil {
				return fmt.Errorf("simplesprite found without matching simplespritedef at offset %d, ref: %d", i, fragData.SpriteRef)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = sprite.Tag
			}
			spriteInstance := SpriteInstance{
				fragID: uint32(i),
				Tag:    tag,
				Flags:  fragData.Flags,
				Sprite: sprite.Tag,
			}
			cm.SpriteInstances = append(cm.SpriteInstances, &spriteInstance)
		case rawfrag.FragCodeBlitSpriteDef: // turns to particle
			fragData, ok := fragment.(*rawfrag.WldFragBlitSpriteDef)
			if !ok {
				return fmt.Errorf("invalid blitspritedef fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_SPB", i)
			}

			particle := &Particle{
				fragID:    uint32(i),
				Tag:       tag,
				Flags:     fragData.Flags,
				SpriteTag: raw.Name(int32(fragData.SpriteInstanceRef)),
				Unknown:   fragData.Unknown,
			}
			cm.Particles = append(cm.Particles, particle)
		case rawfrag.FragCodeParticleCloudDef: // turns to particleinstance
			fragData, ok := fragment.(*rawfrag.WldFragParticleCloudDef)
			if !ok {
				return fmt.Errorf("invalid particleclouddef fragment at offset %d", i)
			}

			particle := cm.particleByFragID(uint32(fragData.ParticleRef))
			if particle == nil {
				return fmt.Errorf("particleclouddef found without matching blitspritedef at offset %d: %d", i, fragData.ParticleRef)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = particle.Tag + "_PCD"
			}

			particleInstance := ParticleInstance{
				fragID:                uint32(i),
				Tag:                   tag,
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
				Color:                 fragData.Color,
				particle:              particle,
			}

			cm.ParticleInstances = append(cm.ParticleInstances, &particleInstance)
		case rawfrag.FragCodeMaterialDef: // turns to material
			fragData, ok := fragment.(*rawfrag.WldFragMaterialDef)
			if !ok {
				return fmt.Errorf("invalid materialdef fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_MDF", i)
			}

			var spriteInstance *SpriteInstance
			if fragData.SimpleSpriteRef > 0 {
				spriteInstance = cm.spriteInstanceByFragID(fragData.SimpleSpriteRef)
				if spriteInstance == nil {
					return fmt.Errorf("materialdef %s refers to missing spriteInstance %d at offset %d", tag, fragData.SimpleSpriteRef, i)
				}
			}

			material := &MaterialDef{
				fragID:         uint32(i),
				Tag:            tag,
				Flags:          fragData.Flags,
				RGBPen:         fragData.RGBPen,
				Brightness:     fragData.Brightness,
				ScaledAmbient:  fragData.ScaledAmbient,
				Pair1:          fragData.Pair1,
				Pair2:          fragData.Pair2,
				spriteInstance: spriteInstance,
			}
			if spriteInstance != nil {
				material.SimpleSpriteInstTag = spriteInstance.Tag
				material.SimpleSpriteInstFlag = spriteInstance.Flags
			}

			material.RenderMethod = model.RenderMethodStr(fragData.RenderMethod)

			cm.MaterialDefs = append(cm.MaterialDefs, material)
		case rawfrag.FragCodeMaterialPalette: // turns to materialinstance
			fragData, ok := fragment.(*rawfrag.WldFragMaterialPalette)
			if !ok {
				return fmt.Errorf("invalid materialpalette fragment at offset %d", i)
			}

			materials := []string{}
			for _, materialRef := range fragData.MaterialRefs {
				material := cm.materialByFragID(materialRef)
				if material == nil {
					return fmt.Errorf("materialInstance found without matching materialdef at offset %d", i)
				}
				materials = append(materials, material.Tag)
			}

			materialInstance := MaterialPalette{
				fragID:    uint32(i),
				Tag:       raw.Name(fragData.NameRef),
				Flags:     fragData.Flags,
				Materials: materials,
			}

			cm.MaterialPalettes = append(cm.MaterialPalettes, &materialInstance)
		case rawfrag.FragCodeDmSpriteDef2: // turns to mesh
			fragData, ok := fragment.(*rawfrag.WldFragDmSpriteDef2)
			if !ok {
				return fmt.Errorf("invalid dmspritedef2 fragment at offset %d", i)
			}

			materialInstance := cm.materialInstanceByFragID(uint32(fragData.MaterialPaletteRef))
			if materialInstance == nil {
				return fmt.Errorf("dmspritedef2 found without matching materialInstance at offset %d", i)
			}

			// animationInstance := cm.animationInstanceByFragID(uint32(fragData.AnimationRef))
			// if animationInstance == nil {
			// 	return fmt.Errorf("dmspritedef2 found without matching animationinstance at offset %d", i)
			// }
			mesh := &DmSpriteDef2{
				fragID:               uint32(i),
				Tag:                  raw.Name(fragData.NameRef),
				Flags:                fragData.Flags,
				MaterialPaletteTag:   materialInstance.Tag,
				DmTrackTag:           "", //animationInstance.Tag,
				Fragment3Ref:         fragData.Fragment3Ref,
				Fragment4Ref:         fragData.Fragment4Ref,
				CenterOffset:         fragData.CenterOffset,
				Params2:              fragData.Params2,
				MaxDistance:          fragData.MaxDistance,
				Min:                  fragData.Min,
				Max:                  fragData.Max,
				Scale:                fragData.Scale,
				Vertices:             fragData.Vertices,
				UVs:                  fragData.UVs,
				VertexNormals:        fragData.VertexNormals,
				Colors:               fragData.Colors,
				FaceMaterialGroups:   fragData.FaceMaterialGroups,
				VertexMaterialGroups: fragData.VertexMaterialGroups,
			}
			for _, face := range fragData.Faces {
				mesh.Faces = append(mesh.Faces, Face{
					Index: face.Index,
					Flags: face.Flags,
				})
			}
			for _, mop := range fragData.MeshOps {
				mesh.MeshOps = append(mesh.MeshOps, MeshOp{
					Index1:    mop.Index1,
					Index2:    mop.Index2,
					Offset:    mop.Offset,
					Param1:    mop.Param1,
					TypeField: mop.TypeField,
				})
			}

			cm.DmSpriteDef2s = append(cm.DmSpriteDef2s, mesh)
		case rawfrag.FragCodeTrackDef: // turns to animation
			fragData, ok := fragment.(*rawfrag.WldFragTrackDef)
			if !ok {
				return fmt.Errorf("invalid trackdef fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_TRACKDEF", i)
			}

			animation := &Animation{
				fragID: uint32(i),
				Tag:    tag,
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
			cm.Animations = append(cm.Animations, animation)
		case rawfrag.FragCodeTrack: // turns to animationinstance
			fragData, ok := fragment.(*rawfrag.WldFragTrack)
			if !ok {
				return fmt.Errorf("invalid track fragment at offset %d", i)
			}

			animation := cm.animationByFragID(uint32(fragData.Track))
			if animation == nil {
				return fmt.Errorf("track found without matching trackdef at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = animation.Tag + "_TRACK"
			}

			animationInstance := AnimationInstance{
				fragID:    uint32(i),
				Tag:       tag,
				Flags:     fragData.Flags,
				Animation: animation.Tag,
				Sleep:     fragData.Sleep,
			}
			cm.AnimationInstances = append(cm.AnimationInstances, &animationInstance)
		case rawfrag.FragCodeDMSprite: // turns to meshinstance
			fragData, ok := fragment.(*rawfrag.WldFragDMSprite)
			if !ok {
				return fmt.Errorf("invalid dmsprite fragment at offset %d", i)
			}

			meshTag := ""
			if fragData.DMSpriteRef > 0 {
				mesh := cm.meshByFragID(uint32(fragData.DMSpriteRef))
				if mesh == nil {
					altMesh := cm.alternateMeshByFragID(uint32(fragData.DMSpriteRef))
					if altMesh == nil {
						return fmt.Errorf("dmsprite found without matching mesh or alternatemesh at offset %d value %d", i, fragData.DMSpriteRef)
					}
					meshTag = altMesh.Tag
				}
				if meshTag == "" {
					meshTag = mesh.Tag
				}
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				if meshTag == "" {
					tag = fmt.Sprintf("%d_DMSPRITE", i)
				} else {
					tag = meshTag + "_DMSPRITE"
				}
			}

			meshInstance := MeshInstance{
				fragID: uint32(i),
				Tag:    tag,
				Mesh:   meshTag,
				Params: fragData.Params,
			}
			cm.MeshInstances = append(cm.MeshInstances, &meshInstance)
		case rawfrag.FragCodeDMSpriteDef: // turns to alternatemesh
			fragData, ok := fragment.(*rawfrag.WldFragDMSpriteDef)
			if !ok {
				return fmt.Errorf("invalid dmspritedef fragment at offset %d", i)
			}

			materialTag := ""
			if fragData.MaterialReference > 0 {
				materialInstance := cm.materialInstanceByFragID(uint32(fragData.MaterialReference))
				if materialInstance == nil {
					return fmt.Errorf("dmspritedef found without matching materialInstance at offset %d", i)
				}
				materialTag = materialInstance.Tag
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_DMSPRITEDEF", i)
			}

			alternateMesh := &AlternateMesh{
				fragID:         uint32(i),
				Tag:            tag,
				Flags:          fragData.Flags,
				Fragment1Maybe: fragData.Fragment1Maybe,
				Material:       materialTag,
				Fragment3:      fragData.Fragment3,
				CenterPosition: fragData.CenterPosition,
				Params2:        fragData.Params2,
				Something2:     fragData.Something2,
				Something3:     fragData.Something3,
				Verticies:      fragData.Vertices,
				TexCoords:      fragData.TexCoords,
				Normals:        fragData.Normals,
				Colors:         fragData.Colors,
				PostVertexFlag: fragData.PostVertexFlag,
				VertexTex:      fragData.VertexTex,
			}

			for _, polygon := range fragData.Polygons {
				alternateMesh.Polygons = append(alternateMesh.Polygons, &AlternateMeshSpritePolygon{
					Flag: polygon.Flag,
					Unk1: polygon.Unk1,
					Unk2: polygon.Unk2,
					Unk3: polygon.Unk3,
					Unk4: polygon.Unk4,
					I1:   polygon.I1,
					I2:   polygon.I2,
					I3:   polygon.I3,
				})
			}

			for _, vertexPiece := range fragData.VertexPieces {
				alternateMesh.VertexPieces = append(alternateMesh.VertexPieces, &AlternateMeshVertexPiece{
					Count:  vertexPiece.Count,
					Offset: vertexPiece.Offset,
				})
			}

			for _, renderGroup := range fragData.RenderGroups {
				alternateMesh.RenderGroups = append(alternateMesh.RenderGroups, &AlternateMeshRenderGroup{
					PolygonCount: renderGroup.PolygonCount,
					MaterialId:   renderGroup.MaterialId,
				})
			}

			for _, size6Piece := range fragData.Size6Pieces {
				alternateMesh.Size6Pieces = append(alternateMesh.Size6Pieces, &AlternateMeshSize6Entry{
					Unk1: size6Piece.Unk1,
					Unk2: size6Piece.Unk2,
					Unk3: size6Piece.Unk3,
					Unk4: size6Piece.Unk4,
					Unk5: size6Piece.Unk5,
				})
			}

			cm.AlternateMeshes = append(cm.AlternateMeshes, alternateMesh)
		case rawfrag.FragCodeActorDef: // turns to actor
			fragData, ok := fragment.(*rawfrag.WldFragActorDef)
			if !ok {
				return fmt.Errorf("invalid actordef fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_ACTORDEF", i)
			}

			actor := &Actor{
				fragID:           uint32(i),
				Tag:              tag,
				Flags:            fragData.Flags,
				CallbackTag:      raw.Name(fragData.CallbackNameRef),
				ActionCount:      fragData.ActionCount,
				FragmentRefCount: fragData.FragmentRefCount,
				BoundsRef:        fragData.BoundsRef,
				CurrentAction:    fragData.CurrentAction,
				Offset:           fragData.Offset,
				Rotation:         fragData.Rotation,
				Unk1:             fragData.Unk1,
				FragmentRefs:     fragData.FragmentRefs,
				Unk2:             fragData.Unk2,
			}

			for _, action := range fragData.Actions {
				actor.Actions = append(actor.Actions, ActorAction{
					LodCount: action.LodCount,
					Unk1:     action.Unk1,
					Lods:     action.Lods,
				})
			}

			cm.Actors = append(cm.Actors, actor)
		case rawfrag.FragCodeActor: // turns to actorinstance
			fragData, ok := fragment.(*rawfrag.WldFragActor)
			if !ok {
				return fmt.Errorf("invalid actor fragment at offset %d", i)
			}

			actorTag := ""
			if fragData.ActorDefRef != -1 {
				actorTag = raw.Name(fragData.ActorDefRef)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				if actorTag == "" {
					tag = fmt.Sprintf("%d_ACTOR_INST", i)
				} else {
					tag = actorTag + "_INST"
				}
			}

			actorInstance := ActorInstance{
				fragID:         uint32(i),
				Tag:            tag,
				ActorTag:       actorTag,
				Flags:          fragData.Flags,
				Sphere:         "",
				CurrentAction:  fragData.CurrentAction,
				Offset:         fragData.Offset,
				Rotation:       fragData.Rotation,
				Unk1:           fragData.Unk1,
				BoundingRadius: fragData.BoundingRadius,
				Scale:          fragData.Scale,
				Sound:          "",
				Unk2:           fragData.Unk2,
			}

			cm.ActorInstances = append(cm.ActorInstances, &actorInstance)
		case rawfrag.FragCodeHierarchialSpriteDef: // turns to skeleton
			fragData, ok := fragment.(*rawfrag.WldFragHierarchialSpriteDef)
			if !ok {
				return fmt.Errorf("invalid hierarchialspritedef fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_HS_DEF", i)
			}

			skeleton := &Skeleton{
				fragID:             uint32(i),
				Tag:                tag,
				Flags:              fragData.Flags,
				CollisionVolumeRef: fragData.CollisionVolumeRef,
				CenterOffset:       fragData.CenterOffset,
				BoundingRadius:     fragData.BoundingRadius,
				Skins:              fragData.Skins,
				SkinLinks:          fragData.SkinLinks,
			}
			for _, bone := range fragData.Bones {
				trackTag := ""
				if bone.TrackRef > 0 {
					track := cm.animationInstanceByFragID(uint32(bone.TrackRef))
					if track == nil {
						return fmt.Errorf("hierarchialspritedef found without matching track at offset %d", i)
					}
					trackTag = track.Tag
				}

				meshOrSpriteOrParticleTag := ""
				if bone.MeshOrSpriteOrParticleRef > 0 {
					meshInstance := cm.meshInstanceByFragID(uint32(bone.MeshOrSpriteOrParticleRef))
					if meshInstance == nil {
						sprite := cm.spriteByFragID(uint32(bone.MeshOrSpriteOrParticleRef))
						if sprite == nil {
							particleInstance := cm.particleInstanceByFragID(uint32(bone.MeshOrSpriteOrParticleRef))
							if particleInstance == nil {
								return fmt.Errorf("hierarchialspritedef found without matching mesh or sprite or particle at offset %d value %d", i, bone.MeshOrSpriteOrParticleRef)
							}
							meshOrSpriteOrParticleTag = particleInstance.Tag
						}
						if meshOrSpriteOrParticleTag == "" {
							meshOrSpriteOrParticleTag = sprite.Tag
						}
					}
					if meshOrSpriteOrParticleTag == "" {
						meshOrSpriteOrParticleTag = meshInstance.Tag
					}
				}

				entry := &SkeletonEntry{
					Tag:          raw.Name(bone.NameRef),
					Flags:        bone.Flags,
					Track:        trackTag,
					MeshOrSprite: meshOrSpriteOrParticleTag,
					SubBones:     bone.SubBones,
				}
				skeleton.Bones = append(skeleton.Bones, entry)
			}
			cm.Skeletons = append(cm.Skeletons, skeleton)
		case rawfrag.FragCodeHierarchialSprite: // turns to skeletoninstance
			fragData, ok := fragment.(*rawfrag.WldFragHierarchialSprite)
			if !ok {
				return fmt.Errorf("invalid hierarchialsprite fragment at offset %d", i)
			}

			skeletonTag := ""
			if fragData.HierarchialSpriteRef > 0 {
				skeleton := cm.skeletonByFragID(uint32(fragData.HierarchialSpriteRef))
				if skeleton == nil {
					return fmt.Errorf("hierarchialsprite found without matching hierarchialspritedef at offset %d value %d", i, fragData.HierarchialSpriteRef)
				}
			}

			tag := raw.Name(int32(fragData.HierarchialSpriteRef))
			if len(tag) == 0 {
				if skeletonTag == "" {
					tag = fmt.Sprintf("%d_HS_INST", i)
				} else {
					tag = skeletonTag + "_INST"
				}
			}

			skeletonInstance := SkeletonInstance{
				fragID:   uint32(i),
				Tag:      tag,
				Skeleton: skeletonTag,
				Flags:    fragData.Flags,
			}
			cm.SkeletonInstances = append(cm.SkeletonInstances, &skeletonInstance)
		case rawfrag.FragCodeLightDef: // turns to light
			fragData, ok := fragment.(*rawfrag.WldFragLightDef)
			if !ok {
				return fmt.Errorf("invalid lightdef fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_LIGHTDEF", i)
			}

			light := &Light{
				fragID:          uint32(i),
				Tag:             tag,
				Flags:           fragData.Flags,
				FrameCurrentRef: fragData.FrameCurrentRef,
				Levels:          fragData.LightLevels,
				Colors:          fragData.Colors,
			}
			cm.Lights = append(cm.Lights, light)
		case rawfrag.FragCodeLight: // turns to lightinstance
			fragData, ok := fragment.(*rawfrag.WldFragLight)
			if !ok {
				return fmt.Errorf("invalid light fragment at offset %d", i)
			}

			lightTag := ""
			if fragData.LightDefRef > 0 {
				light := cm.lightByFragID(uint32(fragData.LightDefRef))
				if light == nil {
					return fmt.Errorf("light found without matching lightdef at offset %d", i)
				}
				lightTag = light.Tag
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				if lightTag == "" {
					tag = fmt.Sprintf("%d_LIGHT", i)
				} else {
					tag = lightTag + "_INST"
				}
			}

			lightInstance := LightInstance{
				fragID: uint32(i),
				Tag:    tag,
				Light:  lightTag,
				Flags:  fragData.Flags,
			}

			cm.LightInstances = append(cm.LightInstances, &lightInstance)
		case rawfrag.FragCodeSprite3DDef: // turns to camera
			fragData, ok := fragment.(*rawfrag.WldFragSprite3DDef)
			if !ok {
				return fmt.Errorf("invalid sprite3ddef fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_SPRITE3DDEF", i)
			}

			camera := &Camera{
				fragID:        uint32(i),
				Tag:           tag,
				Flags:         fragData.Flags,
				SphereListRef: fragData.SphereListRef,
				CenterOffset:  fragData.CenterOffset,
				Radius:        fragData.Radius,
				Vertices:      fragData.Vertices,
			}

			for _, bspNode := range fragData.BspNodes {
				node := &CameraBspNode{
					FrontTree:                   bspNode.FrontTree,
					BackTree:                    bspNode.BackTree,
					VertexIndexes:               bspNode.VertexIndexes,
					RenderMethod:                bspNode.RenderMethod,
					RenderFlags:                 bspNode.RenderFlags,
					RenderPen:                   bspNode.RenderPen,
					RenderBrightness:            bspNode.RenderBrightness,
					RenderScaledAmbient:         bspNode.RenderScaledAmbient,
					RenderSimpleSpriteReference: bspNode.RenderSimpleSpriteReference,
					RenderUVInfoOrigin:          bspNode.RenderUVInfoOrigin,
					RenderUVInfoUAxis:           bspNode.RenderUVInfoUAxis,
					RenderUVInfoVAxis:           bspNode.RenderUVInfoVAxis,
				}

				for _, uvMap := range bspNode.RenderUVMapEntries {
					entry := CameraBspNodeUVMapEntry{
						UvOrigin: uvMap.UvOrigin,
						UAxis:    uvMap.UAxis,
						VAxis:    uvMap.VAxis,
					}
					node.RenderUVMapEntries = append(node.RenderUVMapEntries, entry)
				}
				camera.BspNodes = append(camera.BspNodes, node)

				cm.Cameras = append(cm.Cameras, camera)
			}
		case rawfrag.FragCodeSprite3D: // turns to camerainstance
			fragData, ok := fragment.(*rawfrag.WldFragSprite3D)
			if !ok {
				return fmt.Errorf("invalid sprite3d fragment at offset %d", i)
			}

			cameraTag := ""
			if fragData.Sprite3DDefRef > 0 {
				camera := cm.cameraByFragID(uint32(fragData.Sprite3DDefRef))
				if camera == nil {
					return fmt.Errorf("sprite3d found without matching sprite3ddef at offset %d", i)
				}
				cameraTag = camera.Tag
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				if cameraTag == "" {
					tag = fmt.Sprintf("%d_SPRITE3D", i)
				} else {
					tag = cameraTag + "_INST"
				}
			}

			cameraInstance := CameraInstance{
				fragID:    uint32(i),
				Tag:       tag,
				CameraTag: cameraTag,
				Flags:     fragData.Flags,
			}

			cm.CameraInstances = append(cm.CameraInstances, &cameraInstance)

		case rawfrag.FragCodeSphere: // turns to sphere
			fragData, ok := fragment.(*rawfrag.WldFragSphere)
			if !ok {
				return fmt.Errorf("invalid sphere fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_SPHERE", i)
			}

			sphere := &Sphere{
				fragID: uint32(i),
				Tag:    tag,
				Radius: fragData.Radius,
			}

			cm.Spheres = append(cm.Spheres, sphere)
		case rawfrag.FragCodeZone: // turns to regioninstance
			fragData, ok := fragment.(*rawfrag.WldFragZone)
			if !ok {
				return fmt.Errorf("invalid zone fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_ZONE", i)
			}

			regionInstance := RegionInstance{
				fragID:   uint32(i),
				Tag:      tag,
				Flags:    fragData.Flags,
				UserData: fragData.UserData,
			}

			for _, regionRef := range fragData.Regions {
				if regionRef < 0 || int(regionRef) >= len(cm.Regions) {
					return fmt.Errorf("zone found with invalid region ref %d at offset %d", regionRef, i)
				}
				region := cm.Regions[regionRef]
				regionInstance.RegionTags = append(regionInstance.RegionTags, region.Tag)
			}

			cm.RegionInstances = append(cm.RegionInstances, &regionInstance)

		case rawfrag.FragCodeWorldTree: // turns to bsptree
			fragData, ok := fragment.(*rawfrag.WldFragWorldTree)
			if !ok {
				return fmt.Errorf("invalid worldtree fragment at offset %d", i)
			}

			if len(cm.BspTrees) > 0 {
				return fmt.Errorf("multiple worldtree found at offset %d", i)
			}

			tag := ""
			if fragData.NameRef > 0 {
				tag = raw.Name(fragData.NameRef)
			}
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_WORLD_TREE", i)
			}

			bspTree := &BspTree{
				fragID: uint32(i),
				Tag:    tag,
			}

			for _, node := range fragData.Nodes {
				regionTag := raw.Name(node.RegionRef)
				if len(regionTag) == 0 {
					regionTag = fmt.Sprintf("%d_BSP_TREE", i)
				}

				bspTree.Nodes = append(bspTree.Nodes, &BspTreeNode{
					Normal:    node.Normal,
					Distance:  node.Distance,
					RegionTag: regionTag,
				})
			}

			for i := 0; i < len(fragData.Nodes); i++ {
				fragRegion := fragData.Nodes[i]
				node := bspTree.Nodes[i]
				if fragRegion.FrontRef > 0 {
					if fragRegion.FrontRef-1 >= int32(len(bspTree.Nodes)) {
						return fmt.Errorf("bspTree %s has invalid front ref %d at offset %d", bspTree.Tag, fragRegion.FrontRef, i)
					}

					node.Front = bspTree.Nodes[fragRegion.FrontRef-1]
				}
				if fragRegion.BackRef > 0 {
					if fragRegion.BackRef-1 >= int32(len(bspTree.Nodes)) {
						return fmt.Errorf("bspTree %s has invalid back ref %d at offset %d", bspTree.Tag, fragRegion.BackRef, i)
					}

					node.Back = bspTree.Nodes[fragRegion.BackRef-1]
				}
			}

			cm.BspTrees = append(cm.BspTrees, bspTree)
		case rawfrag.FragCodeRegion: // turns to region
			fragData, ok := fragment.(*rawfrag.WldFragRegion)
			if !ok {
				return fmt.Errorf("invalid region fragment at offset %d", i)
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_REGION", i)
			}

			region := &Region{
				fragID:          uint32(i),
				Tag:             tag,
				Flags:           fragData.Flags,
				AmbientLightRef: fragData.AmbientLightRef,
				//RegionVertexCount:    fragData.RegionVertexCount,
				//RegionProximalCount:  fragData.RegionProximalCount,
				//RenderVertexCount:    fragData.RenderVertexCount,
				//WallCount:            fragData.WallCount,
				//ObstacleCount:        fragData.ObstacleCount,
				CuttingObstacleCount: fragData.CuttingObstacleCount,
				//VisibleNodeCount:     fragData.VisibleNodeCount,
				RegionVertices:  fragData.RegionVertices,
				RegionProximals: fragData.RegionProximals,
				RenderVertices:  fragData.RenderVertices,
			}

			for _, wall := range fragData.Walls {
				region.Walls = append(region.Walls, &RegionWall{
					Flags:                       wall.Flags,
					VertexCount:                 wall.VertexCount,
					RenderMethod:                wall.RenderMethod,
					RenderFlags:                 wall.RenderFlags,
					RenderPen:                   wall.RenderPen,
					RenderBrightness:            wall.RenderBrightness,
					RenderScaledAmbient:         wall.RenderScaledAmbient,
					RenderSimpleSpriteReference: wall.RenderSimpleSpriteReference,
					RenderUVInfoOrigin:          wall.RenderUVInfoOrigin,
					RenderUVInfoUAxis:           wall.RenderUVInfoUAxis,
					RenderUVInfoVAxis:           wall.RenderUVInfoVAxis,
					RenderUVMapEntryCount:       wall.RenderUVMapEntryCount,
					RenderUVMapEntries:          wall.RenderUVMapEntries,
					Normal:                      wall.Normal,
					Vertices:                    wall.Vertices,
				})
			}

			cm.Regions = append(cm.Regions, region)
		case rawfrag.FragCodeAmbientLight: // turns to ambientlightinstance
			fragData, ok := fragment.(*rawfrag.WldFragAmbientLight)
			if !ok {
				return fmt.Errorf("invalid ambientlight fragment at offset %d", i)
			}

			lightTag := ""
			if fragData.LightRef > 0 {
				/*
					light := cm.lightByFragID(uint32(fragData.LightRef))
					 if light == nil {
						return fmt.Errorf("ambientlight found without matching light at offset %d", i)
					}
					lightTag = light.Tag */
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_AMBIENTLIGHT", i)
			}

			ambientLightInstance := AmbientLightInstance{
				fragID:   uint32(i),
				Tag:      tag,
				LightTag: lightTag,
				Flags:    fragData.Flags,
			}
			/*
				for _, regionRef := range fragData.Regions {
						region := cm.regionByFragID(regionRef)
					if region == nil {
						return fmt.Errorf("ambientlight found without matching region at offset %d value %d", i, regionRef)
					}

					ambientLightInstance.RegionTags = append(ambientLightInstance.RegionTags, region.Tag)
				}
			*/

			cm.AmbientLightInstances = append(cm.AmbientLightInstances, &ambientLightInstance)
		case rawfrag.FragCodePointLight: // turns to pointlightinstance
			fragData, ok := fragment.(*rawfrag.WldFragPointLight)
			if !ok {
				return fmt.Errorf("invalid pointlight fragment at offset %d", i)
			}

			lightInstanceTag := ""
			if fragData.LightRef > 0 {
				lightInstance := cm.lightInstanceByFragID(uint32(fragData.LightRef))
				if lightInstance == nil {
					return fmt.Errorf("pointlight found without matching light at offset %d value %d", i, fragData.LightRef)
				}
				lightInstanceTag = lightInstance.Tag
			}

			tag := raw.Name(fragData.NameRef)
			if len(tag) == 0 {
				tag = fmt.Sprintf("%d_POINTLIGHT", i)
			}

			pointLightInstance := PointLightInstance{
				fragID:           uint32(i),
				Tag:              tag,
				LightInstanceTag: lightInstanceTag,
				Flags:            fragData.Flags,
				X:                fragData.X,
				Y:                fragData.Y,
				Z:                fragData.Z,
				Radius:           fragData.Radius,
			}

			cm.PointLightInstances = append(cm.PointLightInstances, &pointLightInstance)
		default:
			return fmt.Errorf("unknown fragment type 0x%x (%s) at offset %d", fragment.FragCode(), raw.FragName(fragment.FragCode()), i)
		}

	}
	return nil
}
