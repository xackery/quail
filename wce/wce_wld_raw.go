package wce

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
	"github.com/xackery/quail/tree"
)

func (wce *Wce) ReadWldRaw(src *raw.Wld) error {
	wce.reset()
	wce.maxMaterialHeads = make(map[string]int)
	wce.maxMaterialTextures = make(map[string]int)
	wce.WorldDef = &WorldDef{}
	if src.IsNewWorld {
		wce.WorldDef.NewWorld = 1
	}
	if src.IsZone {
		wce.WorldDef.Zone = 1
	}

	roots, err := tree.BuildFragReferenceTree(src)
	if err != nil {
		return fmt.Errorf("build frag reference tree: %w", err)
	}

	folders := make(map[int]string)

	// Traverse and print the trees
	fmt.Println("Debug tree:")
	for _, root := range roots {
		fmt.Printf("Root ")
		tree.PrintNode(root, 0)
		setChildrenFolder(folders, root)
	}

	for i := 1; i < len(src.Fragments); i++ {
		fragment := src.Fragments[i]
		folder, ok := folders[i]
		if !ok {
			return fmt.Errorf("fragment %d (%s): folder not found", i, raw.FragName(fragment.FragCode()))
		}

		err := readRawFrag(wce, src, fragment, folder)
		if err != nil {
			return fmt.Errorf("fragment %d (%s): %w", i, raw.FragName(fragment.FragCode()), err)
		}
	}

	return nil
}

