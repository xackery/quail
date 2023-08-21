package wld

type Parserer interface {
	build(e *WLD) error
}

func (e *WLD) initPacks() map[int32]*encoderdecoder {
	return map[int32]*encoderdecoder{
		0x00: {name: "Default", decode: e.defaultRead, encode: e.defaultWrite},                                     // 0
		0x01: {name: "PaletteFile", decode: e.paletteFileRead, encode: e.paletteFileWrite},                         // 1
		0x02: {name: "UserData", decode: e.userDataRead, encode: e.userDataWrite},                                  // 2
		0x03: {name: "TextureList", decode: e.textureListRead, encode: e.textureListWrite},                         // 3
		0x04: {name: "Texture", decode: e.textureRead, encode: e.textureWrite},                                     // 4
		0x05: {name: "TextureRef", decode: e.textureRefRead, encode: e.textureRefWrite},                            // 5
		0x06: {name: "TwoDSpriteDef", decode: e.twoDSpriteDefRead, encode: e.twoDSpriteDefWrite},                   // 6
		0x07: {name: "TwoDSprite", decode: e.twoDSpriteRead, encode: e.twoDSpriteWrite},                            // 7
		0x08: {name: "ThreeDSpriteDef", decode: e.threeDSpriteDefRead, encode: e.threeDSpriteDefWrite},             // 8
		0x09: {name: "ThreeDSprite", decode: e.threeDSpriteRead, encode: e.threeDSpriteWrite},                      // 9
		0x0A: {name: "FourDSpriteDef", decode: e.fourDSpriteDefRead, encode: e.fourDSpriteDefWrite},                // 10
		0x0B: {name: "FourDSprite", decode: e.fourDSpriteRead, encode: e.fourDSpriteWrite},                         // 11
		0x0C: {name: "ParticleSpriteDef", decode: e.particleSpriteDefRead, encode: e.particleSpriteDefWrite},       // 12
		0x0D: {name: "ParticleSprite", decode: e.particleSpriteRead, encode: e.particleSpriteWrite},                // 13
		0x0E: {name: "CompositeSpriteDef", decode: e.compositeSpriteDefRead, encode: e.compositeSpriteDefWrite},    // 14
		0x0F: {name: "CompositeSprite", decode: e.compositeSpriteRead, encode: e.compositeSpriteWrite},             // 15
		0x10: {name: "skeletonTrackDef", decode: e.skeletonTrackDefRead, encode: e.skeletonTrackDefWrite},          // 16
		0x11: {name: "skeletonTrack", decode: e.skeletonTrackRead, encode: e.skeletonTrackWrite},                   // 17
		0x12: {name: "TrackDef", decode: e.trackDefRead, encode: e.trackDefWrite},                                  // 18
		0x13: {name: "Track", decode: e.trackRead, encode: e.trackWrite},                                           // 19
		0x14: {name: "Model", decode: e.modelRead, encode: e.modelWrite},                                           // 20
		0x15: {name: "ObjectLocation", decode: e.objectLocationRead, encode: e.objectLocationWrite},                // 21
		0x16: {name: "Sphere", decode: e.sphereRead, encode: e.sphereWrite},                                        // 22
		0x17: {name: "PolyhedronDef", decode: e.polyhedronDefRead, encode: e.polyhedronDefWrite},                   // 23
		0x18: {name: "Polyhedron", decode: e.polyhedronRead, encode: e.polyhedronWrite},                            // 24
		0x19: {name: "SphereListDef", decode: e.sphereListDefRead, encode: e.sphereListDefWrite},                   // 25
		0x1A: {name: "SphereList", decode: e.sphereListRead, encode: e.sphereListWrite},                            // 26
		0x1B: {name: "LightDef", decode: e.lightDefRead, encode: e.lightDefWrite},                                  // 27
		0x1C: {name: "Light", decode: e.lightRead, encode: e.lightWrite},                                           // 28
		0x1D: {name: "PointLightOld", decode: e.pointLightOldRead, encode: e.pointLightOldWrite},                   // 29
		0x1F: {name: "SoundDef", decode: e.soundDefRead, encode: e.soundDefWrite},                                  // 31
		0x20: {name: "Sound", decode: e.soundRead, encode: e.soundWrite},                                           // 32
		0x21: {name: "WorldTree", decode: e.worldTreeRead, encode: e.worldTreeWrite},                               // 33
		0x22: {name: "Region", decode: e.regionRead, encode: e.regionWrite},                                        // 34
		0x23: {name: "ActiveGeoRegion", decode: e.activeGeoRegionRead, encode: e.activeGeoRegionWrite},             // 35
		0x24: {name: "SkyRegion", decode: e.skyRegionRead, encode: e.skyRegionWrite},                               // 36
		0x25: {name: "DirectionalLightOld", decode: e.directionalLightOldRead, encode: e.directionalLightOldWrite}, // 37
		0x26: {name: "BlitSpriteDef", decode: e.blitSpriteDefRead, encode: e.blitSpriteDefWrite},                   // 38
		0x27: {name: "BlitSprite", decode: e.blitSpriteRead, encode: e.blitSpriteWrite},                            // 39
		0x28: {name: "PointLight", decode: e.pointLightRead, encode: e.pointLightWrite},                            // 40
		0x29: {name: "Zone", decode: e.zoneRead, encode: e.zoneWrite},                                              // 41
		0x2A: {name: "AmbientLight", decode: e.ambientLightRead, encode: e.ambientLightWrite},                      // 42
		0x2B: {name: "DirectionalLight", decode: e.directionalLightRead, encode: e.directionalLightWrite},          // 43
		0x2C: {name: "DmSpriteDef", decode: e.dmSpriteDefRead, encode: e.dmSpriteDefWrite},                         // 44
		0x2D: {name: "DmSprite", decode: e.dmSpriteRead, encode: e.dmSpriteWrite},                                  // 45
		0x2E: {name: "DmTrackDef", decode: e.dmTrackDefRead, encode: e.dmTrackDefWrite},                            // 46
		0x2F: {name: "DmTrack", decode: e.dmTrackRead, encode: e.dmTrackWrite},                                     // 47
		0x30: {name: "Material", decode: e.materialRead, encode: e.materialWrite},                                  // 48
		0x31: {name: "MaterialList", decode: e.materialListRead, encode: e.materialListWrite},                      // 49
		0x32: {name: "DmRGBTrackDef", decode: e.dmRGBTrackDefRead, encode: e.dmRGBTrackDefWrite},                   // 50
		0x33: {name: "DmRGBTrack", decode: e.dmRGBTrackRead, encode: e.dmRGBTrackWrite},                            // 51
		0x34: {name: "ParticleCloudDef", decode: e.particleCloudDefRead, encode: e.particleCloudDefWrite},          // 52
		0x35: {name: "First", decode: e.firstRead, encode: e.firstWrite},                                           // 53
		0x36: {name: "Mesh", decode: e.meshRead, encode: e.meshWrite},                                              // 54
		0x37: {name: "DmTrackDef2", decode: e.dmTrackDef2Read, encode: e.dmTrackDef2Write},                         // 55
	}
}
