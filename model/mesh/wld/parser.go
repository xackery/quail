package wld

type parserer interface {
	build(e *WLD) error
}

func (e *WLD) initParsers() map[int32]*decoder {
	parsers := map[int32]*decoder{
		0x00: {name: "Default", parse: e.defaultRead},                         // 0
		0x01: {name: "PaletteFile", parse: e.paletteFileRead},                 // 1
		0x02: {name: "UserData", parse: e.userDataRead},                       // 2
		0x03: {name: "TextureList", parse: e.textureListRead},                 // 3
		0x04: {name: "Texture", parse: e.textureRead},                         // 4
		0x05: {name: "TextureRef", parse: e.textureRefRead},                   // 5
		0x06: {name: "TwoDSpriteDef", parse: e.twoDSpriteDefRead},             // 6
		0x07: {name: "TwoDSprite", parse: e.twoDSpriteRead},                   // 7
		0x08: {name: "ThreeDSpriteDef", parse: e.threeDSpriteDefRead},         // 8
		0x09: {name: "ThreeDSprite", parse: e.threeDSpriteRead},               // 9
		0x0A: {name: "FourDSpriteDef", parse: e.fourDSpriteDefRead},           // 10
		0x0B: {name: "FourDSprite", parse: e.fourDSpriteRead},                 // 11
		0x0C: {name: "ParticleSpriteDef", parse: e.particleSpriteDefRead},     // 12
		0x0D: {name: "ParticleSprite", parse: e.particleSpriteRead},           // 13
		0x0E: {name: "CompositeSpriteDef", parse: e.compositeSpriteDefRead},   // 14
		0x0F: {name: "CompositeSprite", parse: e.compositeSpriteRead},         // 15
		0x10: {name: "skeletonTrackDef", parse: e.skeletonTrackDefRead},       // 16
		0x11: {name: "skeletonTrack", parse: e.skeletonTrackRead},             // 17
		0x12: {name: "TrackDef", parse: e.trackDefRead},                       // 18
		0x13: {name: "Track", parse: e.trackRead},                             // 19
		0x14: {name: "Model", parse: e.modelRead},                             // 20
		0x15: {name: "ObjectLocation", parse: e.objectLocationRead},           // 21
		0x16: {name: "Sphere", parse: e.sphereRead},                           // 22
		0x17: {name: "PolyhedronDef", parse: e.polyhedronDefRead},             // 23
		0x18: {name: "Polyhedron", parse: e.polyhedronRead},                   // 24
		0x19: {name: "SphereListDef", parse: e.sphereListDefRead},             // 25
		0x1A: {name: "SphereList", parse: e.sphereListRead},                   // 26
		0x1B: {name: "LightDef", parse: e.lightDefRead},                       // 27
		0x1C: {name: "Light", parse: e.lightRead},                             // 28
		0x1D: {name: "PointLightOld", parse: e.pointLightOldRead},             // 29
		0x1F: {name: "SoundDef", parse: e.soundDefRead},                       // 31
		0x20: {name: "Sound", parse: e.soundRead},                             // 32
		0x21: {name: "WorldTree", parse: e.worldTreeRead},                     // 33
		0x22: {name: "Region", parse: e.regionRead},                           // 34
		0x23: {name: "ActiveGeoRegion", parse: e.activeGeoRegionRead},         // 35
		0x24: {name: "SkyRegion", parse: e.skyRegionRead},                     // 36
		0x25: {name: "DirectionalLightOld", parse: e.directionalLightOldRead}, // 37
		0x26: {name: "BlitSpriteDef", parse: e.blitSpriteDefRead},             // 38
		0x27: {name: "BlitSprite", parse: e.blitSpriteRead},                   // 39
		0x28: {name: "PointLight", parse: e.pointLightRead},                   // 40
		0x29: {name: "Zone", parse: e.zoneRead},                               // 41
		0x2A: {name: "AmbientLight", parse: e.ambientLightRead},               // 42
		0x2B: {name: "DirectionalLight", parse: e.directionalLightRead},       // 43
		0x2C: {name: "DmSpriteDef", parse: e.dmSpriteDefRead},                 // 44
		0x2D: {name: "DmSprite", parse: e.dmSpriteRead},                       // 45
		0x2E: {name: "DmTrackDef", parse: e.dmTrackDefRead},                   // 46
		0x2F: {name: "DmTrack", parse: e.dmTrackRead},                         // 47
		0x30: {name: "Material", parse: e.materialRead},                       // 48
		0x31: {name: "MaterialList", parse: e.materialListRead},               // 49
		0x32: {name: "DmRGBTrackDef", parse: e.dmRGBTrackDefRead},             // 50
		0x33: {name: "DmRGBTrack", parse: e.dmRGBTrackRead},                   // 51
		0x34: {name: "ParticleCloudDef", parse: e.particleCloudDefRead},       // 52
		0x35: {name: "First", parse: e.firstRead},                             // 53
		0x36: {name: "mesh", parse: e.meshRead},                               // 54
		0x37: {name: "DmTrackDef2", parse: e.dmTrackDef2Read},                 // 55
	}

	return parsers
}
