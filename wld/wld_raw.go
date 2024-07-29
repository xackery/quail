package wld

import (
	"fmt"
	"io"

	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

func (wld *Wld) ReadRaw(src *raw.Wld) error {

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
		tag := raw.Name(fragData.NameRef)
		if len(tag) == 0 {
			if fragData.NameRef == 0xFF0000 {
				tag = "GLOBALAMBIENT_LIGHTDEF"
			}
		}
		if wld.GlobalAmbientLight != "" {
			return fmt.Errorf("duplicate globalambientlightdef found")
		}
		wld.GlobalAmbientLight = tag
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
		}
		material := &MaterialDef{
			Tag:                  raw.Name(fragData.NameRef),
			Flags:                fragData.Flags,
			RGBPen:               fragData.RGBPen,
			Brightness:           fragData.Brightness,
			ScaledAmbient:        fragData.ScaledAmbient,
			Pair1:                fragData.Pair1,
			Pair2:                fragData.Pair2,
			SimpleSpriteInstTag:  spriteTag,
			SimpleSpriteInstFlag: spriteFlags,
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
			Colors:               fragData.Colors,
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
		for _, color := range fragData.Colors {
			sprite.Colors = append(sprite.Colors, [4]uint8{
				color[0],
				color[1],
				color[2],
				color[3],
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
			Tag:   raw.Name(fragData.NameRef),
			Flags: fragData.Flags,
		}

		for _, transform := range fragData.BoneTransforms {
			frame := TrackFrameTransform{
				PositionDenom: float32(transform.ShiftDenominator),
				Rotation:      transform.Rotation,
			}
			//TODO: fix scale

			//scale := float32(1 / float32(int(1)<<int(frame.PositionDenom)))
			scale := float32(1)
			frame.Position[0] = float32(transform.Shift[0]) / scale
			frame.Position[1] = float32(transform.Shift[1]) / scale
			frame.Position[2] = float32(transform.Shift[2]) / scale

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
			trackInst.Sleep = fragData.Sleep
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
			Tag:           raw.Name(fragData.NameRef),
			Callback:      raw.Name(fragData.CallbackNameRef),
			BoundsRef:     fragData.BoundsRef,
			CurrentAction: fragData.CurrentAction,
			Location:      fragData.Location,
			Unk1:          fragData.Unk1,
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
		if fragData.ActorDefNameRef > 0 {
			if len(src.Fragments) < int(fragData.ActorDefNameRef) {
				return fmt.Errorf("actordef ref %d out of bounds", fragData.ActorDefNameRef)
			}

			actorDef, ok := src.Fragments[fragData.ActorDefNameRef].(*rawfrag.WldFragActorDef)
			if !ok {
				return fmt.Errorf("actordef ref %d not found", fragData.ActorDefNameRef)
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

		trackTag := ""
		if fragData.DMRGBTrackRef > 0 {
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
		}

		actor := &ActorInst{
			Tag:            raw.Name(fragData.NameRef),
			DefinitionTag:  actorDefTag,
			Flags:          fragData.Flags,
			SphereRadius:   sphereRadius,
			DMRGBTrackTag:  trackTag,
			CurrentAction:  fragData.CurrentAction,
			Location:       fragData.Location,
			Unk1:           fragData.Unk1,
			BoundingRadius: fragData.BoundingRadius,
			Scale:          fragData.ScaleFactor,
			SoundTag:       raw.Name(fragData.SoundNameRef),
			UserData:       fragData.UserData,
		}

		wld.ActorInsts = append(wld.ActorInsts, actor)
	case rawfrag.FragCodeHierarchialSpriteDef:
		fragData, ok := fragment.(*rawfrag.WldFragHierarchialSpriteDef)
		if !ok {
			return fmt.Errorf("invalid hierarchialsprite fragment at offset %d", i)
		}

		polyhedronTag := ""
		if fragData.CollisionVolumeRef != 0 && fragData.CollisionVolumeRef != 4294967293 {
			if len(src.Fragments) < int(fragData.CollisionVolumeRef) {
				return fmt.Errorf("collision volume ref %d out of bounds", fragData.CollisionVolumeRef)
			}

			switch collision := src.Fragments[fragData.CollisionVolumeRef].(type) {
			case *rawfrag.WldFragPolyhedron:
				polyhedronTag = raw.Name(collision.NameRef)
			default:
				return fmt.Errorf("unknown collision volume ref %d (%s)", fragData.CollisionVolumeRef, raw.FragName(collision.FragCode()))
			}
		}
		if polyhedronTag != "" {
			return fmt.Errorf("collision volume ref found as %s, report this to xack", polyhedronTag)
		}

		sprite := &HierarchicalSpriteDef{
			Tag:            raw.Name(int32(fragData.NameRef)),
			BoundingRadius: fragData.BoundingRadius,
			CenterOffset:   fragData.CenterOffset,
		}
		if fragData.CollisionVolumeRef == 0xfffffffd {
			sprite.DagCollision = 1
		}
		collisionTag := ""

		for _, srcBone := range fragData.Dags {
			if len(src.Fragments) < int(srcBone.TrackRef) {
				return fmt.Errorf("track ref %d not found", srcBone.TrackRef)
			}
			srcTrack, ok := src.Fragments[srcBone.TrackRef].(*rawfrag.WldFragTrack)
			if !ok {
				return fmt.Errorf("track ref %d not found", srcBone.TrackRef)
			}

			spriteTag := ""
			if srcBone.MeshOrSpriteOrParticleRef > 0 {
				if len(src.Fragments) < int(srcBone.MeshOrSpriteOrParticleRef) {
					return fmt.Errorf("mesh or sprite or particle ref %d not found", srcBone.MeshOrSpriteOrParticleRef)
				}

				spriteInst, ok := src.Fragments[srcBone.MeshOrSpriteOrParticleRef].(*rawfrag.WldFragDMSprite)
				if !ok {
					return fmt.Errorf("sprite ref %d not found", srcBone.MeshOrSpriteOrParticleRef)
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
				case *rawfrag.WldFragHierarchialSpriteDef:
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
				Tag:       raw.Name(srcBone.NameRef),
				Track:     raw.Name(srcTrack.NameRef),
				SubDags:   srcBone.SubBones,
				SpriteTag: spriteTag,
			}

			sprite.Dags = append(sprite.Dags, dag)
		}

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
			case *rawfrag.WldFragHierarchialSpriteDef:
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

			sprite.AttachedSkins = append(sprite.AttachedSkins, skin)
		}

		if sprite.DagCollision == 1 {
			if collisionTag == "" {
				return fmt.Errorf("collision volume ref not found")
			}

			for _, attachedSkin := range sprite.AttachedSkins {
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

		wld.HierarchicalSpriteDefs = append(wld.HierarchicalSpriteDefs, sprite)
	case rawfrag.FragCodeHierarchialSprite:
		return nil
	case rawfrag.FragCodeLightDef:
		fragData, ok := fragment.(*rawfrag.WldFragLightDef)
		if !ok {
			return fmt.Errorf("invalid lightdef fragment at offset %d", i)
		}
		light := &LightDef{
			Tag:             raw.Name(fragData.NameRef),
			Sleep:           fragData.Sleep,
			FrameCurrentRef: fragData.FrameCurrentRef,
			LightLevels:     fragData.LightLevels,
			Colors:          fragData.Colors,
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
			Tag:            raw.Name(fragData.NameRef),
			SphereListTag:  sphereListTag,
			CenterOffset:   fragData.CenterOffset,
			BoundingRadius: fragData.BoundingRadius,
			Vertices:       fragData.Vertices,
		}

		for _, bspNode := range fragData.BspNodes {
			node := &BSPNode{
				FrontTree:     bspNode.FrontTree,
				BackTree:      bspNode.BackTree,
				Vertices:      bspNode.VertexIndexes,
				RenderMethod:  model.RenderMethodStr(bspNode.RenderMethod),
				Flags:         bspNode.RenderFlags,
				Pen:           bspNode.RenderPen,
				Brightness:    bspNode.RenderBrightness,
				ScaledAmbient: bspNode.RenderScaledAmbient,
				Origin:        bspNode.RenderUVInfoOrigin,
				UAxis:         bspNode.RenderUVInfoUAxis,
				VAxis:         bspNode.RenderUVInfoVAxis,
			}

			if bspNode.RenderFlags&0x03 == 0x03 {
				if len(src.Fragments) < int(bspNode.RenderSimpleSpriteReference) {
					return fmt.Errorf("sprite ref %d not found", bspNode.RenderSimpleSpriteReference)
				}
				spriteDef := src.Fragments[bspNode.RenderSimpleSpriteReference]
				switch simpleSprite := spriteDef.(type) {
				case *rawfrag.WldFragSimpleSpriteDef:
					node.SpriteTag = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragDMSpriteDef:
					node.SpriteTag = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragHierarchialSpriteDef:
					node.SpriteTag = raw.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragSprite2D:
					node.SpriteTag = raw.Name(simpleSprite.NameRef)
				default:
					return fmt.Errorf("unhandled render sprite reference fragment type %d at offset %d", spriteDef.FragCode(), i)
				}
			}

			for _, uvMap := range bspNode.RenderUVMapEntries {
				entry := BspNodeUVInfo{
					UvOrigin: uvMap.UvOrigin,
					UAxis:    uvMap.UAxis,
					VAxis:    uvMap.VAxis,
				}
				node.RenderUVMapEntries = append(node.RenderUVMapEntries, entry)
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
			RegionTag:      raw.Name(fragData.NameRef),
			RegionVertices: fragData.RegionVertices,
			Sphere:         fragData.Sphere,
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
	dst := &raw.Wld{}
	if dst.Fragments == nil {
		dst.Fragments = []model.FragmentReadWriter{}
	}
	raw.NameClear()

	globalAmbient := &rawfrag.WldFragGlobalAmbientLightDef{
		//NameRef: raw.NameAdd(wld.GlobalAmbientLight),
	}
	dst.Fragments = append(dst.Fragments, globalAmbient)

	for _, dmSprite := range wld.DMSpriteDef2s {
		isMaterialPaletteFound := false
		for _, materialPalette := range wld.MaterialPalettes {
			if materialPalette.Tag != dmSprite.MaterialPaletteTag {
				continue
			}
			isMaterialPaletteFound = true
			for _, material := range materialPalette.Materials {
				isMaterialDefFound := false
				for _, materialDef := range wld.MaterialDefs {
					if materialDef.Tag != material {
						continue
					}
					isMaterialDefFound = true
					if materialDef.SimpleSpriteInstTag != "" {
						isMaterialFound := false

						for _, sprite := range wld.SimpleSpriteDefs {
							if sprite.Tag != materialDef.SimpleSpriteInstTag {
								continue
							}

							spriteDef := &rawfrag.WldFragSimpleSpriteDef{}
							for _, frame := range sprite.SimpleSpriteFrames {
								if frame.TextureTag == "" {
									continue
								}
								bmInfo := &rawfrag.WldFragBMInfo{
									NameRef:      raw.NameAdd(frame.TextureTag),
									TextureNames: []string{frame.TextureFile},
								}
								dst.Fragments = append(dst.Fragments, bmInfo)
								spriteDef.BitmapRefs = append(spriteDef.BitmapRefs, uint32(len(dst.Fragments)))
							}

							spriteDef.NameRef = raw.NameAdd(sprite.Tag)

							dst.Fragments = append(dst.Fragments, spriteDef)

							inst := &rawfrag.WldFragSimpleSprite{
								//NameRef:   raw.NameAdd(sprite.Tag),
								SpriteRef: int16(len(dst.Fragments)),
								//Flags:    sprite.Flags,
							}
							dst.Fragments = append(dst.Fragments, inst)
							isMaterialFound = true
							break
						}
						if !isMaterialFound {
							return fmt.Errorf("simple sprite %s not found", materialDef.SimpleSpriteInstTag)
						}
					}

					materialDef := &rawfrag.WldFragMaterialDef{
						NameRef:       raw.NameAdd(materialDef.Tag),
						Flags:         materialDef.Flags,
						RenderMethod:  model.RenderMethodInt(materialDef.RenderMethod),
						RGBPen:        materialDef.RGBPen,
						Brightness:    materialDef.Brightness,
						ScaledAmbient: materialDef.ScaledAmbient,
						Pair1:         materialDef.Pair1,
						Pair2:         materialDef.Pair2,
					}
					dst.Fragments = append(dst.Fragments, materialDef)
				}
				if !isMaterialDefFound {
					return fmt.Errorf("materialdef %s not found", material)
				}
			}
			materialPaletteDef := &rawfrag.WldFragMaterialPalette{
				NameRef: raw.NameAdd(materialPalette.Tag),
				Flags:   materialPalette.flags,
			}
			for _, material := range materialPalette.Materials {
				materialPaletteDef.MaterialRefs = append(materialPaletteDef.MaterialRefs, uint32(raw.NameAdd(material)))
			}
			dst.Fragments = append(dst.Fragments, materialPaletteDef)
			break
		}
		if !isMaterialPaletteFound {
			return fmt.Errorf("material %s not found", dmSprite.MaterialPaletteTag)
		}
		dmSpriteDef := &rawfrag.WldFragDmSpriteDef2{
			NameRef:            raw.NameAdd(dmSprite.Tag),
			Flags:              dmSprite.Flags,
			MaterialPaletteRef: uint32(len(dst.Fragments)),
			CenterOffset:       dmSprite.CenterOffset,
			Params2:            dmSprite.Params2,
			MaxDistance:        dmSprite.MaxDistance,
			Min:                dmSprite.Min,
			Max:                dmSprite.Max,
			Scale:              dmSprite.FPScale,
			//Vertices: 		  dmSprite.Vertices,
			Colors: dmSprite.Colors,
		}

		scale := float32(1 / float32(int(1)<<int(dmSprite.FPScale)))

		for _, vert := range dmSprite.Vertices {
			dmSpriteDef.Vertices = append(dmSpriteDef.Vertices, [3]int16{
				int16(vert[0] / scale),
				int16(vert[1] / scale),
				int16(vert[2] / scale),
			})
		}

		for _, uv := range dmSprite.UVs {
			dmSpriteDef.UVs = append(dmSpriteDef.UVs, [2]int16{
				int16(uv[0] / scale),
				int16(uv[1] / scale),
			})
		}

		for _, normal := range dmSprite.VertexNormals {
			dmSpriteDef.VertexNormals = append(dmSpriteDef.VertexNormals, [3]int8{
				int8(normal[0] / scale),
				int8(normal[1] / scale),
				int8(normal[2] / scale),
			})
		}

		dst.Fragments = append(dst.Fragments, dmSpriteDef)
	}

	for _, lightDef := range wld.LightDefs {
		light := &rawfrag.WldFragLightDef{
			NameRef:         raw.NameAdd(lightDef.Tag),
			Sleep:           lightDef.Sleep,
			FrameCurrentRef: lightDef.FrameCurrentRef,
			LightLevels:     lightDef.LightLevels,
			Colors:          lightDef.Colors,
		}
		dst.Fragments = append(dst.Fragments, light)
	}

	for _, sprite := range wld.Sprite3DDefs {
		spriteDef := &rawfrag.WldFragSprite3DDef{
			NameRef:        raw.NameAdd(sprite.Tag),
			CenterOffset:   sprite.CenterOffset,
			BoundingRadius: sprite.BoundingRadius,
			Vertices:       sprite.Vertices,
		}

		for _, node := range sprite.BSPNodes {
			bnode := rawfrag.WldFragThreeDSpriteBspNode{
				FrontTree: node.FrontTree,
				BackTree:  node.BackTree,

				RenderMethod:        model.RenderMethodInt(node.RenderMethod),
				RenderFlags:         node.Flags,
				RenderPen:           node.Pen,
				RenderBrightness:    node.Brightness,
				RenderScaledAmbient: node.ScaledAmbient,
				RenderUVInfoOrigin:  node.Origin,
				RenderUVInfoUAxis:   node.UAxis,
				RenderUVInfoVAxis:   node.VAxis,
			}

			spriteDef.BspNodes = append(spriteDef.BspNodes, bnode)
		}

		dst.Fragments = append(dst.Fragments, spriteDef)
	}

	for _, tree := range wld.WorldTrees {
		worldTree := &rawfrag.WldFragWorldTree{
			Nodes: []rawfrag.WorldTreeNode{},
		}
		for _, node := range tree.WorldNodes {
			worldNode := rawfrag.WorldTreeNode{
				Normal:    node.Normals,
				RegionRef: raw.NameAdd(node.WorldRegionTag),
				FrontRef:  int32(node.FrontTree),
				BackRef:   int32(node.BackTree),
			}
			worldTree.Nodes = append(worldTree.Nodes, worldNode)
		}
		dst.Fragments = append(dst.Fragments, worldTree)
	}

	for _, region := range wld.Regions {
		regionDef := &rawfrag.WldFragRegion{
			NameRef:        raw.NameAdd(region.RegionTag),
			RegionVertices: region.RegionVertices,
			Sphere:         region.Sphere,
			VisNodes:       []rawfrag.VisNode{},
			VisLists:       []rawfrag.VisList{},
		}

		if region.AmbientLightTag != "" {
			isFound := false
			for _, light := range wld.AmbientLights {
				if light.Tag != region.AmbientLightTag {
					continue
				}
				isFound = true

				regionDef.AmbientLightRef = int32(len(dst.Fragments))

				ambientLight := &rawfrag.WldFragAmbientLight{
					NameRef:  raw.NameAdd(light.Tag),
					LightRef: int32(len(dst.Fragments) + 1),
				}
				dst.Fragments = append(dst.Fragments, ambientLight)

				lightDef := &rawfrag.WldFragLight{
					LightDefRef: int32(len(dst.Fragments) + 1),
					Flags:       light.LightFlags,
				}
				dst.Fragments = append(dst.Fragments, lightDef)

				lightDefRef := &rawfrag.WldFragLightDef{

					NameRef: raw.NameAdd(light.LightTag),
				}
				dst.Fragments = append(dst.Fragments, lightDefRef)
				break
			}
			if !isFound {
				return fmt.Errorf("ambient light %s not found", region.AmbientLightTag)
			}
		}

		for _, node := range region.VisTree.VisNodes {
			visNode := rawfrag.VisNode{
				NormalABCD:   node.Normal,
				VisListIndex: node.VisListIndex,
				FrontTree:    node.FrontTree,
				BackTree:     node.BackTree,
			}
			regionDef.VisNodes = append(regionDef.VisNodes, visNode)
		}

		//for _, visList := range region.VisTree.VisLists {
		//visListData := rawfrag.VisList{
		//		Ranges: []byte{},
		//	}
		//for _, rangeVal := range visList.Ranges {
		//	visListData.Ranges = append(visListData.Ranges, byte(rangeVal))
		//}
		//	regionDef.VisLists = append(regionDef.VisLists, visListData)
		//}

		dst.Fragments = append(dst.Fragments, regionDef)
	}

	for _, actor := range wld.ActorInsts {

		if actor.DefinitionTag != "" {
			isFound := false
			for _, actorDef := range wld.ActorDefs {
				if actorDef.Tag != actor.DefinitionTag {
					continue
				}
				isFound = true
				actorDef := &rawfrag.WldFragActorDef{
					NameRef: raw.NameAdd(actorDef.Tag),
					//Flags:  actorDef.Flags,
					CallbackNameRef: raw.NameAdd(actorDef.Callback),
					BoundsRef:       actorDef.BoundsRef,
					CurrentAction:   actorDef.CurrentAction,
					Location:        actorDef.Location,
					Unk1:            actorDef.Unk1,
				}
				dst.Fragments = append(dst.Fragments, actorDef)
				break
			}
			if !isFound {
				return fmt.Errorf("actordef %s not found", actor.DefinitionTag)
			}
		}

		trackRef := int32(0)
		if actor.DMRGBTrackTag != "" {
			isFound := false
			for _, track := range wld.RGBTrackDefs {
				if track.Tag != actor.DMRGBTrackTag {
					continue
				}
				isFound = true
				trackRef = int32(len(dst.Fragments))
				trackDef := &rawfrag.WldFragDmRGBTrackDef{
					NameRef: raw.NameAdd(track.Tag),
					Data1:   track.Data1,
					Data2:   track.Data2,
					Sleep:   track.Sleep,
					Data4:   track.Data4,

					RGBAs: track.RGBAs,
				}
				dst.Fragments = append(dst.Fragments, trackDef)
				break
			}
			if !isFound {
				return fmt.Errorf("dmrgbtrackdef %s not found", actor.DMRGBTrackTag)
			}
		}

		actorInst := &rawfrag.WldFragActor{
			NameRef:         raw.NameAdd(actor.Tag),
			ActorDefNameRef: raw.NameAdd(actor.DefinitionTag),
			DMRGBTrackRef:   trackRef,
			Flags:           actor.Flags,
		}

		if actor.SphereRadius != 0 {

			sphereDef := &rawfrag.WldFragSphere{
				NameRef: raw.NameAdd(actor.Tag),
				Radius:  actor.SphereRadius,
			}
			dst.Fragments = append(dst.Fragments, sphereDef)
		}
		actorInst.SphereRef = uint32(len(dst.Fragments))

		dst.Fragments = append(dst.Fragments, actorInst)
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
	_, exists := animationPrefixesMap[name]
	return exists
}
