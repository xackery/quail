package wce

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
	"github.com/xackery/quail/tree"
)

func (wce *Wce) ReadWldRaw(src *raw.Wld) error {
	wce.reset()
	wce.maxMaterialHeads = make(map[string]int)
	wce.maxMaterialTextures = make(map[string]int)
	wce.WorldDef = &WorldDef{folders: []string{"world"}}
	if src.IsNewWorld {
		wce.WorldDef.NewWorld = 1
	}
	if src.IsZone {
		wce.WorldDef.Zone = 1
	}

	// get a list of root nodes
	roots, nodes, err := tree.BuildFragReferenceTree(wce.isChr, src)
	if err != nil {
		return fmt.Errorf("build frag reference tree: %w", err)
	}

	// make a map of folders each frag contains
	foldersByFrag := make(map[int][]string)

	// Sort root nodes by FragID before processing
	sortedRoots := make([]*tree.Node, 0, len(roots))
	for _, root := range roots {
		sortedRoots = append(sortedRoots, root)
	}

	// Sort by FragID
	sort.Slice(sortedRoots, func(i, j int) bool {
		return sortedRoots[i].FragID < sortedRoots[j].FragID
	})

	// Process the sorted roots
	for _, root := range sortedRoots {
		//fmt.Printf("Processing Root FragID: %d, Tag: %s\n", root.FragID, root.Tag)
		setRootFolder(foldersByFrag, "", root, wce.isChr, nodes, wce)
	}

	for i := 1; i < len(src.Fragments); i++ {
		fragment := src.Fragments[i]
		folders, ok := foldersByFrag[i]
		if !ok {
			return fmt.Errorf("fragment %d (%s): folders not found", i, raw.FragName(fragment.FragCode()))
		}

		err := readRawFrag(wce, src, fragment, folders)
		if err != nil {
			return fmt.Errorf("fragment %d (%s): %w", i, raw.FragName(fragment.FragCode()), err)
		}
	}

	return nil
}