func readRawFrag(e *Wce, rawWld *raw.Wld, fragment model.FragmentReadWriter, folder string) error {
	switch fragment.FragCode() {
	case rawfrag.FragCodeGlobalAmbientLightDef:

		def := &GlobalAmbientLightDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragGlobalAmbientLightDef))
		if err != nil {
			return fmt.Errorf("globalambientlightdef: %w", err)
		}
		e.GlobalAmbientLightDef = def
	case rawfrag.FragCodeBMInfo:
		return nil
	case rawfrag.FragCodeSimpleSpriteDef:
		def := &SimpleSpriteDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragSimpleSpriteDef))
		if err != nil {
			return fmt.Errorf("simplespritedef: %w", err)
		}

		e.SimpleSpriteDefs = append(e.SimpleSpriteDefs, def)
	case rawfrag.FragCodeSimpleSprite:
		//return fmt.Errorf("simplesprite fragment found, but not expected")
	case rawfrag.FragCodeBlitSpriteDef:
		def := &BlitSpriteDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragBlitSpriteDef))
		if err != nil {
			return fmt.Errorf("blitspritedef: %w", err)
		}
		e.BlitSpriteDefs = append(e.BlitSpriteDefs, def)
	case rawfrag.FragCodeBlitSprite:
	case rawfrag.FragCodeParticleCloudDef:
		def := &ParticleCloudDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragParticleCloudDef))
		if err != nil {
			return fmt.Errorf("particleclouddef: %w", err)
		}
		e.ParticleCloudDefs = append(e.ParticleCloudDefs, def)
	case rawfrag.FragCodeMaterialDef:
		def := &MaterialDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragMaterialDef))
		if err != nil {
			return fmt.Errorf("materialdef: %w", err)
		}
		e.MaterialDefs = append(e.MaterialDefs, def)
	case rawfrag.FragCodeMaterialPalette:
		def := &MaterialPalette{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragMaterialPalette))
		if err != nil {
			return fmt.Errorf("materialpalette: %w", err)
		}
		e.MaterialPalettes = append(e.MaterialPalettes, def)
		e.isVariationMaterial = true
	case rawfrag.FragCodeDmSpriteDef2:
		def := &DMSpriteDef2{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragDmSpriteDef2))
		if err != nil {
			return fmt.Errorf("dmspritedef2: %w", err)
		}

		if strings.HasPrefix(def.Tag, "R") {
			tag := strings.TrimSuffix(def.Tag[1:], "_DMSPRITEDEF")
			_, err := strconv.Atoi(tag)
			if err == nil {
				if e.WorldDef == nil {
					e.WorldDef = &WorldDef{}
				}
				e.WorldDef.Zone = 1
			}
		}

		e.DMSpriteDef2s = append(e.DMSpriteDef2s, def)
	case rawfrag.FragCodeTrack:
		def := &TrackInstance{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragTrack))
		if err != nil {
			return fmt.Errorf("track: %w", err)
		}
		e.TrackInstances = append(e.TrackInstances, def)
	case rawfrag.FragCodeTrackDef:
		def := &TrackDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragTrackDef))
		if err != nil {
			return fmt.Errorf("trackdef: %w", err)
		}
		e.TrackDefs = append(e.TrackDefs, def)

	case rawfrag.FragCodeDMTrack:
		frag := fragment.(*rawfrag.WldFragDMTrack)
		if frag.Flags != 0 {
			return fmt.Errorf("dmtrack: unexpected flags %d, report this to xack", frag.Flags)
		}
	case rawfrag.FragCodeDmTrackDef2:
		def := &DMTrackDef2{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragDmTrackDef2))
		if err != nil {
			return fmt.Errorf("dmtrackdef2: %w", err)
		}
		e.DMTrackDef2s = append(e.DMTrackDef2s, def)

	case rawfrag.FragCodeDMSpriteDef:
		def := &DMSpriteDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragDMSpriteDef))
		if err != nil {
			return fmt.Errorf("dmspritedef: %w", err)
		}
		e.DMSpriteDefs = append(e.DMSpriteDefs, def)
	case rawfrag.FragCodeDMSprite:
	case rawfrag.FragCodeActorDef:
		def := &ActorDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragActorDef))
		if err != nil {
			return fmt.Errorf("actordef: %w", err)
		}

		e.ActorDefs = append(e.ActorDefs, def)
		e.isVariationMaterial = false
	case rawfrag.FragCodeActor:
		def := &ActorInst{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragActor))
		if err != nil {
			return fmt.Errorf("actor: %w", err)
		}

		e.ActorInsts = append(e.ActorInsts, def)
	case rawfrag.FragCodeHierarchicalSpriteDef:
		def := &HierarchicalSpriteDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragHierarchicalSpriteDef))
		if err != nil {
			return fmt.Errorf("hierarchicalspritedef: %w", err)
		}
		e.HierarchicalSpriteDefs = append(e.HierarchicalSpriteDefs, def)
	case rawfrag.FragCodeHierarchicalSprite:
		return nil
	case rawfrag.FragCodeLightDef:
		def := &LightDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragLightDef))
		if err != nil {
			return fmt.Errorf("lightdef: %w", err)
		}
		e.LightDefs = append(e.LightDefs, def)
	case rawfrag.FragCodeLight:
		return nil // light instances are ignored, since they're derived from other definitions
	case rawfrag.FragCodeSprite3DDef:
		def := &Sprite3DDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragSprite3DDef))
		if err != nil {
			return fmt.Errorf("sprite3ddef: %w", err)
		}
		e.Sprite3DDefs = append(e.Sprite3DDefs, def)
	case rawfrag.FragCodeSprite3D:
		// sprite instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeZone:
		def := &Zone{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragZone))
		if err != nil {
			return fmt.Errorf("zone: %w", err)
		}
		e.Zones = append(e.Zones, def)

	case rawfrag.FragCodeWorldTree:
		def := &WorldTree{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragWorldTree))
		if err != nil {
			return fmt.Errorf("worldtree: %w", err)
		}
		e.WorldTrees = append(e.WorldTrees, def)

	case rawfrag.FragCodeRegion:
		def := &Region{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragRegion))
		if err != nil {
			return fmt.Errorf("region: %w", err)
		}
		e.Regions = append(e.Regions, def)
	case rawfrag.FragCodeAmbientLight:
		def := &AmbientLight{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragAmbientLight))
		if err != nil {
			return fmt.Errorf("ambientlight: %w", err)
		}
		e.AmbientLights = append(e.AmbientLights, def)
	case rawfrag.FragCodePointLight:
		def := &PointLight{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragPointLight))
		if err != nil {
			return fmt.Errorf("pointlight: %w", err)
		}
		e.PointLights = append(e.PointLights, def)
	case rawfrag.FragCodePolyhedronDef:
		def := &PolyhedronDefinition{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragPolyhedronDef))
		if err != nil {
			return fmt.Errorf("polyhedrondefinition: %w", err)
		}
		e.PolyhedronDefs = append(e.PolyhedronDefs, def)
	case rawfrag.FragCodePolyhedron:
		// polyhedron instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeSphere:
		// sphere instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeDmRGBTrackDef:
		def := &RGBTrackDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragDmRGBTrackDef))
		if err != nil {
			return fmt.Errorf("dmrgbtrackdef: %w", err)
		}
		e.RGBTrackDefs = append(e.RGBTrackDefs, def)
	case rawfrag.FragCodeDmRGBTrack:
	case rawfrag.FragCodeSprite2DDef:
		def := &Sprite2DDef{folder: folder}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragSprite2DDef))
		if err != nil {
			return fmt.Errorf("sprite2ddef: %w", err)
		}
		e.Sprite2DDefs = append(e.Sprite2DDefs, def)
	case rawfrag.FragCodeSprite2D:
	default:
		return fmt.Errorf("unhandled fragment type %d (%s)", fragment.FragCode(), raw.FragName(fragment.FragCode()))
	}

	return nil
}

