// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package wld

import (
	"strings"

	"github.com/xackery/quail/raw"
)

var AsciiVersion = "v0.0.1"

// Wld is a struct representing a Wld file
type Wld struct {
	isVariationMaterial    bool   // set true while writing or reading variations
	lastReadModelTag       string // last model tag read
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
}

type WldDefinitioner interface {
	Definition() string
	ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error)
	Write(token *AsciiWriteToken) error
}

// ByTag returns a instance by tag
func (wld *Wld) ByTag(tag string) WldDefinitioner {
	if tag == "" {
		return nil
	}
	if strings.HasSuffix(tag, "_SPRITE") {
		for _, sprite := range wld.SimpleSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}
	if strings.HasSuffix(tag, "_PCD") {
		for _, cloud := range wld.ParticleCloudDefs {
			if cloud.Tag == tag {
				return cloud
			}
		}
	}
	if strings.HasSuffix(tag, "_SPB") {
		for _, sprite := range wld.ParticleCloudDefs {
			if sprite.BlitSpriteDefTag == tag {
				return sprite
			}
		}
	}
	if strings.HasSuffix(tag, "_MDF") {
		for _, material := range wld.MaterialDefs {
			if material.Tag == tag {
				return material
			}
		}
	}
	if strings.HasSuffix(tag, "_MP") {
		for _, palette := range wld.MaterialPalettes {
			if palette.Tag == tag {
				return palette
			}
		}
	}
	if strings.HasSuffix(tag, "_DMSPRITEDEF") {
		for _, sprite := range wld.DMSpriteDef2s {
			if sprite.Tag == tag {
				return sprite
			}
		}
		for _, sprite := range wld.DMSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}
	if strings.HasSuffix(tag, "_LIGHTDEF") {
		for _, light := range wld.LightDefs {
			if light.Tag == tag {
				return light
			}
		}
	}
	if strings.HasSuffix(tag, "_LDEF") {
		for _, light := range wld.LightDefs {
			if light.Tag == tag {
				return light
			}
		}
	}

	if strings.HasSuffix(tag, "_TRACKDEF") {
		for _, track := range wld.TrackDefs {
			if track.Tag == tag {
				return track
			}
		}
	}

	if strings.HasSuffix(tag, "_HS_DEF") {
		for _, sprite := range wld.HierarchicalSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}

	if strings.HasSuffix(tag, "_POLYHDEF") {
		for _, polyhedron := range wld.PolyhedronDefs {
			if polyhedron.Tag == tag {
				return polyhedron
			}
		}
	}

	if strings.HasSuffix(tag, "_DMT") {
		for _, track := range wld.RGBTrackDefs {
			if track.Tag == tag {
				return track
			}
		}
	}

	for _, sprite := range wld.Sprite3DDefs {
		if sprite.Tag == tag {
			return sprite
		}
	}
	for _, region := range wld.Regions {
		if region.Tag == tag {
			return region
		}
	}

	for _, actor := range wld.ActorDefs {
		if actor.Tag == tag {
			return actor
		}
	}

	for _, track := range wld.TrackInstances {
		if track.Tag == tag {
			return track
		}
	}

	for _, sprite := range wld.Sprite2DDefs {
		if sprite.Tag == tag {
			return sprite
		}
	}

	for _, sprite := range wld.SimpleSpriteDefs {
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
func (wld *Wld) ByTagWithIndex(tag string, index int) WldDefinitioner {
	if tag == "" {
		return nil
	}

	if strings.HasSuffix(tag, "_TRACK") {
		for _, track := range wld.TrackInstances {
			if track.Tag == tag && track.TagIndex == index {
				return track
			}
		}
	}

	if strings.HasSuffix(tag, "_TRACKDEF") {
		for _, track := range wld.TrackDefs {
			if track.Tag == tag && track.TagIndex == index {
				return track
			}
		}
	}

	if strings.HasSuffix(tag, "_MDF") {
		for _, material := range wld.MaterialDefs {
			if material.Tag == tag && material.Variation == index {
				return material
			}
		}
	}

	if strings.HasSuffix(tag, "_SPRITE") {
		for _, sprite := range wld.SimpleSpriteDefs {
			if sprite.Tag != tag {
				continue
			}
			if sprite.Variation != index {
				continue
			}
			return sprite
		}
	}

	return nil
}

// NextTagIndex returns the next available index for a tag
func (wld *Wld) NextTagIndex(tag string) int {
	if tag == "" {
		return 0
	}

	_, ok := wld.tagIndexes[tag]
	if !ok {
		wld.tagIndexes[tag] = 0
		return 0
	}

	wld.tagIndexes[tag]++
	return wld.tagIndexes[tag]
}

func (wld *Wld) reset() {
	wld.GlobalAmbientLightDef = nil
	wld.lastReadModelTag = ""
	wld.tagIndexes = make(map[string]int)
	wld.SimpleSpriteDefs = []*SimpleSpriteDef{}
	wld.MaterialDefs = []*MaterialDef{}
	wld.variationMaterialDefs = make(map[string][]*MaterialDef)
	wld.MaterialPalettes = []*MaterialPalette{}
	wld.DMSpriteDefs = []*DMSpriteDef{}
	wld.DMSpriteDef2s = []*DMSpriteDef2{}
	wld.ActorDefs = []*ActorDef{}
	wld.ActorInsts = []*ActorInst{}
	wld.LightDefs = []*LightDef{}
	wld.PointLights = []*PointLight{}
	wld.Sprite3DDefs = []*Sprite3DDef{}
	wld.TrackInstances = []*TrackInstance{}
	wld.TrackDefs = []*TrackDef{}
	wld.HierarchicalSpriteDefs = []*HierarchicalSpriteDef{}
	wld.PolyhedronDefs = []*PolyhedronDefinition{}
	wld.WorldTrees = []*WorldTree{}
	wld.Regions = []*Region{}
	wld.AmbientLights = []*AmbientLight{}
	wld.Zones = []*Zone{}
	wld.RGBTrackDefs = []*RGBTrackDef{}
	wld.ParticleCloudDefs = []*ParticleCloudDef{}
	wld.Sprite2DDefs = []*Sprite2DDef{}
}
