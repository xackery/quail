// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package wld

import (
	"strings"

	"github.com/xackery/quail/raw"
)

var AsciiVersion = "v0.0.1"

// Wld is a struct representing a Wld file
type Wld struct {
	FileName               string
	GlobalAmbientLight     string
	Version                uint32
	SimpleSpriteDefs       []*SimpleSpriteDef
	MaterialDefs           []*MaterialDef
	MaterialPalettes       []*MaterialPalette
	DMSpriteDefs           []*DMSpriteDef
	DMSpriteInsts          []*DMSprite
	DMSpriteDef2s          []*DMSpriteDef2
	ActorDefs              []*ActorDef
	ActorInsts             []*ActorInst
	LightDefs              []*LightDef
	PointLights            []*PointLight
	Sprite3DDefs           []*Sprite3DDef
	TrackInstances         []*TrackInstance
	TrackDefs              []*TrackDef
	HierarchicalSpriteDefs []*HierarchicalSpriteDef
	PolyhedronDefs         []*PolyhedronDefinition
	WorldTrees             []*WorldTree
	Regions                []*Region
	AmbientLights          []*AmbientLight
	Zones                  []*Zone
	RGBTrackDefs           []*RGBTrackDef
	RGBTrackInsts          []*RGBTrack
}

type WldDefinitioner interface {
	Definition() string
	ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error)
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
	}
	if strings.HasSuffix(tag, "_LIGHTDEF") {
		for _, light := range wld.LightDefs {
			if light.Tag == tag {
				return light
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
	return nil
}