func (wce *Wce) WriteWldRaw(w io.Writer) error {
	var err error
	dst := &raw.Wld{
		IsNewWorld: false,
	}
	if wce.WorldDef != nil && wce.WorldDef.NewWorld == 1 {
		dst.IsNewWorld = true
	}
	if dst.Fragments == nil {
		dst.Fragments = []model.FragmentReadWriter{}
	}
	dst.NameClear()

	if wce.GlobalAmbientLightDef != nil {
		_, err = wce.GlobalAmbientLightDef.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("global ambient light: %w", err)
		}
	}

	wce.modelTags = []string{}

	if wce.WorldDef.Zone != 1 {
		baseTags := []string{}
		for _, actorDef := range wce.ActorDefs {
			if actorDef.Tag == "" {
				return fmt.Errorf("dmspritedef tag is empty")
			}
			isUnique := true
			for _, baseTag := range baseTags {
				if baseTag == baseTagTrim(wce.isObj, actorDef.Tag) {
					isUnique = false
					break
				}
			}
			if isUnique {
				baseTags = append(baseTags, baseTagTrim(wce.isObj, actorDef.Tag))
			}
			wce.modelTags = append(wce.modelTags, baseTagTrim(wce.isObj, actorDef.Tag))
		}

		//sort.Strings(baseTags)

		clouds := []string{}
		for _, cloud := range wce.ParticleCloudDefs {
			isUnique := true
			for _, bstr := range clouds {
				if bstr == cloud.Tag {
					isUnique = false
					break
				}
			}
			if !isUnique {
				continue
			}
			clouds = append(clouds, cloud.Tag)
		}
		sort.Strings(clouds)

		for _, cloud := range clouds {
			for _, cloudDef := range wce.ParticleCloudDefs {
				if cloud != cloudDef.Tag {
					continue
				}
				_, err = cloudDef.ToRaw(wce, dst)
				if err != nil {
					return fmt.Errorf("cloud %s: %w", cloudDef.Tag, err)
				}
			}
		}

		clks := make(map[string]bool)
		for _, matDef := range wce.MaterialDefs {
			if !strings.HasPrefix(matDef.Tag, "CLK") {
				continue
			}

			_, err = strconv.Atoi(matDef.Tag[3:6])
			if err != nil {
				continue
			}
			if clks[matDef.Tag] {
				continue
			}
			clks[matDef.Tag] = true

			_, err = matDef.ToRaw(wce, dst)
			if err != nil {
				return fmt.Errorf("materialdef %s: %w", matDef.Tag, err)
			}
		}

		for _, baseTag := range baseTags {

			for _, actorDef := range wce.ActorDefs {
				if baseTag != baseTagTrim(wce.isObj, actorDef.Tag) {
					continue
				}
				_, err = actorDef.ToRaw(wce, dst)
				if err != nil {
					return fmt.Errorf("actordef %s: %w", actorDef.Tag, err)
				}
			}

			for _, hiSprite := range wce.HierarchicalSpriteDefs {
				hiBaseTag := baseTagTrim(wce.isObj, hiSprite.Tag)
				if baseTag != hiBaseTag {
					continue
				}
				_, err = hiSprite.ToRaw(wce, dst)
				if err != nil {
					return fmt.Errorf("hierarchicalsprite %s: %w", hiSprite.Tag, err)
				}
			}

			for _, dmSprite := range wce.DMSpriteDef2s {
				dmBaseTag := baseTagTrim(wce.isObj, dmSprite.Tag)
				if baseTag != dmBaseTag {
					continue
				}
				_, err = dmSprite.ToRaw(wce, dst)
				if err != nil {
					return fmt.Errorf("dmspritedef2 %s: %w", dmSprite.Tag, err)
				}
			}

			for _, dmSprite := range wce.DMSpriteDefs {
				dmBaseTag := baseTagTrim(wce.isObj, dmSprite.Tag)
				if baseTag != dmBaseTag {
					continue
				}
				_, err = dmSprite.ToRaw(wce, dst)
				if err != nil {
					return fmt.Errorf("dmspritedef %s: %w", dmSprite.Tag, err)
				}
			}

			for _, track := range wce.TrackInstances {

				if wce.isTrackAni(track.Tag) {
					continue
				}

				if track.SpriteTag != baseTag {
					continue
				}

				_, err = track.ToRaw(wce, dst)
				if err != nil {
					return fmt.Errorf("track %s: %w", track.Tag, err)
				}
			}

		}
	} else {
		for _, dmSprite := range wce.DMSpriteDef2s {
			_, err = dmSprite.ToRaw(wce, dst)
			if err != nil {
				return fmt.Errorf("dmspritedef2 %s: %w", dmSprite.Tag, err)
			}
		}
		for _, dmSprite := range wce.DMSpriteDefs {
			_, err = dmSprite.ToRaw(wce, dst)
			if err != nil {
				return fmt.Errorf("dmspritedef %s: %w", dmSprite.Tag, err)
			}
		}
		for _, hiSprite := range wce.HierarchicalSpriteDefs {
			_, err = hiSprite.ToRaw(wce, dst)
			if err != nil {
				return fmt.Errorf("hierarchicalsprite %s: %w", hiSprite.Tag, err)
			}
		}

	}

	for _, light := range wce.PointLights {
		_, err = light.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("pointlight %s: %w", light.Tag, err)
		}
	}

	for _, sprite := range wce.Sprite3DDefs {
		_, err = sprite.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
		}
	}

	for _, tree := range wce.WorldTrees {
		_, err = tree.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("worldtree: %w", err)
		}
	}

	for _, region := range wce.Regions {
		_, err = region.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("region %s: %w", region.Tag, err)
		}
	}

	for _, alight := range wce.AmbientLights {
		_, err = alight.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("ambientlight %s: %w", alight.Tag, err)
		}
	}

	for _, actor := range wce.ActorInsts {
		_, err = actor.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("actor %s: %w", actor.Tag, err)
		}
	}

	for _, track := range wce.TrackInstances {
		if track.fragID > 0 {
			continue
		}

		_, err = track.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("track %s: %w", track.Tag, err)
		}

	}

	for _, actorDef := range wce.ActorDefs {
		if actorDef.fragID > 0 {
			continue
		}

		_, err = actorDef.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("actordef %s: %w", actorDef.Tag, err)
		}
	}

	// Write out BlitSpriteDefs
	for _, blitSprite := range wce.BlitSpriteDefs {
		_, err = blitSprite.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("blitsprite %s: %w", blitSprite.Tag, err)
		}
	}

	for _, zone := range wce.Zones {
		_, err = zone.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("zone %s: %w", zone.Tag, err)
		}
	}

	dst.Fragments = append([]model.FragmentReadWriter{&rawfrag.WldFragDefault{}}, dst.Fragments...)
	return dst.Write(w)
}