func readRawFrag(e *Wce, rawWld *raw.Wld, fragment model.FragmentReadWriter, folders []string) error {
	if len(folders) == 0 {
		folders = []string{"world"}
	}
	switch fragment.FragCode() {
	case rawfrag.FragCodeGlobalAmbientLightDef:

		def := &GlobalAmbientLightDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragGlobalAmbientLightDef))
		if err != nil {
			return fmt.Errorf("globalambientlightdef: %w", err)
		}
		e.GlobalAmbientLightDef = def
	case rawfrag.FragCodeBMInfo:
		return nil
	case rawfrag.FragCodeSimpleSpriteDef:
		def := &SimpleSpriteDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragSimpleSpriteDef))
		if err != nil {
			return fmt.Errorf("simplespritedef: %w", err)
		}

		e.SimpleSpriteDefs = append(e.SimpleSpriteDefs, def)
	case rawfrag.FragCodeSimpleSprite:
		//return fmt.Errorf("simplesprite fragment found, but not expected")
	case rawfrag.FragCodeBlitSpriteDef:
		def := &BlitSpriteDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragBlitSpriteDef))
		if err != nil {
			return fmt.Errorf("blitspritedef: %w", err)
		}
		e.BlitSpriteDefs = append(e.BlitSpriteDefs, def)
	case rawfrag.FragCodeBlitSprite:
	case rawfrag.FragCodeParticleCloudDef:
		def := &ParticleCloudDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragParticleCloudDef))
		if err != nil {
			return fmt.Errorf("particleclouddef: %w", err)
		}
		e.ParticleCloudDefs = append(e.ParticleCloudDefs, def)
	case rawfrag.FragCodeMaterialDef:
		def := &MaterialDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragMaterialDef))
		if err != nil {
			return fmt.Errorf("materialdef: %w", err)
		}
		e.MaterialDefs = append(e.MaterialDefs, def)
	case rawfrag.FragCodeMaterialPalette:
		def := &MaterialPalette{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragMaterialPalette))
		if err != nil {
			return fmt.Errorf("materialpalette: %w", err)
		}
		e.MaterialPalettes = append(e.MaterialPalettes, def)
		// e.isVariationMaterial = true
	case rawfrag.FragCodeDmSpriteDef2:
		def := &DMSpriteDef2{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragDmSpriteDef2))
		if err != nil {
			return fmt.Errorf("dmspritedef2: %w", err)
		}

		if strings.HasPrefix(def.Tag, "R") {
			tag := strings.TrimSuffix(def.Tag[1:], "_DMSPRITEDEF")
			_, err := strconv.Atoi(tag)
			if err == nil {
				if e.WorldDef == nil {
					e.WorldDef = &WorldDef{folders: []string{"world"}}
				}
				e.WorldDef.Zone = 1
			}
		}

		e.DMSpriteDef2s = append(e.DMSpriteDef2s, def)
	case rawfrag.FragCodeTrack:
		def := &TrackInstance{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragTrack))
		if err != nil {
			return fmt.Errorf("track: %w", err)
		}
		e.TrackInstances = append(e.TrackInstances, def)
	case rawfrag.FragCodeTrackDef:
		def := &TrackDef{folders: folders}
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
		def := &DMTrackDef2{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragDmTrackDef2))
		if err != nil {
			return fmt.Errorf("dmtrackdef2: %w", err)
		}
		e.DMTrackDef2s = append(e.DMTrackDef2s, def)

	case rawfrag.FragCodeDMSpriteDef:
		def := &DMSpriteDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragDMSpriteDef))
		if err != nil {
			return fmt.Errorf("dmspritedef: %w", err)
		}
		e.DMSpriteDefs = append(e.DMSpriteDefs, def)
	case rawfrag.FragCodeDMSprite:
	case rawfrag.FragCodeActorDef:
		def := &ActorDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragActorDef))
		if err != nil {
			return fmt.Errorf("actordef: %w", err)
		}

		e.ActorDefs = append(e.ActorDefs, def)
		// e.isVariationMaterial = false
	case rawfrag.FragCodeActor:
		def := &ActorInst{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragActor))
		if err != nil {
			return fmt.Errorf("actor: %w", err)
		}

		e.ActorInsts = append(e.ActorInsts, def)
	case rawfrag.FragCodeHierarchicalSpriteDef:
		def := &HierarchicalSpriteDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragHierarchicalSpriteDef))
		if err != nil {
			return fmt.Errorf("hierarchicalspritedef: %w", err)
		}
		e.HierarchicalSpriteDefs = append(e.HierarchicalSpriteDefs, def)
	case rawfrag.FragCodeHierarchicalSprite:
		return nil
	case rawfrag.FragCodeLightDef:
		def := &LightDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragLightDef))
		if err != nil {
			return fmt.Errorf("lightdef: %w", err)
		}
		e.LightDefs = append(e.LightDefs, def)
	case rawfrag.FragCodeLight:
		return nil // light instances are ignored, since they're derived from other definitions
	case rawfrag.FragCodeSprite3DDef:
		def := &Sprite3DDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragSprite3DDef))
		if err != nil {
			return fmt.Errorf("sprite3ddef: %w", err)
		}
		e.Sprite3DDefs = append(e.Sprite3DDefs, def)
	case rawfrag.FragCodeSprite3D:
		// sprite instances are ignored, since they're derived from other definitions
		return nil
	case rawfrag.FragCodeZone:
		def := &Zone{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragZone))
		if err != nil {
			return fmt.Errorf("zone: %w", err)
		}
		e.Zones = append(e.Zones, def)

	case rawfrag.FragCodeWorldTree:
		def := &WorldTree{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragWorldTree))
		if err != nil {
			return fmt.Errorf("worldtree: %w", err)
		}
		e.WorldTrees = append(e.WorldTrees, def)

	case rawfrag.FragCodeRegion:
		def := &Region{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragRegion))
		if err != nil {
			return fmt.Errorf("region: %w", err)
		}
		e.Regions = append(e.Regions, def)
	case rawfrag.FragCodeAmbientLight:
		def := &AmbientLight{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragAmbientLight))
		if err != nil {
			return fmt.Errorf("ambientlight: %w", err)
		}
		e.AmbientLights = append(e.AmbientLights, def)
	case rawfrag.FragCodePointLight:
		def := &PointLight{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragPointLight))
		if err != nil {
			return fmt.Errorf("pointlight: %w", err)
		}
		e.PointLights = append(e.PointLights, def)
	case rawfrag.FragCodePolyhedronDef:
		def := &PolyhedronDefinition{folders: folders}
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
		def := &RGBTrackDef{folders: folders}
		err := def.FromRaw(e, rawWld, fragment.(*rawfrag.WldFragDmRGBTrackDef))
		if err != nil {
			return fmt.Errorf("dmrgbtrackdef: %w", err)
		}
		e.RGBTrackDefs = append(e.RGBTrackDefs, def)
	case rawfrag.FragCodeDmRGBTrack:
	case rawfrag.FragCodeSprite2DDef:
		def := &Sprite2DDef{folders: folders}
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
	// regexAniNormal    = regexp.MustCompile(`^([A-Z])([0-9]{2})([A-Z]{3}).*`)
	// regexAniAlt       = regexp.MustCompile(`^([A-Z])([0-9]{2})([A-Z])([A-Z]{3}).*`)
	// regexAniAltSuffix = regexp.MustCompile(`^([A-Z])([0-9]{2}).*_([A-Z]{3})$`)
	// regexTrackNormal  = regexp.MustCompile(`^([A-Z]{3}).*`)
	regexAniPrefix = regexp.MustCompile(`^[CDLOPST](0[1-9]|[1-9][0-9])`)
)

