package wce

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/xackery/quail/helper"
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

func readRawFrag(e *Wce, rawWld *raw.Wld, fragment helper.FragmentReadWriter, folders []string) error {
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

	err = wce.convertEQGToWld()
	if err != nil {
		return fmt.Errorf("convert eqg to wld: %w", err)
	}

	dst := &raw.Wld{
		IsNewWorld: false,
	}
	if wce.WorldDef != nil && wce.WorldDef.NewWorld == 1 {
		dst.IsNewWorld = true
	}
	if dst.Fragments == nil {
		dst.Fragments = []helper.FragmentReadWriter{}
	}
	dst.NameClear()

	if wce.GlobalAmbientLightDef != nil {
		_, err = wce.GlobalAmbientLightDef.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("global ambient light: %w", err)
		}
	}

	// Write spell blit particles? (SPB)
	for _, blitSprite := range wce.BlitSpriteDefs {
		if !strings.HasSuffix(blitSprite.Tag, "_SPB") {
			continue
		}
		_, err = blitSprite.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("blitsprite %s: %w", blitSprite.Tag, err)
		}
	}

	// Write other particle cloud blits
	for _, blitSprite := range wce.BlitSpriteDefs {
		if !strings.HasSuffix(blitSprite.Tag, "_SPRITE") {
			continue
		}
		_, err = blitSprite.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("blitsprite %s: %w", blitSprite.Tag, err)
		}
	}

	// Write spell effect actordefs
	for _, actorDef := range wce.ActorDefs {
		if actorDef.fragID > 0 {
			continue
		}

		// Check if any LevelOfDetail's SpriteTag ends with "_SPRITE"
		hasSpriteTag := false
		for _, action := range actorDef.Actions {
			for _, lod := range action.LevelOfDetails {
				if strings.HasSuffix(lod.SpriteTag, "_SPRITE") {
					hasSpriteTag = true
					break
				}
			}
			if hasSpriteTag {
				break
			}
		}

		if !hasSpriteTag {
			continue
		}

		_, err := actorDef.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("actordef %s: %w", actorDef.Tag, err)
		}
	}

	// Write particle clouds
	for _, cloudDef := range wce.ParticleCloudDefs {
		_, err = cloudDef.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("cloud %s: %w", cloudDef.Tag, err)
		}
	}

	// Write other blits (for 2D Sprites and stuff)
	for _, blitSprite := range wce.BlitSpriteDefs {
		if strings.HasSuffix(blitSprite.Tag, "_SPB") || strings.HasSuffix(blitSprite.Tag, "_SPRITE") {
			// Skip if the Tag ends with _SPB or _SPRITE
			continue
		}
		// Process the blitSprite if it doesn't end with _SPB or _SPRITE
		_, err := blitSprite.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("blitsprite %s: %w", blitSprite.Tag, err)
		}
	}

	// Write out CHR_EYE materials
	for _, matDef := range wce.MaterialDefs {
		if !strings.HasPrefix(matDef.Tag, "CHR_EYE") {
			continue
		}
		_, err = matDef.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("materialdef %s: %w", matDef.Tag, err)
		}
	}

	// Write out "variation" materials
	for _, matDef := range wce.MaterialDefs {
		if matDef.Variation == 0 {
			continue
		}
		_, err = matDef.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("materialdef %s: %w", matDef.Tag, err)
		}
	}

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

	// Write non-spell effect actordefs
	for _, actorDef := range wce.ActorDefs {
		if actorDef.fragID > 0 {
			continue
		}

		hasSpriteTag := false
		for _, action := range actorDef.Actions {
			for _, lod := range action.LevelOfDetails {
				if strings.HasSuffix(lod.SpriteTag, "_SPRITE") {
					hasSpriteTag = true
					break
				}
			}
			if hasSpriteTag {
				break
			}
		}

		if hasSpriteTag {
			continue
		}

		// fmt.Printf("Processing ActorDef: %s\n", actorDef.Tag)

		_, err := actorDef.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("actordef %s: %w", actorDef.Tag, err)
		}
	}

	for _, zone := range wce.Zones {
		_, err = zone.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("zone %s: %w", zone.Tag, err)
		}
	}

	dst.Fragments = append([]helper.FragmentReadWriter{&rawfrag.WldFragDefault{}}, dst.Fragments...)
	return dst.Write(w)
}

