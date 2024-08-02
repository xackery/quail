package wld

import (
	"fmt"
	"io"

	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

func (wld *Wld) ReadRaw(src *raw.Wld) error {
	wld.reset()
	for i := 1; i < len(src.Fragments); i++ {
		fragment := src.Fragments[i]
		err := readRawFrag(wld, src, fragment)
		if err != nil {
			return fmt.Errorf("fragment %d (%s): %w", i, raw.FragName(fragment.FragCode()), err)
		}
	}

	return nil
}

func readRawFrag(wld *Wld, src *raw.Wld, fragment model.FragmentReadWriter) error {
	i := 0

	switch fragment.FragCode() {
	case rawfrag.FragCodeGlobalAmbientLightDef:
		fragData, ok := fragment.(*rawfrag.WldFragGlobalAmbientLightDef)
		if !ok {
			return fmt.Errorf("invalid globalambientlightdef fragment at offset %d", i)
		}
		if wld.GlobalAmbientLight != nil {
			return fmt.Errorf("duplicate globalambientlightdef found")
		}

		wld.GlobalAmbientLight = &GlobalAmbientLightDef{}
		wld.GlobalAmbientLight.Tag = "DEFAULT_AMBIENTLIGHT"
		if fragData.NameRef != -16777216 {
			wld.GlobalAmbientLight.Tag = raw.Name(fragData.NameRef)
		}

	case rawfrag.FragCodeBMInfo:
		return nil
	case rawfrag.FragCodeSimpleSpriteDef:
		fragData, ok := fragment.(*rawfrag.WldFragSimpleSpriteDef)
		if !ok {
			return fmt.Errorf("invalid simplespritedef fragment at offset %d", i)
		}
		tag := raw.Name(fragData.NameRef)
		if len(tag) == 0 {
			tag = fmt.Sprintf("%d_SPRITEDEF", i)
		}
		sprite := &SimpleSpriteDef{
			Tag: tag,
		}
		if fragData.Flags&0x02 == 0x02 {
			sprite.SkipFrames.Valid = true
		}
		if fragData.Flags&0x04 == 0x04 {
			sprite.Animated.Valid = true
		}
		if fragData.Flags&0x10 == 0x10 && sprite.Animated.Valid {
			sprite.Sleep.Valid = true
			sprite.Sleep.Uint32 = fragData.Sleep
		}
		if fragData.Flags&0x20 == 0x20 {
			sprite.CurrentFrame.Valid = true
			sprite.CurrentFrame.Int32 = fragData.CurrentFrame
		}

		for _, bitmapRef := range fragData.BitmapRefs {
			if bitmapRef == 0 {
				return nil
			}
			if len(src.Fragments) < int(bitmapRef) {
				return fmt.Errorf("bitmap ref %d not found", bitmapRef)
			}
			bitmap := src.Fragments[bitmapRef]
			bmInfo, ok := bitmap.(*rawfrag.WldFragBMInfo)
			if !ok {
				return fmt.Errorf("invalid bitmap ref %d", bitmapRef)
			}
			sprite.SimpleSpriteFrames = append(sprite.SimpleSpriteFrames, SimpleSpriteFrame{
				TextureTag:  raw.Name(bmInfo.NameRef),
				TextureFile: bmInfo.TextureNames[0],
			})
		}
		wld.SimpleSpriteDefs = append(wld.SimpleSpriteDefs, sprite)
	case rawfrag.FragCodeSimpleSprite:
		//return fmt.Errorf("simplesprite fragment found, but not expected")
	case rawfrag.FragCodeBlitSpriteDef:
		return fmt.Errorf("blitsprite fragment found, but not expected")
	case rawfrag.FragCodeParticleCloudDef:
		return fmt.Errorf("particlecloud fragment found, but not expected")
	case rawfrag.FragCodeMaterialDef:
		fragData, ok := fragment.(*rawfrag.WldFragMaterialDef)
		if !ok {
			return fmt.Errorf("invalid materialdef fragment at offset %d", i)
		}
		spriteTag := ""
		spriteFlags := uint32(0)
		if fragData.SimpleSpriteRef > 0 {
			if len(src.Fragments) < int(fragData.SimpleSpriteRef) {
				return fmt.Errorf("simplesprite ref %d out of bounds", fragData.SimpleSpriteRef)
			}
			simpleSprite, ok := src.Fragments[fragData.SimpleSpriteRef].(*rawfrag.WldFragSimpleSprite)
			if !ok {
				return fmt.Errorf("simplesprite ref %d not found", fragData.SimpleSpriteRef)
			}
			if len(src.Fragments) < int(simpleSprite.SpriteRef) {
				return fmt.Errorf("sprite ref %d out of bounds", simpleSprite.SpriteRef)
			}
			spriteDef, ok := src.Fragments[simpleSprite.SpriteRef].(*rawfrag.WldFragSimpleSpriteDef)
			if !ok {
				return fmt.Errorf("sprite ref %d not found", simpleSprite.SpriteRef)
			}

			spriteTag = raw.Name(spriteDef.NameRef)
			spriteFlags = simpleSprite.Flags
			if spriteFlags != 0x50 {
				fmt.Printf("unknown sprite flag %d (ignored)\n", spriteFlags)
				//	return fmt.Errorf("unknown sprite flag %d", spriteFlags)
			}
		}
		material := &MaterialDef{
			Tag:             raw.Name(fragData.NameRef),
			RenderMethod:    model.RenderMethodStr(fragData.RenderMethod),
			RGBPen:          fragData.RGBPen,
			Brightness:      fragData.Brightness,
			ScaledAmbient:   fragData.ScaledAmbient,
			SimpleSpriteTag: spriteTag,
		}
		if fragData.Flags&0x02 == 0x02 {
			material.Pair1.Valid = true
			material.Pair1.Uint32 = fragData.Pair1
			material.Pair2.Valid = true
			material.Pair2.Float32 = fragData.Pair2
		}
		wld.MaterialDefs = append(wld.MaterialDefs, material)
	case rawfrag.FragCodeMaterialPalette:
		fragData, ok := fragment.(*rawfrag.WldFragMaterialPalette)
		if !ok {
			return fmt.Errorf("invalid materialpalette fragment at offset %d", i)
		}

		tag := raw.Name(fragData.NameRef)
		if len(tag) == 0 {
			tag = fmt.Sprintf("%d_MPL", i)
		}

		materialPalette := &MaterialPalette{
			Tag:   tag,
			flags: fragData.Flags,
		}
		for _, materialRef := range fragData.MaterialRefs {
			if len(src.Fragments) < int(materialRef) {
				return fmt.Errorf("material ref %d not found", materialRef)
			}
			material, ok := src.Fragments[materialRef].(*rawfrag.WldFragMaterialDef)
			if !ok {
				return fmt.Errorf("invalid materialdef fragment at offset %d", materialRef)
			}
			materialPalette.Materials = append(materialPalette.Materials, raw.Name(material.NameRef))
		}
		wld.MaterialPalettes = append(wld.MaterialPalettes, materialPalette)
	case rawfrag.FragCodeDmSpriteDef2:
		fragData, ok := fragment.(*rawfrag.WldFragDmSpriteDef2)
		if !ok {
			return fmt.Errorf("invalid dmspritedef2 fragment at offset %d", i)
		}
		sprite := &DMSpriteDef2{
			Tag:                  raw.Name(fragData.NameRef),
			Flags:                fragData.Flags,
			DmTrackTag:           raw.Name(fragData.DMTrackRef),
			Fragment3Ref:         fragData.Fragment3Ref,
			Fragment4Ref:         fragData.Fragment4Ref,
			CenterOffset:         fragData.CenterOffset,
			Params2:              fragData.Params2,
			MaxDistance:          fragData.MaxDistance,
			Min:                  fragData.Min,
			Max:                  fragData.Max,
			FPScale:              fragData.Scale,
			VertexColors:         fragData.Colors,
			FaceMaterialGroups:   fragData.FaceMaterialGroups,
			VertexMaterialGroups: fragData.VertexMaterialGroups,
		}
		if fragData.MaterialPaletteRef > 0 {
			if len(src.Fragments) < int(fragData.MaterialPaletteRef) {
				return fmt.Errorf("materialpalette ref %d out of bounds", fragData.MaterialPaletteRef)
			}
			materialPalette, ok := src.Fragments[fragData.MaterialPaletteRef].(*rawfrag.WldFragMaterialPalette)
			if !ok {
				return fmt.Errorf("materialpalette ref %d not found", fragData.MaterialPaletteRef)
			}
			sprite.MaterialPaletteTag = raw.Name(materialPalette.NameRef)
		}

		scale := 1.0 / float32(int(1<<fragData.Scale))

		for _, vert := range fragData.Vertices {
			sprite.Vertices = append(sprite.Vertices, [3]float32{
				float32(vert[0]) * scale,
				float32(vert[1]) * scale,
				float32(vert[2]) * scale,
			})
		}
		for _, uv := range fragData.UVs {
			sprite.UVs = append(sprite.UVs, [2]float32{
				float32(uv[0]) * scale,
				float32(uv[1]) * scale,
			})
		}
		for _, vn := range fragData.VertexNormals {
			sprite.VertexNormals = append(sprite.VertexNormals, [3]float32{
				float32(vn[0]) * scale,
				float32(vn[1]) * scale,
				float32(vn[2]) * scale,
			})
		}
		for _, face := range fragData.Faces {
			sprite.Faces = append(sprite.Faces, &Face{
				Flags:    face.Flags,
				Triangle: face.Index,
			})
		}
		for _, mop := range fragData.MeshOps {
			sprite.MeshOps = append(sprite.MeshOps, &MeshOp{
				Index1:    mop.Index1,
				Index2:    mop.Index2,
				Offset:    mop.Offset,
				Param1:    mop.Param1,
				TypeField: mop.TypeField,
			})
		}
		wld.DMSpriteDef2s = append(wld.DMSpriteDef2s, sprite)
	case rawfrag.FragCodeTrackDef:
		fragData, ok := fragment.(*rawfrag.WldFragTrackDef)
		if !ok {
			return fmt.Errorf("invalid trackdef fragment at offset %d", i)
		}

		track := &TrackDef{
			Tag: raw.Name(fragData.NameRef),
		}

		for _, transform := range fragData.FrameTransforms {
			frame := TrackFrameTransform{
				XYZScale: transform.ShiftDenominator,
			}
			scale := 1.0 / float32(int(1<<transform.ShiftDenominator))

			frame.XYZ = [3]float32{
				float32(transform.Shift[0]) / scale,
				float32(transform.Shift[1]) / scale,
				float32(transform.Shift[2]) / scale,
			}

			if fragData.Flags&0x08 == 0x08 {
				frame.RotScale.Valid = true
				frame.RotScale.Int16 = transform.RotateDenominator
				scale = 1.0 / float32(int(1<<transform.RotateDenominator))
				frame.Rotation.Valid = true
				frame.Rotation.Float32Slice3 = [3]float32{
					float32(transform.Rotation[0]) / scale,
					float32(transform.Rotation[1]) / scale,
					float32(transform.Rotation[2]) / scale,
				}
			} else {
				frame.RotScale.Valid = false
				frame.LegacyRotation.Valid = true
				frame.LegacyRotation.Float32Slice4 = [4]float32{
					float32(transform.Rotation[0]) / scale,
					float32(transform.Rotation[1]) / scale,
					float32(transform.Rotation[2]) / scale,
					float32(transform.Rotation[3]) / scale,
				}
			}

			track.FrameTransforms = append(track.FrameTransforms, frame)
		}

		wld.TrackDefs = append(wld.TrackDefs, track)
	case rawfrag.FragCodeTrack:
		fragData, ok := fragment.(*rawfrag.WldFragTrack)
		if !ok {
			return fmt.Errorf("invalid track fragment at offset %d", i)
		}

		if len(src.Fragments) < int(fragData.TrackRef) {
			return fmt.Errorf("trackdef ref %d not found", fragData.TrackRef)
		}

		trackDef, ok := src.Fragments[fragData.TrackRef].(*rawfrag.WldFragTrackDef)
		if !ok {
			return fmt.Errorf("trackdef ref %d not found", fragData.TrackRef)
		}

		trackInst := &TrackInstance{
			Tag:           raw.Name(fragData.NameRef),
			DefinitionTag: raw.Name(trackDef.NameRef),
		}
		if fragData.Flags&0x01 == 0x01 {
			trackInst.Sleep.Valid = true
			trackInst.Sleep.Uint32 = fragData.Sleep
		}
		if fragData.Flags&0x02 == 0x02 {
			trackInst.Reverse = 1
		}
		if fragData.Flags&0x04 == 0x04 {
			trackInst.Interpolate = 1
		}

		wld.TrackInstances = append(wld.TrackInstances, trackInst)
	case rawfrag.FragCodeDMSpriteDef:
		fragData, ok := fragment.(*rawfrag.WldFragDMSpriteDef)
		if !ok {
			return fmt.Errorf("invalid dmspritedef fragment at offset %d", i)
		}

		sprite := &DMSpriteDef{
			Tag:            raw.Name(fragData.NameRef),
			Flags:          fragData.Flags,
			Fragment1Maybe: fragData.Fragment1Maybe,
			Material:       raw.Name(int32(fragData.MaterialReference)),
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
			sprite.Polygons = append(sprite.Polygons, &DMSpriteDefSpritePolygon{
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
			sprite.VertexPieces = append(sprite.VertexPieces, &DMSpriteDefVertexPiece{
				Count:  vertexPiece.Count,
				Offset: vertexPiece.Offset,
			})
		}

		for _, renderGroup := range fragData.RenderGroups {
			sprite.RenderGroups = append(sprite.RenderGroups, &DMSpriteDefRenderGroup{
				PolygonCount: renderGroup.PolygonCount,
				MaterialId:   renderGroup.MaterialId,
			})
		}

		for _, size6Piece := range fragData.Size6Pieces {
			sprite.Size6Pieces = append(sprite.Size6Pieces, &DMSpriteDefSize6Entry{
				Unk1: size6Piece.Unk1,
				Unk2: size6Piece.Unk2,
				Unk3: size6Piece.Unk3,
				Unk4: size6Piece.Unk4,
				Unk5: size6Piece.Unk5,
			})
		}

		wld.DMSpriteDefs = append(wld.DMSpriteDefs, sprite)

	case rawfrag.FragCodeDMSprite:
		fragData, ok := fragment.(*rawfrag.WldFragDMSprite)
		if !ok {
			return fmt.Errorf("invalid dmsprite fragment at offset %d", i)
		}

		if len(src.Fragments) < int(fragData.DMSpriteRef) {
			return fmt.Errorf("dmspritedef ref %d not found", fragData.DMSpriteRef)
		}

		dmSpriteDef, ok := src.Fragments[fragData.DMSpriteRef].(*rawfrag.WldFragDmSpriteDef2)
		if !ok {
			return fmt.Errorf("dmspritedef ref %d not found", fragData.DMSpriteRef)
		}

		dmsprite := &DMSprite{
			Tag:           raw.Name(fragData.NameRef),
			DefinitionTag: raw.Name(dmSpriteDef.NameRef),
			Param:         fragData.Params,
		}

		wld.DMSpriteInsts = append(wld.DMSpriteInsts, dmsprite)
	case rawfrag.FragCodeActorDef:
		fragData, ok := fragment.(*rawfrag.WldFragActorDef)
		if !ok {
			return fmt.Errorf("invalid actordef fragment at offset %d", i)
		}

		actor := &ActorDef{
			Tag:       raw.Name(fragData.NameRef),
			Callback:  raw.Name(fragData.CallbackNameRef),
			BoundsRef: fragData.BoundsRef,
			Unk1:      fragData.Unk1,
		}

		if fragData.Flags&0x01 == 0x01 {
			actor.CurrentAction.Valid = true
			actor.CurrentAction.Uint32 = fragData.CurrentAction
		}
		if fragData.Flags&0x02 == 0x02 {
			actor.Location.Valid = true
			actor.Location.Float32Slice6 = fragData.Location
		}
		if fragData.Flags&0x40 == 0x40 {
			actor.ActiveGeometry.Valid = true
		}

		if len(fragData.Actions) != len(fragData.FragmentRefs) {
			return fmt.Errorf("actordef actions and fragmentrefs mismatch at offset %d", i)
		}

		fragRefIndex := 0
		for _, srcAction := range fragData.Actions {
			lods := []ActorLevelOfDetail{}
			for _, srcLod := range srcAction.Lods {
				spriteTag := ""
				if len(fragData.FragmentRefs) > fragRefIndex {
					spriteRef := fragData.FragmentRefs[fragRefIndex]
					if len(src.Fragments) < int(spriteRef) {
						return fmt.Errorf("actordef fragment ref %d not found at offset %d", spriteRef, i)
					}
					switch sprite := src.Fragments[spriteRef].(type) {
					case *rawfrag.WldFragSprite3D:
						if len(src.Fragments) < int(sprite.Sprite3DDefRef) {
							return fmt.Errorf("sprite3ddef ref %d not found", sprite.Sprite3DDefRef)
						}
						spriteDef, ok := src.Fragments[sprite.Sprite3DDefRef].(*rawfrag.WldFragSprite3DDef)
						if !ok {
							return fmt.Errorf("sprite3ddef ref %d not found", sprite.Sprite3DDefRef)
						}
						spriteTag = raw.Name(spriteDef.NameRef)
					case *rawfrag.WldFragDMSprite:
						if len(src.Fragments) < int(sprite.DMSpriteRef) {
							return fmt.Errorf("dmsprite ref %d not found", sprite.DMSpriteRef)
						}
						spriteDef := src.Fragments[sprite.DMSpriteRef].(*rawfrag.WldFragDmSpriteDef2)
						if !ok {
							return fmt.Errorf("dmsprite ref %d not found", sprite.DMSpriteRef)
						}
						spriteTag = raw.Name(spriteDef.NameRef)
					case *rawfrag.WldFragHierarchicalSprite:
						if len(src.Fragments) < int(sprite.HierarchicalSpriteRef) {
							return fmt.Errorf("hierarchicalsprite def ref %d not found", sprite.HierarchicalSpriteRef)
						}
						spriteDef, ok := src.Fragments[sprite.HierarchicalSpriteRef].(*rawfrag.WldFragHierarchicalSpriteDef)
						if !ok {
							return fmt.Errorf("hierarchicalsprite def ref %d not found", sprite.HierarchicalSpriteRef)
						}
						spriteTag = raw.Name(spriteDef.NameRef)

					default:
						return fmt.Errorf("unhandled sprite instance fragment type %d (%s) at offset %d", sprite.FragCode(), raw.FragName(sprite.FragCode()), i)
					}
				}
				lod := ActorLevelOfDetail{
					SpriteTag:   spriteTag,
					MinDistance: srcLod,
				}

				lods = append(lods, lod)
				fragRefIndex++
			}

			actor.Actions = append(actor.Actions, ActorAction{
				Unk1:           srcAction.Unk1,
				LevelOfDetails: lods,
			})
		}

		wld.ActorDefs = append(wld.ActorDefs, actor)
	case rawfrag.FragCodeActor:
		fragData, ok := fragment.(*rawfrag.WldFragActor)
		if !ok {
			return fmt.Errorf("invalid actor fragment at offset %d", i)
		}

		actorDefTag := ""
		if fragData.ActorDefRef > 0 {
			if len(src.Fragments) < int(fragData.ActorDefRef) {
				return fmt.Errorf("actordef ref %d out of bounds", fragData.ActorDefRef)
			}

			actorDef, ok := src.Fragments[fragData.ActorDefRef].(*rawfrag.WldFragActorDef)
			if !ok {
				return fmt.Errorf("actordef ref %d not found", fragData.ActorDefRef)
			}
			actorDefTag = raw.Name(actorDef.NameRef)
		}

		if len(src.Fragments) < int(fragData.SphereRef) {
			return fmt.Errorf("sphere ref %d not found", fragData.SphereRef)
		}

		sphereRadius := float32(0)
		if fragData.SphereRef > 0 {
			sphereDef, ok := src.Fragments[fragData.SphereRef].(*rawfrag.WldFragSphere)
			if !ok {
				return fmt.Errorf("sphere ref %d not found", fragData.SphereRef)
			}
			sphereRadius = sphereDef.Radius
		}

		actor := &ActorInst{
			Tag:           raw.Name(fragData.NameRef),
			DefinitionTag: actorDefTag,
			SphereRadius:  sphereRadius,
			UserData:      fragData.UserData,
		}

		if fragData.Flags&0x01 == 0x01 {
			actor.CurrentAction.Valid = true
			actor.CurrentAction.Uint32 = fragData.CurrentAction
		}

		if fragData.Flags&0x02 == 0x02 {
			actor.Location.Valid = true
			actor.Location.Float32Slice6 = fragData.Location
		}

		if fragData.Flags&0x04 == 0x04 {
			actor.BoundingRadius.Valid = true
			actor.BoundingRadius.Float32 = fragData.BoundingRadius
		}

		if fragData.Flags&0x08 == 0x08 {
			actor.Scale.Valid = true
			actor.Scale.Float32 = fragData.ScaleFactor
		}

		if fragData.Flags&0x10 == 0x10 {
			actor.SoundTag.Valid = true
			actor.SoundTag.String = raw.Name(fragData.SoundNameRef)
		}

		if fragData.Flags&0x20 == 0x20 {
			actor.Active.Valid = true
		}

		// 0x40 unknown
		if fragData.Flags&0x80 == 0x80 {
			actor.SpriteVolumeOnly.Valid = true
		}

		if fragData.Flags&0x100 == 0x100 {
			actor.DMRGBTrackTag.Valid = true

			trackTag := ""
			if fragData.DMRGBTrackRef == 0 {
				return fmt.Errorf("dmrgbtrack flag set, but ref is 0")
			}
			if len(src.Fragments) < int(fragData.DMRGBTrackRef) {
				return fmt.Errorf("dmrgbtrack ref %d out of bounds", fragData.DMRGBTrackRef)
			}

			track, ok := src.Fragments[fragData.DMRGBTrackRef].(*rawfrag.WldFragDmRGBTrack)
			if !ok {
				return fmt.Errorf("dmrgbtrack ref %d not found", fragData.DMRGBTrackRef)
			}
			if len(src.Fragments) < int(track.TrackRef) {
				return fmt.Errorf("dmrgbtrackdef ref %d not found", track.TrackRef)
			}

			trackDef, ok := src.Fragments[track.TrackRef].(*rawfrag.WldFragDmRGBTrackDef)
			if !ok {
				return fmt.Errorf("dmrgbtrackdef ref %d not found", track.TrackRef)
			}
			if trackDef.NameRef != 0 {
				trackTag = raw.Name(trackDef.NameRef)
			}
			actor.DMRGBTrackTag.String = trackTag
		}

		wld.ActorInsts = append(wld.ActorInsts, actor)
	case rawfrag.FragCodeHierarchicalSpriteDef:
		fragData, ok := fragment.(*rawfrag.WldFragHierarchicalSpriteDef)
		if !ok {
			return fmt.Errorf("invalid hierarchicalsprite fragment at offset %d", i)
		}

		collisionTag := ""
		if fragData.CollisionVolumeNameRef != 0 && fragData.CollisionVolumeNameRef != 4294967293 {
			if len(src.Fragments) < int(fragData.CollisionVolumeNameRef) {
				return fmt.Errorf("collision volume ref %d out of bounds", fragData.CollisionVolumeNameRef)
			}

			switch collision := src.Fragments[fragData.CollisionVolumeNameRef].(type) {
			case *rawfrag.WldFragPolyhedron:
				collisionTag = raw.Name(collision.NameRef)
			default:
				return fmt.Errorf("unknown collision volume ref %d (%s)", fragData.CollisionVolumeNameRef, raw.FragName(collision.FragCode()))
			}
		}
		if fragData.CollisionVolumeNameRef == 4294967293 {
			collisionTag = "SPECIAL_COLLISION"
		}
		if collisionTag != "" {
			return fmt.Errorf("collision volume ref found as %s, report this to xack", collisionTag)
		}

		spriteDef := &HierarchicalSpriteDef{
			Tag:                raw.Name(int32(fragData.NameRef)),
			CollisionVolumeTag: collisionTag,
		}
		if fragData.Flags&0x01 == 0x01 {
			spriteDef.CenterOffset.Valid = true
			spriteDef.CenterOffset.Float32Slice3 = fragData.CenterOffset
		}
		if fragData.Flags&0x02 == 0x02 {
			spriteDef.BoundingRadius.Valid = true
			spriteDef.BoundingRadius.Float32 = fragData.BoundingRadius
		}

		for _, dag := range fragData.Dags {
			if len(src.Fragments) < int(dag.TrackRef) {
				return fmt.Errorf("track ref %d not found", dag.TrackRef)
			}
			srcTrack, ok := src.Fragments[dag.TrackRef].(*rawfrag.WldFragTrack)
			if !ok {
				return fmt.Errorf("track ref %d not found", dag.TrackRef)
			}

			spriteTag := ""
			if dag.MeshOrSpriteOrParticleRef > 0 {
				if len(src.Fragments) < int(dag.MeshOrSpriteOrParticleRef) {
					return fmt.Errorf("mesh or sprite or particle ref %d not found", dag.MeshOrSpriteOrParticleRef)
				}

				spriteInst, ok := src.Fragments[dag.MeshOrSpriteOrParticleRef].(*rawfrag.WldFragDMSprite)
				if !ok {
					return fmt.Errorf("sprite ref %d not found", dag.MeshOrSpriteOrParticleRef)
				}

				if len(src.Fragments) < int(spriteInst.DMSpriteRef) {
					return fmt.Errorf("dmsprite ref %d not found", spriteInst.DMSpriteRef)
				}

				spriteDef := src.Fragments[spriteInst.DMSpriteRef]
				switch simpleSprite := spriteDef.(type) {
				case *rawfrag.WldFragSimpleSpriteDef:
					spriteTag = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragDMSpriteDef:
					spriteTag = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragHierarchicalSpriteDef:
					spriteTag = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragSprite2D:
					spriteTag = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragDmSpriteDef2:
					spriteTag = raw.Name(simpleSprite.NameRef)
				default:
					return fmt.Errorf("unhandled mesh or sprite or particle reference fragment type %d (%s) at offset %d", spriteDef.FragCode(), raw.FragName(spriteDef.FragCode()), i)
				}
			}
			if spriteTag != "" && collisionTag == "" {
				collisionTag = spriteTag
			}

			dag := Dag{
				Tag:       raw.Name(dag.NameRef),
				Track:     raw.Name(srcTrack.NameRef),
				SubDags:   dag.SubDags,
				SpriteTag: spriteTag,
			}

			spriteDef.Dags = append(spriteDef.Dags, dag)
		}

		// based on fragData.Flags&0x100 == 0x100 {
		for i := 0; i < len(fragData.DMSprites); i++ {
			dmSpriteTag := ""
			if len(src.Fragments) < int(fragData.DMSprites[i]) {
				return fmt.Errorf("dmsprite ref %d not found", fragData.DMSprites[i])
			}
			dmSprite, ok := src.Fragments[fragData.DMSprites[i]].(*rawfrag.WldFragDMSprite)
			if !ok {
				return fmt.Errorf("dmsprite ref %d not found", fragData.DMSprites[i])
			}
			if len(src.Fragments) < int(dmSprite.DMSpriteRef) {
				return fmt.Errorf("dmsprite ref %d not found", dmSprite.DMSpriteRef)
			}
			switch spriteDef := src.Fragments[dmSprite.DMSpriteRef].(type) {
			case *rawfrag.WldFragSimpleSpriteDef:
				dmSpriteTag = raw.Name(spriteDef.NameRef)
			case *rawfrag.WldFragDMSpriteDef:
				dmSpriteTag = raw.Name(spriteDef.NameRef)
			case *rawfrag.WldFragHierarchicalSpriteDef:
				dmSpriteTag = raw.Name(spriteDef.NameRef)
			case *rawfrag.WldFragSprite2D:
				dmSpriteTag = raw.Name(spriteDef.NameRef)
			case *rawfrag.WldFragDmSpriteDef2:
				dmSpriteTag = raw.Name(spriteDef.NameRef)
			default:
				return fmt.Errorf("unhandled dmsprite reference fragment type %d (%s) at offset %d", spriteDef.FragCode(), raw.FragName(spriteDef.FragCode()), i)
			}

			skin := AttachedSkin{
				DMSpriteTag:               dmSpriteTag,
				LinkSkinUpdatesToDagIndex: fragData.LinkSkinUpdatesToDagIndexes[i],
			}

			spriteDef.AttachedSkins = append(spriteDef.AttachedSkins, skin)
		}

		if spriteDef.CollisionVolumeTag != "" {
			if collisionTag == "" {
				return fmt.Errorf("collision volume ref not found")
			}

			for _, attachedSkin := range spriteDef.AttachedSkins {
				isFound := false
				for _, dmSprite := range wld.DMSpriteDef2s {
					if dmSprite.Tag != attachedSkin.DMSpriteTag {
						continue
					}
					dmSprite.PolyhedronTag = collisionTag

					isFound = true
					break
				}
				if !isFound {
					return fmt.Errorf("dmsprite %s not found", attachedSkin.DMSpriteTag)
				}
			}
		}

		wld.HierarchicalSpriteDefs = append(wld.HierarchicalSpriteDefs, spriteDef)
	case rawfrag.FragCodeHierarchicalSprite:
		return nil
	case rawfrag.FragCodeLightDef:
		fragData, ok := fragment.(*rawfrag.WldFragLightDef)
		if !ok {
			return fmt.Errorf("invalid lightdef fragment at offset %d", i)
		}
		light := &LightDef{
			Tag:         raw.Name(fragData.NameRef),
			Flags:       fragData.Flags,
			LightLevels: fragData.LightLevels,
			Colors:      fragData.Colors,
		}
		if fragData.Flags&0x01 == 0x01 {
			light.CurrentFrame.Valid = true
			light.CurrentFrame.Uint32 = fragData.FrameCurrentRef
		}
		if fragData.Flags&0x02 == 0x02 {
			light.Sleep.Valid = true
			light.Sleep.Uint32 = fragData.Sleep
		}
		if fragData.Flags&0x04 == 0x04 {
			light.LightLevels = fragData.LightLevels
		} else {
			if len(fragData.LightLevels) > 0 {
				return fmt.Errorf("light levels found but flag 0x04 not set")
			}
		}

		wld.LightDefs = append(wld.LightDefs, light)
	case rawfrag.FragCodeLight:
		return nil // light instances are ignored, since they're derived from other definitions
	case rawfrag.FragCodeSprite3DDef:
		fragData, ok := fragment.(*rawfrag.WldFragSprite3DDef)
		if !ok {
			return fmt.Errorf("invalid sprite3ddef fragment at offset %d", i)
		}

		if len(src.Fragments) < int(fragData.SphereListRef) {
			return fmt.Errorf("spherelist ref %d out of bounds", fragData.SphereListRef)
		}

		sphereListTag := ""
		if fragData.SphereListRef > 0 {
			sphereList, ok := src.Fragments[fragData.SphereListRef].(*rawfrag.WldFragSphereList)
			if !ok {
				return fmt.Errorf("spherelist ref %d not found", fragData.SphereListRef)
			}
			sphereListTag = raw.Name(sphereList.NameRef)
		}

		sprite := &Sprite3DDef{
			Tag:           raw.Name(fragData.NameRef),
			SphereListTag: sphereListTag,
			Vertices:      fragData.Vertices,
		}

		if fragData.Flags&0x01 == 0x01 {
			sprite.CenterOffset.Valid = true
			sprite.CenterOffset.Float32Slice3 = fragData.CenterOffset
		}

		if fragData.Flags&0x02 == 0x02 {
			sprite.BoundingRadius.Valid = true
			sprite.BoundingRadius.Float32 = fragData.BoundingRadius
		}

		for _, bspNode := range fragData.BspNodes {
			node := &BSPNode{
				FrontTree:    bspNode.FrontTree,
				BackTree:     bspNode.BackTree,
				Vertices:     bspNode.VertexIndexes,
				RenderMethod: model.RenderMethodStr(bspNode.RenderMethod),
			}

			if bspNode.RenderFlags&0x01 == 0x01 {
				node.Pen.Valid = true
				node.Pen.Uint32 = bspNode.RenderPen
			}

			if bspNode.RenderFlags&0x02 == 0x02 {
				node.Brightness.Valid = true
				node.Brightness.Float32 = bspNode.RenderBrightness
			}

			if bspNode.RenderFlags&0x04 == 0x04 {
				node.ScaledAmbient.Valid = true
				node.ScaledAmbient.Float32 = bspNode.RenderScaledAmbient
			}

			if bspNode.RenderFlags&0x08 == 0x08 {
				node.SpriteTag.Valid = true
				if len(src.Fragments) < int(bspNode.RenderSimpleSpriteReference) {
					return fmt.Errorf("sprite ref %d not found", bspNode.RenderSimpleSpriteReference)
				}
				spriteDef := src.Fragments[bspNode.RenderSimpleSpriteReference]
				switch simpleSprite := spriteDef.(type) {
				case *rawfrag.WldFragSimpleSpriteDef:
					node.SpriteTag.String = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragDMSpriteDef:
					node.SpriteTag.String = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragHierarchicalSpriteDef:
					node.SpriteTag.String = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragSprite2D:
					node.SpriteTag.String = raw.Name(simpleSprite.NameRef)
				default:
					return fmt.Errorf("unhandled render sprite reference fragment type %d at offset %d", spriteDef.FragCode(), i)
				}
			}

			if bspNode.RenderFlags&0x10 == 0x10 {
				// has uvinfo
				node.UvOrigin.Valid = true
				node.UAxis.Valid = true
				node.VAxis.Valid = true
				node.UvOrigin.Float32Slice3 = bspNode.RenderUVInfoOrigin
				node.UAxis.Float32Slice3 = bspNode.RenderUVInfoUAxis
				node.VAxis.Float32Slice3 = bspNode.RenderUVInfoVAxis
			}

			if bspNode.RenderFlags&0x20 == 0x20 {
				node.Uvs = bspNode.Uvs
			}

			if bspNode.RenderFlags&0x40 == 0x40 {
				node.TwoSided = 1
			}

			sprite.BSPNodes = append(sprite.BSPNodes, node)
		}

		wld.Sprite3DDefs = append(wld.Sprite3DDefs, sprite)
	case rawfrag.FragCodeSprite3D:
		// sprite instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeZone:
		fragData, ok := fragment.(*rawfrag.WldFragZone)
		if !ok {
			return fmt.Errorf("invalid zone fragment at offset %d", i)
		}

		zone := &Zone{
			Tag:      raw.Name(fragData.NameRef),
			Regions:  fragData.Regions,
			UserData: fragData.UserData,
		}

		wld.Zones = append(wld.Zones, zone)
	case rawfrag.FragCodeWorldTree:
		fragData, ok := fragment.(*rawfrag.WldFragWorldTree)
		if !ok {
			return fmt.Errorf("invalid worldtree fragment at offset %d", i)
		}

		worldTree := &WorldTree{}
		for _, srcNode := range fragData.Nodes {
			node := &WorldNode{
				Normals:        srcNode.Normal,
				WorldRegionTag: raw.Name(srcNode.RegionRef),
				FrontTree:      uint32(srcNode.FrontRef),
				BackTree:       uint32(srcNode.BackRef),
			}
			worldTree.WorldNodes = append(worldTree.WorldNodes, node)

		}
		wld.WorldTrees = append(wld.WorldTrees, worldTree)
	case rawfrag.FragCodeRegion:
		fragData, ok := fragment.(*rawfrag.WldFragRegion)
		if !ok {
			return fmt.Errorf("invalid region fragment at offset %d", i)
		}

		region := &Region{
			VisTree:        &VisTree{},
			Tag:            raw.Name(fragData.NameRef),
			RegionVertices: fragData.RegionVertices,
			Sphere:         fragData.Sphere,
			ReverbVolume:   fragData.ReverbVolume,
			ReverbOffset:   fragData.ReverbOffset,
		}
		// 0x01 is sphere, we just copy
		// 0x02 has reverb volume, we just copy
		// 0x04 has reverb offset, we just copy
		if fragData.Flags&0x08 == 0x08 {
			region.RegionFog = 1
		}
		if fragData.Flags&0x10 == 0x10 {
			region.Gouraud2 = 1
		}
		if fragData.Flags&0x20 == 0x20 {
			region.EncodedVisibility = 1
		}
		// 0x40 unknown
		if fragData.Flags&0x80 == 0x80 {
			region.VisListBytes = 1
		}

		if fragData.MeshReference > 0 && fragData.Flags&0x100 != 0x100 {
			fmt.Printf("region mesh ref %d but flag 0x100 not set\n", fragData.MeshReference)
		}

		if fragData.AmbientLightRef > 0 {
			if len(src.Fragments) < int(fragData.AmbientLightRef) {
				return fmt.Errorf("ambient light ref %d not found", fragData.AmbientLightRef)
			}

			ambientLight, ok := src.Fragments[fragData.AmbientLightRef].(*rawfrag.WldFragGlobalAmbientLightDef)
			if !ok {
				return fmt.Errorf("ambient light ref %d not found", fragData.AmbientLightRef)
			}

			region.AmbientLightTag = raw.Name(ambientLight.NameRef)
		}

		for _, node := range fragData.VisNodes {

			visNode := &VisNode{
				Normal:       node.NormalABCD,
				VisListIndex: node.VisListIndex,
				FrontTree:    node.FrontTree,
				BackTree:     node.BackTree,
			}

			region.VisTree.VisNodes = append(region.VisTree.VisNodes, visNode)
		}

		for _, visList := range fragData.VisLists {
			visListData := &VisList{}
			for _, rangeVal := range visList.Ranges {
				visListData.Ranges = append(visListData.Ranges, int8(rangeVal))
			}

			region.VisTree.VisLists = append(region.VisTree.VisLists, visListData)
		}

		if fragData.MeshReference > 0 {
			if len(src.Fragments) < int(fragData.MeshReference) {
				return fmt.Errorf("mesh ref %d not found", fragData.MeshReference)
			}

			rawMesh := src.Fragments[fragData.MeshReference]
			switch mesh := rawMesh.(type) {
			case *rawfrag.WldFragDmSpriteDef2:
				region.SpriteTag = raw.Name(mesh.NameRef)
			default:
				return fmt.Errorf("unhandled mesh reference fragment type %d (%s) at offset %d", rawMesh.FragCode(), raw.FragName(rawMesh.FragCode()), i)
			}
		}

		wld.Regions = append(wld.Regions, region)
	case rawfrag.FragCodeAmbientLight:
		fragData, ok := fragment.(*rawfrag.WldFragAmbientLight)
		if !ok {
			return fmt.Errorf("invalid ambientlight fragment at offset %d", i)
		}

		lightTag := ""
		lightFlags := uint32(0)
		if fragData.LightRef > 0 {
			if len(src.Fragments) < int(fragData.LightRef) {
				return fmt.Errorf("lightdef ref %d out of bounds", fragData.LightRef)
			}

			light, ok := src.Fragments[fragData.LightRef].(*rawfrag.WldFragLight)
			if !ok {
				return fmt.Errorf("lightdef ref %d not found", fragData.LightRef)
			}

			lightFlags = light.Flags

			lightDef, ok := src.Fragments[light.LightDefRef].(*rawfrag.WldFragLightDef)
			if !ok {
				return fmt.Errorf("lightdef ref %d not found", light.LightDefRef)
			}

			lightTag = raw.Name(lightDef.NameRef)
		}

		light := &AmbientLight{
			Tag:        raw.Name(fragData.NameRef),
			LightTag:   lightTag,
			LightFlags: lightFlags,
			Regions:    fragData.Regions,
		}

		wld.AmbientLights = append(wld.AmbientLights, light)
	case rawfrag.FragCodePointLight:
		fragData, ok := fragment.(*rawfrag.WldFragPointLight)
		if !ok {
			return fmt.Errorf("invalid pointlight fragment at offset %d", i)
		}

		lightTag := ""
		lightFlags := uint32(0)
		if fragData.LightRef > 0 {
			if len(src.Fragments) < int(fragData.LightRef) {
				return fmt.Errorf("light ref %d not found", fragData.LightRef)
			}

			light, ok := src.Fragments[fragData.LightRef].(*rawfrag.WldFragLight)
			if !ok {
				return fmt.Errorf("light ref %d not found", fragData.LightRef)
			}

			lightFlags = light.Flags

			if len(src.Fragments) < int(light.LightDefRef) {
				return fmt.Errorf("lightdef ref %d not found", light.LightDefRef)
			}

			lightDef, ok := src.Fragments[light.LightDefRef].(*rawfrag.WldFragLightDef)
			if !ok {
				return fmt.Errorf("lightdef ref %d not found", light.LightDefRef)
			}

			lightTag = raw.Name(lightDef.NameRef)
		}

		light := &PointLight{
			Tag:        raw.Name(fragData.NameRef),
			LightTag:   lightTag,
			LightFlags: lightFlags,
			Location:   fragData.Location,
			Radius:     fragData.Radius,
		}

		wld.PointLights = append(wld.PointLights, light)
	case rawfrag.FragCodePolyhedronDef:
		fragData, ok := fragment.(*rawfrag.WldFragPolyhedronDef)
		if !ok {
			return fmt.Errorf("invalid polyhedrondef fragment at offset %d", i)
		}

		polyhedron := &PolyhedronDefinition{
			Tag:            raw.Name(fragData.NameRef),
			BoundingRadius: fragData.BoundingRadius,
			ScaleFactor:    fragData.ScaleFactor,
			Vertices:       fragData.Vertices,
		}

		for _, srcFace := range fragData.Faces {
			face := &PolyhedronDefinitionFace{
				Vertices: srcFace.Vertices,
			}

			polyhedron.Faces = append(polyhedron.Faces, face)
		}

		wld.PolyhedronDefs = append(wld.PolyhedronDefs, polyhedron)
	case rawfrag.FragCodePolyhedron:
		// polyhedron instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeSphere:
		// sphere instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeDmRGBTrackDef:
		fragData, ok := fragment.(*rawfrag.WldFragDmRGBTrackDef)
		if !ok {
			return fmt.Errorf("invalid dmrgbtrackdef fragment at offset %d", i)
		}

		track := &RGBTrackDef{
			Tag:   raw.Name(fragData.NameRef),
			Data1: fragData.Data1,
			Data2: fragData.Data2,
			Sleep: fragData.Sleep,
			Data4: fragData.Data4,
			RGBAs: fragData.RGBAs,
		}

		wld.RGBTrackDefs = append(wld.RGBTrackDefs, track)
		return nil
	case rawfrag.FragCodeDmRGBTrack:
		fragData, ok := fragment.(*rawfrag.WldFragDmRGBTrack)
		if !ok {
			return fmt.Errorf("invalid dmrgbtrack fragment at offset %d", i)
		}

		definitionTag := ""
		if fragData.TrackRef > 0 {
			if len(src.Fragments) < int(fragData.TrackRef) {
				return fmt.Errorf("dmrgbtrackdef ref %d not found", fragData.TrackRef)
			}

			trackDef, ok := src.Fragments[fragData.TrackRef].(*rawfrag.WldFragDmRGBTrackDef)
			if !ok {
				return fmt.Errorf("dmrgbtrackdef ref %d not found", fragData.TrackRef)
			}
			definitionTag = raw.Name(trackDef.NameRef)
		}

		track := &RGBTrack{
			Tag:           raw.Name(fragData.NameRef),
			DefinitionTag: definitionTag,
			Flags:         fragData.Flags,
		}

		wld.RGBTrackInsts = append(wld.RGBTrackInsts, track)
		return nil
	default:
		return fmt.Errorf("unhandled fragment type %d (%s)", fragment.FragCode(), raw.FragName(fragment.FragCode()))
	}

	return nil
}

func (wld *Wld) WriteRaw(w io.Writer) error {
	var err error
	dst := &raw.Wld{
		IsOldWorld: true,
	}
	if dst.Fragments == nil {
		dst.Fragments = []model.FragmentReadWriter{}
	}
	raw.NameClear()

	if wld.GlobalAmbientLight != nil {
		wld.isZone = true
		_, err = wld.GlobalAmbientLight.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("global ambient light: %w", err)
		}

	}

	for _, dmSprite := range wld.DMSpriteDef2s {
		_, err = dmSprite.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("dmspritedef2 %s: %w", dmSprite.Tag, err)
		}
	}
	for _, hiSprite := range wld.HierarchicalSpriteDefs {
		_, err = hiSprite.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("hierarchicalsprite %s: %w", hiSprite.Tag, err)
		}
	}

	for _, lightDef := range wld.LightDefs {
		_, err = lightDef.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("light %s: %w", lightDef.Tag, err)
		}

	}

	for _, sprite := range wld.Sprite3DDefs {
		_, err = sprite.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
		}
	}

	for _, tree := range wld.WorldTrees {
		_, err = tree.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("worldtree: %w", err)
		}
	}

	for _, region := range wld.Regions {
		_, err = region.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("region %s: %w", region.Tag, err)
		}
	}

	for _, alight := range wld.AmbientLights {
		_, err = alight.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("ambientlight %s: %w", alight.Tag, err)
		}
	}

	for _, actor := range wld.ActorInsts {
		_, err = actor.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("actor %s: %w", actor.Tag, err)
		}
	}

	for _, track := range wld.TrackInstances {
		if track.fragID > 0 {
			continue
		}

		_, err = track.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("track %s: %w", track.Tag, err)
		}

	}

	for _, actorDef := range wld.ActorDefs {
		if actorDef.fragID > 0 {
			continue
		}

		_, err = actorDef.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("actordef %s: %w", actorDef.Tag, err)
		}
	}

	for _, zone := range wld.Zones {
		_, err = zone.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("zone %s: %w", zone.Tag, err)
		}
	}

	return dst.Write(w)
}

