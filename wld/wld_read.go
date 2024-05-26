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
