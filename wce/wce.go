// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package wce

import (
	"strings"

	"github.com/xackery/quail/raw"
)

var AsciiVersion = "v0.0.1"

// Wce is a struct representing a Wce file
type Wce struct {
	isVariationMaterial    bool   // set true while writing or reading variations
	lastReadFolder         string // used during wce parsing to remember context
	isObj                  bool   // true when a _obj suffix is found in path
	isChr                  bool   // true when a _chr suffix is found in path
	modelTags              []string
	maxMaterialHeads       map[string]int
	maxMaterialTextures    map[string]int
	tagIndexes             map[string]int // used when parsing to keep track of indexes
	FileName               string
	WorldDef               *WorldDef
	GlobalAmbientLightDef  *GlobalAmbientLightDef
	Version                uint32
	ActorDefs              []*ActorDef
	ActorInsts             []*ActorInst
	AmbientLights          []*AmbientLight
	BlitSpriteDefs         []*BlitSpriteDef
	DMSpriteDef2s          []*DMSpriteDef2
	DMSpriteDefs           []*DMSpriteDef
	DMTrackDef2s           []*DMTrackDef2
	HierarchicalSpriteDefs []*HierarchicalSpriteDef
	LightDefs              []*LightDef
	MaterialDefs           []*MaterialDef
	MaterialPalettes       []*MaterialPalette
	ParticleCloudDefs      []*ParticleCloudDef
	PointLights            []*PointLight
	PolyhedronDefs         []*PolyhedronDefinition
	Regions                []*Region
	RGBTrackDefs           []*RGBTrackDef
	SimpleSpriteDefs       []*SimpleSpriteDef
	Sprite2DDefs           []*Sprite2DDef
	Sprite3DDefs           []*Sprite3DDef
	TrackDefs              []*TrackDef
	TrackInstances         []*TrackInstance
	variationMaterialDefs  map[string][]*MaterialDef
	WorldTrees             []*WorldTree
	Zones                  []*Zone
	AniDefs                []*AniDef
	MdsDefs                []*MdsDef
	ModDefs                []*ModDef
	TerDefs                []*TerDef
	LayDefs                []*LayDef
}

type WldDefinitioner interface {
	Definition() string
	ToRaw(src *Wce, dst *raw.Wld) (int16, error)
	Write(token *AsciiWriteToken) error
}

func New(filename string) *Wce {
	isObj := strings.Contains(filename, "_obj")
	isChr := strings.Contains(filename, "_chr")

	return &Wce{
		FileName:              filename,
		isObj:                 isObj,
		isChr:                 isChr,
		maxMaterialHeads:      make(map[string]int),
		maxMaterialTextures:   make(map[string]int),
		variationMaterialDefs: make(map[string][]*MaterialDef),
		WorldDef:              &WorldDef{folders: []string{"world"}},
	}
}

// ByTag returns a instance by tag
func (wce *Wce) ByTag(tag string) WldDefinitioner {
	if tag == "" {
		return nil
	}
	if strings.HasSuffix(tag, "_SPRITE") {
		for _, sprite := range wce.SimpleSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
		for _, sprite := range wce.BlitSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}
	if strings.HasSuffix(tag, "_PCD") {
		for _, cloud := range wce.ParticleCloudDefs {
			if cloud.Tag == tag {
				return cloud
			}
		}
	}
	if strings.HasSuffix(tag, "_SPB") {
		for _, sprite := range wce.BlitSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}
	if strings.HasSuffix(tag, "_MDF") {
		for _, material := range wce.MaterialDefs {
			if material.Tag == tag {
				return material
			}
		}
	}
	if strings.HasSuffix(tag, "_MP") {
		for _, palette := range wce.MaterialPalettes {
			if palette.Tag == tag {
				return palette
			}
		}
	}
	if strings.HasSuffix(tag, "_DMSPRITEDEF") {
		for _, sprite := range wce.DMSpriteDef2s {
			if sprite.Tag == tag {
				return sprite
			}
		}
		for _, sprite := range wce.DMSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}
	if strings.HasSuffix(tag, "_DMTRACKDEF") {
		for _, track := range wce.DMTrackDef2s {
			if track.Tag == tag {
				return track
			}
		}
	}
	if strings.HasSuffix(tag, "_LIGHTDEF") {
		for _, light := range wce.LightDefs {
			if light.Tag == tag {
				return light
			}
		}
	}
	if strings.HasSuffix(tag, "_LDEF") {
		for _, light := range wce.LightDefs {
			if light.Tag == tag {
				return light
			}
		}
	}

	if strings.HasSuffix(tag, "_TRACKDEF") {
		for _, track := range wce.TrackDefs {
			if track.Tag == tag {
				return track
			}
		}
	}

	if strings.HasSuffix(tag, "_HS_DEF") {
		for _, sprite := range wce.HierarchicalSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}

	if strings.HasSuffix(tag, "_POLYHDEF") {
		for _, polyhedron := range wce.PolyhedronDefs {
			if polyhedron.Tag == tag {
				return polyhedron
			}
		}
	}

	if strings.HasSuffix(tag, "_DMT") {
		for _, track := range wce.RGBTrackDefs {
			if track.Tag == tag {
				return track
			}
		}
	}

	for _, sprite := range wce.Sprite3DDefs {
		if sprite.Tag == tag {
			return sprite
		}
	}
	for _, region := range wce.Regions {
		if region.Tag == tag {
			return region
		}
	}

	for _, actor := range wce.ActorDefs {
		if actor.Tag == tag {
			return actor
		}
	}

	for _, track := range wce.TrackInstances {
		if track.Tag == tag {
			return track
		}
	}

	for _, sprite := range wce.Sprite2DDefs {
		if sprite.Tag == tag {
			return sprite
		}
	}

	for _, sprite := range wce.SimpleSpriteDefs {
		if sprite.Tag == tag {
			return sprite
		}
		if strings.HasSuffix(sprite.Tag, "_SPRITE") && !strings.HasSuffix(tag, "_SPRITE") {
			if sprite.Tag == tag+"_SPRITE" {
				return sprite
			}
		}
	}
	return nil
}