var animationPrefixesMap = map[string]struct{}{
	"C01": {}, "C02": {}, "C03": {}, "C04": {}, "C05": {}, "C06": {}, "C07": {}, "C08": {}, "C09": {}, "C10": {}, "C11": {},
	"D01": {}, "D02": {}, "D03": {}, "D04": {}, "D05": {},
	"L01": {}, "L02": {}, "L03": {}, "L04": {}, "L05": {}, "L06": {}, "L07": {}, "L08": {}, "L09": {},
	"O01": {},
	"S01": {}, "S02": {}, "S03": {}, "S04": {}, "S05": {}, "S06": {}, "S07": {}, "S08": {}, "S09": {}, "S10": {},
	"S11": {}, "S12": {}, "S13": {}, "S14": {}, "S15": {}, "S16": {}, "S17": {}, "S18": {}, "S19": {}, "S20": {},
	"S21": {}, "S22": {}, "S23": {}, "S24": {}, "S25": {}, "S26": {}, "S27": {}, "S28": {},
	"P01": {}, "P02": {}, "P03": {}, "P04": {}, "P05": {}, "P06": {}, "P07": {}, "P08": {},
	"O02": {}, "O03": {},
	"T01": {}, "T02": {}, "T03": {}, "T04": {}, "T05": {}, "T06": {}, "T07": {}, "T08": {}, "T09": {},
}

func isAnimationPrefix(name string) bool {
	if len(name) < 3 {
		return false
	}
	prefix := name[:3]

	_, exists := animationPrefixesMap[prefix]
	return exists
}