var (
	regexAniNormal    = regexp.MustCompile(`^([A-Z])([0-9]{2})([A-Z]{3}).*`)
	regexAniAlt       = regexp.MustCompile(`^([A-Z])([0-9]{2})([A-Z])([A-Z]{3}).*`)
	regexAniAltSuffix = regexp.MustCompile(`^([A-Z])([0-9]{2}).*_([A-Z]{3})$`)
	regexTrackNormal  = regexp.MustCompile(`^([A-Z]{3}).*`)
	regexAniPrefix    = regexp.MustCompile(`^[CDLOPST](0[1-9]|[1-9][0-9])`)
)

// returns model name (ELF, etc), sequence tag (C, P, etc), subsequence, sequence number
// if sequence number is -1, it's a bone
func (wce *Wce) trackTagAndSequence(tag string) (string, string, string, int) {
	tag = strings.TrimSuffix(tag, "_TRACK")
	tag = strings.TrimSuffix(tag, "_TRACKDEF")
	m := regexTrackNormal.FindStringSubmatch(tag)
	if len(m) > 1 {
		isFound := false
		for _, modelTag := range wce.modelTags {
			if modelTag != m[1] {
				continue
			}
			isFound = true
			break
		}
		if isFound {
			return m[1], "", "", -1
		}
	}
	m = regexAniNormal.FindStringSubmatch(tag)
	if len(m) > 3 {
		isFound := false
		for _, modelTag := range wce.modelTags {
			if modelTag != m[3] {
				continue
			}
			isFound = true
			break
		}
		if isFound {
			seq, err := strconv.Atoi(m[2])
			if err == nil {
				return m[3], m[1], "", seq
			}
		}
	}

	m = regexAniAlt.FindStringSubmatch(tag)
	if len(m) > 4 {
		isFound := false
		for _, modelTag := range wce.modelTags {
			if modelTag != m[4] {
				continue
			}
			isFound = true
			break
		}
		if isFound {
			seq, err := strconv.Atoi(m[2])
			if err == nil {
				return m[4], m[1], m[3], seq
			}
		}
	}
	m = regexAniAltSuffix.FindStringSubmatch(tag)
	if len(m) > 1 {
		isFound := false
		for _, modelTag := range wce.modelTags {
			if modelTag != m[3] {
				continue
			}
			isFound = true
			break
		}
		if isFound {
			seq, err := strconv.Atoi(m[2])
			if err == nil {
				return m[3], m[1], "", -seq
			}
		}
	}

	return "", "", "", -1
}