// ByTagWithIndex returns a instance by tag with index included
func (wce *Wce) ByTagWithIndex(tag string, index int) WldDefinitioner {
	if tag == "" {
		return nil
	}

	if strings.HasSuffix(tag, "_DMSPRITEDEF") {
		for _, dmsprite := range wce.DMSpriteDef2s {
			if dmsprite.Tag == tag && dmsprite.TagIndex == index {
				return dmsprite
			}
		}
		for _, dmsprite := range wce.DMSpriteDefs {
			if dmsprite.Tag == tag && dmsprite.TagIndex == index {
				return dmsprite
			}
		}
	}

	if strings.HasSuffix(tag, "_TRACK") {
		for _, track := range wce.TrackInstances {
			if track.Tag == tag && track.TagIndex == index {
				return track
			}
		}
	}

	if strings.HasSuffix(tag, "_TRACKDEF") {
		for _, track := range wce.TrackDefs {
			if track.Tag == tag && track.TagIndex == index {
				return track
			}
		}
	}

	if strings.HasSuffix(tag, "_MDF") {
		for _, material := range wce.MaterialDefs {
			if material.Tag == tag && material.TagIndex == index {
				return material
			}
		}
	}

	if strings.HasSuffix(tag, "_SPRITE") {
		for _, sprite := range wce.SimpleSpriteDefs {
			if sprite.Tag == tag && sprite.TagIndex == index {
				return sprite
			}
		}
	}

	if strings.HasSuffix(tag, "_PCD") {
		for _, sprite := range wce.ParticleCloudDefs {
			if sprite.Tag == tag && sprite.TagIndex == index {
				return sprite
			}
		}
	}

	return nil
}

// NextTagIndex returns the next available index for a tag
func (wce *Wce) NextTagIndex(tag string) int {
	if tag == "" {
		return 0
	}

	_, ok := wce.tagIndexes[tag]
	if !ok {
		wce.tagIndexes[tag] = 0
		return 0
	}

	wce.tagIndexes[tag]++
	return wce.tagIndexes[tag]
}

func (wce *Wce) reset() {
	wce.GlobalAmbientLightDef = nil
	wce.lastReadFolder = ""
	wce.tagIndexes = make(map[string]int)
	wce.SimpleSpriteDefs = []*SimpleSpriteDef{}
	wce.MaterialDefs = []*MaterialDef{}
	wce.variationMaterialDefs = make(map[string][]*MaterialDef)
	wce.MaterialPalettes = []*MaterialPalette{}
	wce.DMSpriteDefs = []*DMSpriteDef{}
	wce.DMSpriteDef2s = []*DMSpriteDef2{}
	wce.ActorDefs = []*ActorDef{}
	wce.ActorInsts = []*ActorInst{}
	wce.LightDefs = []*LightDef{}
	wce.PointLights = []*PointLight{}
	wce.Sprite3DDefs = []*Sprite3DDef{}
	wce.TrackInstances = []*TrackInstance{}
	wce.TrackDefs = []*TrackDef{}
	wce.HierarchicalSpriteDefs = []*HierarchicalSpriteDef{}
	wce.PolyhedronDefs = []*PolyhedronDefinition{}
	wce.WorldTrees = []*WorldTree{}
	wce.Regions = []*Region{}
	wce.AmbientLights = []*AmbientLight{}
	wce.Zones = []*Zone{}
	wce.RGBTrackDefs = []*RGBTrackDef{}
	wce.ParticleCloudDefs = []*ParticleCloudDef{}
	wce.Sprite2DDefs = []*Sprite2DDef{}
	wce.MdsDefs = []*MdsDef{}
	wce.ModDefs = []*ModDef{}
	wce.TerDefs = []*TerDef{}
}
