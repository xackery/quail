package wld

import (
	"fmt"

	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wld/cache"
)

func (wld *Wld) Read(src *raw.Wld) error {
	cm := &cache.CacheManager{}
	err := cm.Load(src)
	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	defer cm.Close()

	wld.Version = src.Version
	wld.GlobalAmbientLight = cm.GlobalAmbientLight
	wld.FileName = src.MetaFileName

	err = wld.readDMSpriteDef2(cm)
	if err != nil {
		return fmt.Errorf("readDMSpriteDef2: %w", err)
	}

	err = wld.readMaterialDef(cm)
	if err != nil {
		return fmt.Errorf("readMaterialDef: %w", err)
	}

	err = wld.readMaterialPalette(cm)
	if err != nil {
		return fmt.Errorf("readMaterialPalette: %w", err)
	}

	err = wld.readSimpleSpriteDef(cm)
	if err != nil {
		return fmt.Errorf("readSimpleSpriteDef: %w", err)
	}

	err = wld.readActorDef(cm)
	if err != nil {
		return fmt.Errorf("readActorDef: %w", err)
	}

	err = wld.readActorInst(cm)
	if err != nil {
		return fmt.Errorf("readActorInst: %w", err)
	}

	err = wld.readLightDef(cm)
	if err != nil {
		return fmt.Errorf("readLightDef: %w", err)
	}

	err = wld.readPointLight(cm)
	if err != nil {
		return fmt.Errorf("readPointLight: %w", err)
	}

	err = wld.readSprite3DDef(cm)
	if err != nil {
		return fmt.Errorf("readSprite3DDef: %w", err)
	}

	return nil
}

func (wld *Wld) readDMSpriteDef2(cm *cache.CacheManager) error {
	for _, src := range cm.DmSpriteDef2s {
		scale := float32(1 / float32(int(1)<<int(src.Scale)))

		dst := &DMSpriteDef2{
			Tag:                  src.Tag,
			Flags:                src.Flags,
			MaterialPaletteTag:   src.MaterialPaletteTag,
			DmTrackTag:           src.DmTrackTag,
			Fragment3Ref:         src.Fragment3Ref,
			Fragment4Ref:         src.Fragment4Ref,
			CenterOffset:         src.CenterOffset,
			Params2:              src.Params2,
			MaxDistance:          src.MaxDistance,
			Min:                  src.Min,
			Max:                  src.Max,
			FPScale:              src.Scale,
			SkinAssignmentGroups: src.SkinAssignmentGroups,
			FaceMaterialGroups:   src.FaceMaterialGroups,
			VertexMaterialGroups: src.VertexMaterialGroups,
		}

		for _, vert := range src.Vertices {
			dst.Vertices = append(dst.Vertices, [3]float32{
				float32(vert[0]) * scale,
				float32(vert[1]) * scale,
				float32(vert[2]) * scale,
			})
		}
		for _, uv := range src.UVs {
			dst.UVs = append(dst.UVs, [2]float32{
				float32(uv[0]) * scale,
				float32(uv[1]) * scale,
			})
		}
		for _, vn := range src.VertexNormals {
			dst.VertexNormals = append(dst.VertexNormals, [3]float32{
				float32(vn[0]) * scale,
				float32(vn[1]) * scale,
				float32(vn[2]) * scale,
			})
		}

		for _, color := range src.Colors {
			dst.Colors = append(dst.Colors, [4]uint8{
				color[0],
				color[1],
				color[2],
				color[3],
			})
		}

		for _, face := range src.Faces {
			dst.Faces = append(dst.Faces, &Face{
				Flags:    face.Flags,
				Triangle: [3]uint16{face.Index[0], face.Index[1], face.Index[2]},
			})
		}

		for _, mop := range src.MeshOps {
			dst.MeshOps = append(dst.MeshOps, &MeshOp{
				Index1:    mop.Index1,
				Index2:    mop.Index2,
				Offset:    mop.Offset,
				Param1:    mop.Param1,
				TypeField: mop.TypeField,
			})
		}

		wld.DMSpriteDef2s = append(wld.DMSpriteDef2s, dst)
	}
	return nil

}

func (wld *Wld) readMaterialDef(cm *cache.CacheManager) error {
	for _, src := range cm.MaterialDefs {

		dst := &MaterialDef{
			Tag:                  src.Tag,
			Flags:                src.Flags,
			RenderMethod:         src.RenderMethod,
			RGBPen:               src.RGBPen,
			Brightness:           src.Brightness,
			ScaledAmbient:        src.ScaledAmbient,
			SimpleSpriteInstTag:  src.SimpleSpriteInstTag,
			SimpleSpriteInstFlag: src.SimpleSpriteInstFlag,
		}

		wld.MaterialDefs = append(wld.MaterialDefs, dst)
	}
	return nil
}