func (wce *Wce) isTrackAni(tag string) bool {
	// If isObj is true, it's not a track animation
	if wce.isObj {
		return false
	}

	// Check if the tag starts with the specified regex pattern
	if regexAniPrefix.MatchString(tag) {
		return true
	}

	return false
}

func baseTagTrim(isObj bool, tag string) string {
	tag = strings.ReplaceAll(tag, " ", "")
	if len(tag) < 2 {
		return tag
	}

	index := strings.Index(tag, "_")
	if index > 0 {
		tag = tag[:index]
	}
	/*
		if !isObj && !strings.HasPrefix(tag, "IT") {
			// find suffix first number
			for i := 0; i < len(tag); i++ {
				if tag[i] >= '0' && tag[i] <= '9' {
					tag = tag[:i]
					break
				}
			}
		} */

	if len(tag) > 4 && strings.HasSuffix(tag, "HE") {
		tag = tag[:len(tag)-2]
	}

	if tag == "PREPE" {
		tag = "PRE"
	}

	if strings.HasSuffix(tag, "EYE") && len(tag) >= 6 {
		tag = tag[:len(tag)-3]
	}

	if len(tag) == 7 {
		tag = strings.TrimSuffix(tag, "MESH")
	}
	if len(tag) == 6 {
		tag = strings.TrimSuffix(tag, "BOD")
	}
	return tag
}

// Dummy strings used in tag matching
var dummyStrings = []string{
	"10404P0", "2HNSWORD", "BARDING", "BELT", "BODY", "BONE",
	"BOW", "BOX", "DUMMY", "HUMEYE", "MESH", "POINT", "POLYSURF",
	"RIDER", "SHOULDER",
}

// Root patterns for animation and model parsing
var rootPatterns = []string{
	`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A-Z]{3}_TRACK$`,
	`^([C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])){2}_[A-Z]{3}_TRACK$`,
	`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A-Z]{3}[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A-Z]{3}_TRACK$`,
	`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A-Z]{3}[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])_[A-Z]{3}_TRACK$`,
	`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A,B,G][A-Z]{3}[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A,B,G]_[A-Z]{3}_TRACK$`,
	`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A,B,G][C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])_[A-Z]{3}_TRACK$`,
}

// Item patterns for non-character cases
var itemPatterns = []string{
	`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])IT\d+_TRACK$`,
	`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])_IT\d+_TRACK$`,
	`^([C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])){2}_IT\d+_TRACK$`,
}

