package wld

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

func (wld *Wld) ReadRaw(src *raw.Wld) error {
	wld.reset()
	wld.WorldDef = &WorldDef{}
	if src.IsNewWorld {
		wld.WorldDef.NewWorld = 1
	}
	if src.IsZone {
		wld.WorldDef.Zone = 1
	}

	for i := 1; i < len(src.Fragments); i++ {
		fragment := src.Fragments[i]
		err := readRawFrag(wld, src, fragment)
		if err != nil {
			return fmt.Errorf("fragment %d (%s): %w", i, raw.FragName(fragment.FragCode()), err)
		}
	}

	return nil
}

func readRawFrag(wld *Wld, rawWld *raw.Wld, fragment model.FragmentReadWriter) error {

	switch fragment.FragCode() {
	case rawfrag.FragCodeGlobalAmbientLightDef:

		def := &GlobalAmbientLightDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragGlobalAmbientLightDef))
		if err != nil {
			return fmt.Errorf("globalambientlightdef: %w", err)
		}
		wld.GlobalAmbientLightDef = def
	case rawfrag.FragCodeBMInfo:
		return nil
	case rawfrag.FragCodeSimpleSpriteDef:
		def := &SimpleSpriteDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragSimpleSpriteDef))
		if err != nil {
			return fmt.Errorf("simplespritedef: %w", err)
		}

		wld.SimpleSpriteDefs = append(wld.SimpleSpriteDefs, def)
	case rawfrag.FragCodeSimpleSprite:
		//return fmt.Errorf("simplesprite fragment found, but not expected")
	case rawfrag.FragCodeBlitSpriteDef:
		def := &BlitSpriteDefinition{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragBlitSpriteDef))
		if err != nil {
			return fmt.Errorf("blitspritedef: %w", err)
		}
		wld.BlitSpriteDefinitions = append(wld.BlitSpriteDefinitions, def)
	case rawfrag.FragCodeBlitSprite:

	case rawfrag.FragCodeParticleCloudDef:
		def := &ParticleCloudDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragParticleCloudDef))
		if err != nil {
			return fmt.Errorf("particleclouddef: %w", err)
		}
		wld.ParticleCloudDefs = append(wld.ParticleCloudDefs, def)
	case rawfrag.FragCodeMaterialDef:
		def := &MaterialDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragMaterialDef))
		if err != nil {
			return fmt.Errorf("materialdef: %w", err)
		}
		wld.MaterialDefs = append(wld.MaterialDefs, def)
	case rawfrag.FragCodeMaterialPalette:
		def := &MaterialPalette{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragMaterialPalette))
		if err != nil {
			return fmt.Errorf("materialpalette: %w", err)
		}
		wld.MaterialPalettes = append(wld.MaterialPalettes, def)
	case rawfrag.FragCodeDmSpriteDef2:
		def := &DMSpriteDef2{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragDmSpriteDef2))
		if err != nil {
			return fmt.Errorf("dmspritedef2: %w", err)
		}
		wld.lastReadModelTag = baseTagTrim(def.Tag)

		if strings.HasPrefix(def.Tag, "R") {
			tag := strings.TrimSuffix(def.Tag[1:], "_DMSPRITEDEF")
			_, err := strconv.Atoi(tag)
			if err == nil {
				if wld.WorldDef == nil {
					wld.WorldDef = &WorldDef{}
				}
				wld.WorldDef.Zone = 1
			}
		}

		wld.DMSpriteDef2s = append(wld.DMSpriteDef2s, def)
	case rawfrag.FragCodeTrackDef:
		def := &TrackDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragTrackDef))
		if err != nil {
			return fmt.Errorf("trackdef: %w", err)
		}
		wld.TrackDefs = append(wld.TrackDefs, def)
	case rawfrag.FragCodeTrack:
		def := &TrackInstance{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragTrack))
		if err != nil {
			return fmt.Errorf("track: %w", err)
		}
		wld.TrackInstances = append(wld.TrackInstances, def)
	case rawfrag.FragCodeDMSpriteDef:
		def := &DMSpriteDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragDMSpriteDef))
		if err != nil {
			return fmt.Errorf("dmspritedef: %w", err)
		}
		wld.DMSpriteDefs = append(wld.DMSpriteDefs, def)
	case rawfrag.FragCodeDMSprite:
	case rawfrag.FragCodeActorDef:
		def := &ActorDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragActorDef))
		if err != nil {
			return fmt.Errorf("actordef: %w", err)
		}

		wld.ActorDefs = append(wld.ActorDefs, def)
	case rawfrag.FragCodeActor:
		def := &ActorInst{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragActor))
		if err != nil {
			return fmt.Errorf("actor: %w", err)
		}

		wld.ActorInsts = append(wld.ActorInsts, def)
	case rawfrag.FragCodeHierarchicalSpriteDef:
		def := &HierarchicalSpriteDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragHierarchicalSpriteDef))
		if err != nil {
			return fmt.Errorf("hierarchicalspritedef: %w", err)
		}
		wld.HierarchicalSpriteDefs = append(wld.HierarchicalSpriteDefs, def)
	case rawfrag.FragCodeHierarchicalSprite:
		return nil
	case rawfrag.FragCodeLightDef:
		def := &LightDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragLightDef))
		if err != nil {
			return fmt.Errorf("lightdef: %w", err)
		}
		wld.LightDefs = append(wld.LightDefs, def)
	case rawfrag.FragCodeLight:
		return nil // light instances are ignored, since they're derived from other definitions
	case rawfrag.FragCodeSprite3DDef:
		def := &Sprite3DDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragSprite3DDef))
		if err != nil {
			return fmt.Errorf("sprite3ddef: %w", err)
		}
		wld.Sprite3DDefs = append(wld.Sprite3DDefs, def)
	case rawfrag.FragCodeSprite3D:
		// sprite instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeZone:
		def := &Zone{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragZone))
		if err != nil {
			return fmt.Errorf("zone: %w", err)
		}
		wld.Zones = append(wld.Zones, def)

	case rawfrag.FragCodeWorldTree:
		def := &WorldTree{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragWorldTree))
		if err != nil {
			return fmt.Errorf("worldtree: %w", err)
		}
		wld.WorldTrees = append(wld.WorldTrees, def)

	case rawfrag.FragCodeRegion:
		def := &Region{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragRegion))
		if err != nil {
			return fmt.Errorf("region: %w", err)
		}
		wld.Regions = append(wld.Regions, def)
	case rawfrag.FragCodeAmbientLight:
		def := &AmbientLight{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragAmbientLight))
		if err != nil {
			return fmt.Errorf("ambientlight: %w", err)
		}
		wld.AmbientLights = append(wld.AmbientLights, def)
	case rawfrag.FragCodePointLight:
		def := &PointLight{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragPointLight))
		if err != nil {
			return fmt.Errorf("pointlight: %w", err)
		}
		wld.PointLights = append(wld.PointLights, def)
	case rawfrag.FragCodePolyhedronDef:
		def := &PolyhedronDefinition{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragPolyhedronDef))
		if err != nil {
			return fmt.Errorf("polyhedrondefinition: %w", err)
		}
		wld.PolyhedronDefs = append(wld.PolyhedronDefs, def)
	case rawfrag.FragCodePolyhedron:
		// polyhedron instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeSphere:
		// sphere instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeDmRGBTrackDef:
		def := &RGBTrackDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragDmRGBTrackDef))
		if err != nil {
			return fmt.Errorf("dmrgbtrackdef: %w", err)
		}
		wld.RGBTrackDefs = append(wld.RGBTrackDefs, def)
	case rawfrag.FragCodeDmRGBTrack:
	case rawfrag.FragCodeSprite2DDef:
		def := &Sprite2DDef{}
		err := def.FromRaw(wld, rawWld, fragment.(*rawfrag.WldFragSprite2DDef))
		if err != nil {
			return fmt.Errorf("sprite2ddef: %w", err)
		}
		wld.Sprite2DDefs = append(wld.Sprite2DDefs, def)
	case rawfrag.FragCodeSprite2D:
	default:
		return fmt.Errorf("unhandled fragment type %d (%s)", fragment.FragCode(), raw.FragName(fragment.FragCode()))
	}

	return nil
}