func (wld *Wld) readMaterialPalette(cm *cache.CacheManager) error {
	for _, src := range cm.MaterialPalettes {
		dst := &MaterialPalette{
			Tag:       src.Tag,
			Flags:     src.Flags,
			Materials: src.Materials,
		}
		wld.MaterialPalettes = append(wld.MaterialPalettes, dst)
	}
	return nil
}

func (wld *Wld) readSimpleSpriteDef(cm *cache.CacheManager) error {
	for _, src := range cm.SimpleSpriteDefs {
		dst := &SimpleSpriteDef{
			Tag:     src.Tag,
			BMInfos: src.BMInfos,
		}

		wld.SimpleSpriteDefs = append(wld.SimpleSpriteDefs, dst)
	}
	return nil
}

func (wld *Wld) readActorDef(cm *cache.CacheManager) error {
	for _, src := range cm.ActorDefs {
		dst := &ActorDef{
			Tag:           src.Tag,
			Callback:      src.Callback,
			BoundsRef:     src.BoundsRef,
			CurrentAction: src.CurrentAction,
			Location:      src.Location,
			Unk1:          src.Unk1,
			FragmentRefs:  src.FragmentRefs,
		}
		for _, srcAction := range src.Actions {
			dstAction := ActorAction{}

			for i := 0; i < len(srcAction.Lods); i++ {
				dstLod := ActorLevelOfDetail{
					MinDistance: srcAction.Lods[i],
				}

				dstAction.LevelOfDetails = append(dstAction.LevelOfDetails, dstLod)
			}
			dst.Actions = append(dst.Actions, dstAction)
		}

		wld.ActorDefs = append(wld.ActorDefs, dst)
	}
	return nil
}

func (wld *Wld) readActorInst(cm *cache.CacheManager) error {
	for _, src := range cm.ActorInsts {
		dst := &ActorInst{
			Tag:            src.Tag,
			DefinitionTag:  src.ActorDefTag,
			Flags:          src.Flags,
			SphereTag:      src.SphereTag,
			CurrentAction:  src.CurrentAction,
			Location:       src.Location,
			Unk1:           src.Unk1,
			BoundingRadius: src.BoundingRadius,
			Scale:          src.Scale,
			Unk2:           src.Unk2,
		}

		wld.ActorInsts = append(wld.ActorInsts, dst)
	}
	return nil
}

func (wld *Wld) readLightDef(cm *cache.CacheManager) error {
	for _, src := range cm.LightDefs {
		dst := &LightDef{
			Tag:             src.Tag,
			Flags:           src.Flags,
			FrameCurrentRef: src.FrameCurrentRef,
			LightLevels:     src.LightLevels,
			Colors:          src.Colors,
		}

		wld.LightDefs = append(wld.LightDefs, dst)
	}
	return nil
}

func (wld *Wld) readPointLight(cm *cache.CacheManager) error {
	for _, src := range cm.PointLights {
		dst := &PointLight{
			Tag:         src.Tag,
			LightDefTag: src.LightDefTag,
			Flags:       src.Flags,
			Location:    src.Location,
			Radius:      src.Radius,
		}

		wld.PointLights = append(wld.PointLights, dst)
	}
	return nil
}

func (wld *Wld) readSprite3DDef(cm *cache.CacheManager) error {
	for _, src := range cm.Sprite3DDefs {
		dst := &Sprite3DDef{
			Tag:      src.Tag,
			Vertices: src.Vertices,
		}
		for _, srcBspNode := range src.BspNodes {
			dstBspNode := &BSPNode{
				Vertices:                    srcBspNode.VertexIndexes,
				RenderMethod:                srcBspNode.RenderMethod,
				RenderFlags:                 srcBspNode.RenderFlags,
				RenderPen:                   srcBspNode.RenderPen,
				RenderBrightness:            srcBspNode.RenderBrightness,
				RenderScaledAmbient:         srcBspNode.RenderScaledAmbient,
				RenderSimpleSpriteReference: srcBspNode.RenderSimpleSpriteReference,
				RenderUVInfoOrigin:          srcBspNode.RenderUVInfoOrigin,
				RenderUVInfoUAxis:           srcBspNode.RenderUVInfoUAxis,
				RenderUVInfoVAxis:           srcBspNode.RenderUVInfoVAxis,
				FrontTree:                   srcBspNode.FrontTree,
				BackTree:                    srcBspNode.BackTree,
			}

			for _, srcUVMapMethod := range srcBspNode.RenderUVMapEntries {
				dstUVMapMethod := BspNodeUVInfo{
					UvOrigin: srcUVMapMethod.UvOrigin,
					UAxis:    srcUVMapMethod.UAxis,
					VAxis:    srcUVMapMethod.VAxis,
				}

				dstBspNode.RenderUVMapEntries = append(dstBspNode.RenderUVMapEntries, dstUVMapMethod)
			}

			dst.BSPNodes = append(dst.BSPNodes, dstBspNode)
		}

		wld.Sprite3DDefs = append(wld.Sprite3DDefs, dst)
	}
	return nil
}