func (wce *Wce) trackAnimationParse(tag string) (string, string) {
	// Check if the tag starts with currentAniCode + currentAniModelCode
	combinedCode := wce.currentAniCode + wce.currentAniModelCode
	if wce.currentAniCode != "" && wce.currentAniModelCode != "" && strings.HasPrefix(tag, combinedCode) {
		return wce.currentAniCode, wce.currentAniModelCode
	}

	// Check against previousAnimations
	for previous := range wce.previousAnimations {
		if strings.HasPrefix(tag, previous) {
			parts := strings.Split(previous, ":")
			if len(parts) == 2 {
				return parts[0], parts[1]
			}
		}
	}

	// Check if the tag starts with the currentAniCode and contains a dummy string
	for _, dummy := range dummyStrings {
		if strings.HasPrefix(tag, wce.currentAniCode) && strings.Contains(tag, dummy) {
			return wce.currentAniCode, wce.currentAniModelCode
		}
	}

	// Handle special cases when isChr is true
	if wce.isChr {
		if strings.HasPrefix(tag, wce.currentAniCode) {
			if wce.currentAniModelCode == "SED" && len(tag) >= 6 && tag[3:6] == "FDD" {
				return wce.currentAniCode, wce.currentAniModelCode
			}
			if wce.currentAniModelCode == "FMP" && len(tag) >= len(wce.currentAniCode)+2 {
				suffixStartIndex := len(wce.currentAniCode)
				for _, suffix := range []string{"PE", "CH", "NE", "HE", "BI", "FO", "TH", "CA", "BO"} {
					if strings.HasPrefix(tag[suffixStartIndex:], suffix) {
						return wce.currentAniCode, wce.currentAniModelCode
					}
				}
			}
			if wce.currentAniModelCode == "SKE" && len(tag) >= len(wce.currentAniCode)+2 {
				suffixStartIndex := len(wce.currentAniCode)
				for _, suffix := range []string{"BI", "BO", "CA", "CH", "FA", "FI", "FO", "HA", "HE", "L_POINT", "NE", "PE", "R_POINT", "SH", "TH", "TO", "TU"} {
					if strings.HasPrefix(tag[suffixStartIndex:], suffix) {
						return wce.currentAniCode, wce.currentAniModelCode
					}
				}
			}
		}

		// Attempt to match root patterns
		for _, pattern := range rootPatterns {
			matched, _ := regexp.MatchString(pattern, tag)
			if matched {
				switch pattern {
				case rootPatterns[0]: // Pattern 1
					return handleNewAniModelCode(wce, tag[:3], tag[3:6])
				case rootPatterns[1]: // Pattern 2
					return handleNewAniModelCode(wce, tag[:3], tag[7:10])
				case rootPatterns[2], rootPatterns[3]: // Pattern 3 and 4
					return handleNewAniModelCode(wce, tag[:3], tag[3:6])
				case rootPatterns[4]: // Pattern 5
					return handleNewAniModelCode(wce, tag[:4], tag[4:7])
				case rootPatterns[5]: // Pattern 6
					return handleNewAniModelCode(wce, tag[:4], tag[8:11])
				}
			}
		}

		// Fallback for isChr
		if len(tag) >= 6 {
			newAniCode := tag[:3]
			newModelCode := tag[3:6]

			return handleNewAniModelCode(wce, newAniCode, newModelCode)
		}

		// If the tag is too short, return empty values
		return "", ""
	}

	// Special cases for isChr == false
	if strings.HasPrefix(tag, wce.currentAniCode) {
		if wce.currentAniModelCode == "IT157" && len(tag) >= 6 && tag[3:6] == "SNA" {
			return wce.currentAniCode, wce.currentAniModelCode
		}
		if wce.currentAniModelCode == "IT61" && len(tag) >= 6 && tag[3:6] == "WIP" {
			return wce.currentAniCode, wce.currentAniModelCode
		}
	}

	// Handle item patterns if isChr is false
	for _, pattern := range itemPatterns {
		matched, _ := regexp.MatchString(pattern, tag)
		if matched {
			newAniCode := tag[:3]
			modelCodeStart := strings.Index(tag, "IT") + 2
			modelCodeEnd := modelCodeStart
			for modelCodeEnd < len(tag) && tag[modelCodeEnd] >= '0' && tag[modelCodeEnd] <= '9' {
				modelCodeEnd++
			}
			return handleNewAniModelCode(wce, newAniCode, "IT"+tag[modelCodeStart:modelCodeEnd])
		}
	}

	// Default fallback for isChr == false
	if len(tag) >= 6 {
		aniCode := tag[:3]
		modelCode := "IT"
		for i := 3; i < len(tag); i++ {
			if tag[i] >= '0' && tag[i] <= '9' {
				modelCode += string(tag[i])
			} else {
				break
			}
		}
		return handleNewAniModelCode(wce, aniCode, modelCode)
	}

	return "", ""
}

// Helper function to handle new animation and model codes
func handleNewAniModelCode(wce *Wce, newAniCode, newModelCode string) (string, string) {
	wce.previousAnimations[wce.currentAniCode+wce.currentAniModelCode] = struct{}{}
	wce.currentAniCode = newAniCode
	wce.currentAniModelCode = newModelCode
	return newAniCode, newModelCode
}

func setChildrenFolder(folders map[int]string, node *tree.Node) {
	tag := node.Parent
	if tag == "" {
		tag = node.Tag
	}
	if strings.Contains(tag, "_") {
		tag = strings.Split(tag, "_")[0]
	}
	folders[int(node.FragID)] = tag
	for _, child := range node.Children {
		setChildrenFolder(folders, child)
	}
}