func (wld *Wld) WriteRaw(w io.Writer) error {
	var err error
	dst := &raw.Wld{
		IsNewWorld: false,
	}
	if wld.WorldDef != nil && wld.WorldDef.NewWorld == 1 {
		dst.IsNewWorld = true
	}
	if dst.Fragments == nil {
		dst.Fragments = []model.FragmentReadWriter{}
	}
	raw.NameClear()

	if wld.GlobalAmbientLightDef != nil {
		_, err = wld.GlobalAmbientLightDef.ToRaw(wld, dst)
		if err != nil {
			return fmt.Errorf("global ambient light: %w", err)
		}
	}

	if wld.WorldDef.Zone != 1 {
		baseTags := []string{}
		for _, dmSprite := range wld.DMSpriteDef2s {
			if dmSprite.Tag == "" {
				return fmt.Errorf("dmspritedef tag is empty")
			}
			isUnique := true
			for _, baseTag := range baseTags {
				if baseTag == baseTagTrim(dmSprite.Tag) {
					isUnique = false
					break
				}
			}
			if isUnique {
				baseTags = append(baseTags, baseTagTrim(dmSprite.Tag))
			}
		}

		//sort.Strings(baseTags)

		for _, baseTag := range baseTags {

			for _, dmSprite := range wld.DMSpriteDef2s {
				dmBaseTag := baseTagTrim(dmSprite.Tag)
				if baseTag != dmBaseTag {
					continue
				}
				_, err = dmSprite.ToRaw(wld, dst)
				if err != nil {
					return fmt.Errorf("dmspritedef2 %s: %w", dmSprite.Tag, err)
				}
			}

			for _, dmSprite := range wld.DMSpriteDefs {
				dmBaseTag := baseTagTrim(dmSprite.Tag)
				if baseTag != dmBaseTag {
					continue
				}
				_, err = dmSprite.ToRaw(wld, dst)
				if err != nil {
					return fmt.Errorf("dmspritedef %s: %w", dmSprite.Tag, err)
				}
			}

			for _, hiSprite := range wld.HierarchicalSpriteDefs {
				hiBaseTag := baseTagTrim(hiSprite.Tag)
				if baseTag != hiBaseTag {
					continue
				}
				_, err = hiSprite.ToRaw(wld, dst)
				if err != nil {
					return fmt.Errorf("hierarchicalsprite %s: %w", hiSprite.Tag, err)
				}
			}

			for _, matDef := range wld.MaterialDefs {
				shortTag := baseTag
				if len(baseTag) > 3 {
					shortTag = baseTag[:3]
				}
				if !strings.HasPrefix(matDef.Tag, shortTag) {
					continue
				}

				_, err = matDef.ToRaw(wld, dst)
				if err != nil {
					return fmt.Errorf("materialdef %s: %w", matDef.Tag, err)
				}
			}

			for _, track := range wld.TrackInstances {
				if !track.Sleep.Valid {
					continue
				}

				if track.SpriteTag != baseTag {
					continue
				}

				_, err = track.ToRaw(wld, dst)
				if err != nil {
					return fmt.Errorf("track %s: %w", track.Tag, err)
				}
			}

			for _, actorDef := range wld.ActorDefs {
				if baseTag != baseTagTrim(actorDef.Tag) {
					continue
				}
				_, err = actorDef.ToRaw(wld, dst)
				if err != nil {
					return fmt.Errorf("actordef %s: %w", actorDef.Tag, err)
				}
			}
		}
	} else {
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

	dst.Fragments = append([]model.FragmentReadWriter{&rawfrag.WldFragDefault{}}, dst.Fragments...)
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
	"P01": {}, "P02": {}, "P03": {}, "P04": {}, "P05": {}, "P06": {}, "P07": {}, "P08": {}, "P09": {},
	"O02": {}, "O03": {},
	"T01": {}, "T02": {}, "T03": {}, "T04": {}, "T05": {}, "T06": {}, "T07": {}, "T08": {}, "T09": {},
}

func isAnimationPrefix(name string) bool {
	if len(name) < 3 {
		return false
	}

	name = strings.TrimSuffix(name, "_ani")

	firstUnderscore := strings.Index(name, "_")
	if firstUnderscore < 0 {
		return false
	}

	name = name[:firstUnderscore]
	if len(name) < 3 {
		return false
	}

	index := name[len(name)-2:]
	prefix := name[len(name)-3 : len(name)-2]

	isPrefix := false
	letterPrefixes := []string{"C", "D", "L", "O", "S", "P", "T"}
	for _, p := range letterPrefixes {
		if prefix == p {
			isPrefix = true
			break
		}
	}

	_, err := strconv.Atoi(index)
	if err != nil {
		return false
	}

	return isPrefix
}

func baseTagTrim(tag string) string {
	tag = strings.ReplaceAll(tag, " ", "")
	if len(tag) < 2 {
		return tag
	}

	index := strings.Index(tag, "_")
	if index > 0 {
		tag = tag[:index]
	}
	// find suffix first number
	for i := 0; i < len(tag); i++ {
		if tag[i] >= '0' && tag[i] <= '9' {
			tag = tag[:i]
			break
		}
	}

	if len(tag) > 4 && strings.HasSuffix(tag, "HE") {
		tag = tag[:len(tag)-2]
	}

	if tag == "PREPE" {
		tag = "PRE"
	}

	if strings.HasSuffix(tag, "EYE") && len(tag) >= 6 {
		tag = tag[:len(tag)-3]
	}

	return tag
}