// returns model name (ELF, etc), sequence tag (C, P, etc), subsequence, sequence number
// if sequence number is -1, it's a bone
// func (wce *Wce) trackTagAndSequence(tag string) (string, string, string, int) {
// 	tag = strings.TrimSuffix(tag, "_TRACK")
// 	tag = strings.TrimSuffix(tag, "_TRACKDEF")
// 	m := regexTrackNormal.FindStringSubmatch(tag)
// 	if len(m) > 1 {
// 		isFound := false
// 		for _, modelTag := range wce.modelTags {
// 			if modelTag != m[1] {
// 				continue
// 			}
// 			isFound = true
// 			break
// 		}
// 		if isFound {
// 			return m[1], "", "", -1
// 		}
// 	}
// 	m = regexAniNormal.FindStringSubmatch(tag)
// 	if len(m) > 3 {
// 		isFound := false
// 		for _, modelTag := range wce.modelTags {
// 			if modelTag != m[3] {
// 				continue
// 			}
// 			isFound = true
// 			break
// 		}
// 		if isFound {
// 			seq, err := strconv.Atoi(m[2])
// 			if err == nil {
// 				return m[3], m[1], "", seq
// 			}
// 		}
// 	}

// 	m = regexAniAlt.FindStringSubmatch(tag)
// 	if len(m) > 4 {
// 		isFound := false
// 		for _, modelTag := range wce.modelTags {
// 			if modelTag != m[4] {
// 				continue
// 			}
// 			isFound = true
// 			break
// 		}
// 		if isFound {
// 			seq, err := strconv.Atoi(m[2])
// 			if err == nil {
// 				return m[4], m[1], m[3], seq
// 			}
// 		}
// 	}
// 	m = regexAniAltSuffix.FindStringSubmatch(tag)
// 	if len(m) > 1 {
// 		isFound := false
// 		for _, modelTag := range wce.modelTags {
// 			if modelTag != m[3] {
// 				continue
// 			}
// 			isFound = true
// 			break
// 		}
// 		if isFound {
// 			seq, err := strconv.Atoi(m[2])
// 			if err == nil {
// 				return m[3], m[1], "", -seq
// 			}
// 		}
// 	}

// 	return "", "", "", -1
// }

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

func setRootFolder(foldersByFrag map[int][]string, folder string, node *tree.Node, isChr bool, nodes map[int32]*tree.Node, wce *Wce) {

	// If no folder is assigned, handle specific cases based on FragType
	if len(foldersByFrag[int(node.FragID)]) == 0 {
		switch node.FragType {
		case "DmSpriteDef2":
			prefix, err := helper.DmSpriteDefTagParse(isChr, node.Tag)
			if err == nil && prefix != "" {
				folder = prefix
			}
		case "BlitSpriteDef":
			if strings.HasPrefix(node.Tag, "I_") {
				// Remove "I_" and take the part before the next "_"
				strippedTag := strings.TrimPrefix(node.Tag, "I_")
				if strings.Contains(strippedTag, "_") {
					folder = strings.Split(strippedTag, "_")[0]
				} else {
					folder = strippedTag
				}
			}
		case "MaterialDef":
			prefix, err := helper.MaterialTagParse(isChr, node.Tag)
			if err == nil && prefix != "" {
				if prefix == "CLK04" {
					for _, potentialNode := range nodes {
						if potentialNode.FragType == "MaterialPalette" {
							// Check the child nodes of the MaterialPalette node
							for _, childNode := range potentialNode.Children {
								if strings.HasPrefix(childNode.Tag, "CLK04") {
									// Add the first 3 characters of each MaterialPalette node's tags to foldersByFrag
									folderToAdd := potentialNode.Tag[:3]
									foldersByFrag[int(node.FragID)] = appendUnique(foldersByFrag[int(node.FragID)], folderToAdd)
								}
							}
						}
					}
				} else {
					// Use the returned prefix directly as the folder
					folder = prefix
				}
			}
		case "Track":
			if wce.isTrackAni(node.Tag) {
				// Call TrackAnimationParse only if isTrackAni returns true
				_, prefix := helper.TrackAnimationParse(isChr, node.Tag)
				if prefix != "" {
					folder = prefix
				}
			} else {
				// If isTrackAni returns false, set folder to the first 3 characters of node.Tag
				if len(node.Tag) >= 3 {
					folder = node.Tag[:3]
				} else {
					folder = node.Tag // Fallback to the full tag if it's shorter than 3 characters
				}
			}
		default:
			folder = node.Tag
			if strings.Contains(folder, "_") {
				folder = strings.Split(folder, "_")[0]
			}
		}
	}

	foldersByFrag[int(node.FragID)] = append(foldersByFrag[int(node.FragID)], folder)

	// Pass the folder down to child nodes
	addChildrenFolder(foldersByFrag, folder, node)
}

func addChildrenFolder(foldersByFrag map[int][]string, folder string, node *tree.Node) {
	// Propagate the folder to the children
	for _, child := range node.Children {
		// Add the folder to the child node's folder list if it's not already there
		isUnique := true
		for _, existingFolder := range foldersByFrag[int(child.FragID)] {
			if existingFolder == folder {
				isUnique = false
				break
			}
		}
		if isUnique {
			foldersByFrag[int(child.FragID)] = append(foldersByFrag[int(child.FragID)], folder)
		}

		// Recursively process the child nodes
		addChildrenFolder(foldersByFrag, folder, child)
	}
}

func appendUnique(slice []string, value string) []string {
	for _, v := range slice {
		if v == value {
			return slice
		}
	}
	return append(slice, value)
}
