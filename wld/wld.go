// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package wld

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