var (
	regexAniPrefix = regexp.MustCompile(`^[CDLOPST](0[1-9]|[1-9][0-9])`)
)

func (wce *Wce) isTrackAni(tag string) bool {
	// If isObj is true, it's not an animation track
	if wce.isObj {
		return false
	}

	// Check if the tag starts with the specified regex pattern
	if regexAniPrefix.MatchString(tag) {
		return true
	}

	return false
}

func setRootFolder(foldersByFrag map[int][]string, folder string, node *tree.Node, isChr bool, nodes map[int32]*tree.Node, wce *Wce) {

	// If no folder is assigned, handle specific cases based on FragType
	if len(foldersByFrag[int(node.FragID)]) == 0 {
		switch node.FragType {
		case "ActorDef":
			folder = node.Tag
			if strings.Contains(folder, "_ACTORDEF") {
				folder = strings.Split(folder, "_ACTORDEF")[0]
			}
		case "AmbientLight":
			folder = "ZONE"
		case "DmSpriteDef2":
			prefix, err := helper.DmSpriteDefTagParse(isChr, node.Tag)
			if err == nil && prefix != "" {
				folder = prefix
			}
		case "BlitSpriteDef":
			if strings.HasPrefix(node.Tag, "I_") {
				// If the tag starts with "I_", split the tag after "I_" by the next "_"
				strippedTag := strings.TrimPrefix(node.Tag, "I_")
				if strings.Contains(strippedTag, "_") {
					folder = "I_" + strings.Split(strippedTag, "_")[0]
				} else {
					folder = node.Tag
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
									addChildrenFolder(foldersByFrag, folderToAdd, node)
								}
							}
						}
					}
				} else {
					// Use the returned prefix directly as the folder
					folder = prefix
				}
			}
		case "Region":
			folder = "REGION"
		case "Track":
			if wce.isTrackAni(node.Tag) {
				_, prefix := helper.TrackAnimationParse(isChr, node.Tag)
				if prefix != "" {
					folder = prefix
				}
			} else {
				if strings.HasPrefix(node.Tag, "IT") {
					if strings.Contains(node.Tag, "_") {
						folder = strings.SplitN(node.Tag, "_", 2)[0]
					} else {
						folder = node.Tag // Use the whole tag if there's no "_"
					}
				} else {
					if len(node.Tag) >= 3 {
						folder = node.Tag[:3]
					} else {
						folder = node.Tag // Use the full tag if it's shorter than 3 characters
					}
				}
			}
		case "Zone":
			folder = "ZONE"
		default:
			folder = node.Tag
			if strings.Contains(folder, "_") {
				folder = strings.Split(folder, "_")[0]
			}
		}
	}

	foldersByFrag[int(node.FragID)] = appendUnique(foldersByFrag[int(node.FragID)], folder)

	// Pass the folder down to child nodes
	addChildrenFolder(foldersByFrag, folder, node)
}

func addChildrenFolder(foldersByFrag map[int][]string, folder string, node *tree.Node) {
	// Propagate the folder to the children
	for _, child := range node.Children {

		foldersByFrag[int(child.FragID)] = appendUnique(foldersByFrag[int(child.FragID)], folder)

		// Recursively process the child nodes
		addChildrenFolder(foldersByFrag, folder, child)
	}
}

func appendUnique(slice []string, value string) []string {
	// Skip adding empty strings
	if strings.TrimSpace(value) == "" {
		return slice
	}
	for _, v := range slice {
		if v == value {
			return slice
		}
	}
	return append(slice, value)
}

func (wce *Wce) convertEQGToWld() error {

	//for _, mds := range wce.MdsDefs {
	//}

	//for _, mod := range wce.ModDefs {
	//}

	//for _, ter := range wce.TerDefs {
	//}

	//for _, ani := range wce.AniDefs {
	//}

	//for _, lay := range wce.LayDefs {
	//}

	//for _, pts := range wce.PtsDefs {

	//}

	//for _, prt := range wce.PrtDefs {
	//}
	return nil
}
