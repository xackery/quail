package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

var (
	fragNames = map[int]string{
		0x00: "Default",
		0x01: "DefaultPaletteFile",
		0x02: "UserData",
		0x03: "BMInfo",
		0x04: "SimpleSpriteDef",
		0x05: "SimpleSprite",
		0x06: "Sprite2DDef",
		0x07: "Sprite2D",
		0x08: "Sprite3DDef",
		0x09: "Sprite3D",
		0x0A: "Sprite4DDef",
		0x0B: "Sprite4D",
		0x0C: "ParticleSpriteDef",
		0x0D: "ParticleSprite",
		0x0E: "CompositeSpriteDef",
		0x0F: "CompositeSprite",
		0x10: "HierarchialSpriteDef",
		0x11: "HierarchialSprite",
		0x12: "TrackDef",
		0x13: "Track",
		0x14: "ActorDef",
		0x15: "Actor",
		0x16: "Sphere",
		0x17: "PolyhedronDef",
		0x18: "Polyhedron",
		0x19: "SphereListDef",
		0x1A: "SphereList",
		0x1B: "LightDef",
		0x1C: "Light",
		0x1D: "PointLightOld",
		0x1E: "PointLightOldDef",
		0x1F: "Sound",
		0x20: "SoundDef",
		0x21: "WorldTree",
		0x22: "Region",
		0x23: "ActiveGeoRegion",
		0x24: "SkyRegion",
		0x25: "DirectionalLightOld",
		0x26: "BlitSpriteDef",
		0x27: "BlitSprite",
		0x28: "PointLight",
		0x29: "Zone",
		0x2A: "AmbientLight",
		0x2B: "DirectionalLight",
		0x2C: "DMSpriteDef",
		0x2D: "DMSprite",
		0x2E: "DMTrackDef",
		0x2F: "DMTrack",
		0x30: "MaterialDef",
		0x31: "MaterialPalette",
		0x32: "DmRGBTrackDef",
		0x33: "DmRGBTrack",
		0x34: "ParticleCloudDef",
		0x35: "GlobalAmbientLightDef",
		0x36: "DmSpriteDef2",
		0x37: "DmTrackDef2",
	}
)

const (
	FragCodeDefault               = 0x00 // 0
	FragCodeDefaultPaletteFile    = 0x01 // 1
	FragCodeUserData              = 0x02 // 2
	FragCodeBMInfo                = 0x03 // 3 BMINFO (property of SIMPLESPRITEDEFINITION)
	FragCodeSimpleSpriteDef       = 0x04 // 4 SIMPLESPRITEDEF
	FragCodeSimpleSprite          = 0x05 // 5 SIMPLESPRITE (property of BLITSPRITEDEFINITION, MATERIALDEFINITION)
	FragCodeSprite2DDef           = 0x06 // 6
	FragCodeSprite2D              = 0x07 // 7
	FragCodeSprite3DDef           = 0x08 // 8 SPRITE3DDEF
	FragCodeSprite3D              = 0x09 // 9 3DSPRITE (property of ACTORDEF)
	FragCodeSprite4DDef           = 0x0A // 10
	FragCodeSprite4D              = 0x0B // 11
	FragCodeParticleSpriteDef     = 0x0C // 12 PARTICLESPRITEDEF
	FragCodeParticleSprite        = 0x0D // 13
	FragCodeCompositeSpriteDef    = 0x0E // 14
	FragCodeCompositeSprite       = 0x0F // 15
	FragCodeHierarchialSpriteDef  = 0x10 // 16 HIERARCHICALSPRITEDEF
	FragCodeHierarchialSprite     = 0x11 // 17
	FragCodeTrackDef              = 0x12 // 18 TRACKDEFINITION
	FragCodeTrack                 = 0x13 // 19 TRACKINSTANCE
	FragCodeActorDef              = 0x14 // 20 ACTORDEF
	FragCodeActor                 = 0x15 // 21 ACTORINST
	FragCodeSphere                = 0x16 // 22 SPHERE (may be child of region?)
	FragCodePolyhedronDef         = 0x17 // 23 POLYHEDRONDEFINITION
	FragCodePolyhedron            = 0x18 // 24
	FragCodeSphereListDef         = 0x19 // 25
	FragCodeSphereList            = 0x1A // 26
	FragCodeLightDef              = 0x1B // 27 LIGHTDEFINITION
	FragCodeLight                 = 0x1C // 28
	FragCodePointLightOld         = 0x1D // 29
	FragCodePointLightOldDef      = 0x1E // 30
	FragCodeSound                 = 0x1F // 31
	FragCodeSoundDef              = 0x20 // 32
	FragCodeWorldTree             = 0x21 // 33 WORLDTREE
	FragCodeRegion                = 0x22 // 34 REGION
	FragCodeActiveGeoRegion       = 0x23 // 35
	FragCodeSkyRegion             = 0x24 // 36
	FragCodeDirectionalLightOld   = 0x25 // 37
	FragCodeBlitSpriteDef         = 0x26 // 38 BLITSPRITEDEFINITION
	FragCodeBlitSprite            = 0x27 // 39
	FragCodePointLight            = 0x28 // 40 POINTLIGHT
	FragCodeZone                  = 0x29 // 41 ZONE
	FragCodeAmbientLight          = 0x2A // 42 AMBIENTLIGHT
	FragCodeDirectionalLight      = 0x2B // 43
	FragCodeDMSpriteDef           = 0x2C // 44
	FragCodeDMSprite              = 0x2D // 45 DMSPRITE (property of HIERARCHICALSPRITEDEF)
	FragCodeDMTrackDef            = 0x2E // 46 TRACKDEFINITION
	FragCodeDMTrack               = 0x2F // 47 TRACKINSTANCE
	FragCodeMaterialDef           = 0x30 // 48 MATERIALDEFINITION
	FragCodeMaterialPalette       = 0x31 // 49 MATERIALPALETTE
	FragCodeDmRGBTrackDef         = 0x32 // 50
	FragCodeDmRGBTrack            = 0x33 // 51
	FragCodeParticleCloudDef      = 0x34 // 52
	FragCodeGlobalAmbientLightDef = 0x35 // 53
	FragCodeDmSpriteDef2          = 0x36 // 54 DMSPRITEDEF2
	FragCodeDmTrackDef2           = 0x37 // 55
)

// FragName returns the name of a fragment
func FragName(fragCode int) string {
	name, ok := fragNames[fragCode]
	if ok {
		return name
	}
	return fmt.Sprintf("unknownFrag%d", fragCode)
}

// FragIndex returns the index of a fragment
func FragIndex(name string) int {
	for k, v := range fragNames {
		if v == name {
			return k
		}
	}
	return -1
}

// NewFrag takes a reader, analyzes the first 4 bytes, and returns a new fragment struct based on it
func NewFrag(r io.ReadSeeker) model.FragmentReadWriter {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	fragCode := dec.Int32()
	err := dec.Error()
	if err != nil {
		return nil
	}
	//r.Seek(0, io.SeekStart)
	switch fragCode {
	case FragCodeDefault:
		return &WldFragDefault{}
	case FragCodeDefaultPaletteFile:
		return &WldFragDefaultPaletteFile{}
	case FragCodeUserData:
		return &WldFragUserData{}
	case FragCodeBMInfo:
		return &WldFragBMInfo{}
	case FragCodeSimpleSpriteDef:
		return &WldFragSimpleSpriteDef{}
	case FragCodeSimpleSprite:
		return &WldFragSimpleSprite{}
	case FragCodeSprite2DDef:
		return &WldFragSprite2DDef{}
	case FragCodeSprite2D:
		return &WldFragSprite2D{}
	case FragCodeSprite3DDef:
		return &WldFragSprite3DDef{}
	case FragCodeSprite3D:
		return &WldFragSprite3D{}
	case FragCodeSprite4DDef:
		return &WldFragSprite4DDef{}
	case FragCodeSprite4D:
		return &WldFragSprite4D{}
	case FragCodeParticleSpriteDef:
		return &WldFragParticleSpriteDef{}
	case FragCodeParticleSprite:
		return &WldFragParticleSprite{}
	case FragCodeCompositeSpriteDef:
		return &WldFragCompositeSpriteDef{}
	case FragCodeCompositeSprite:
		return &WldFragCompositeSprite{}
	case FragCodeHierarchialSpriteDef:
		return &WldFragHierarchialSpriteDef{}
	case FragCodeHierarchialSprite:
		return &WldFragHierarchialSprite{}
	case FragCodeTrackDef:
		return &WldFragTrackDef{}
	case FragCodeTrack:
		return &WldFragTrack{}
	case FragCodeActorDef:
		return &WldFragActorDef{}
	case FragCodeActor:
		return &WldFragActor{}
	case FragCodeSphere:
		return &WldFragSphere{}
	case FragCodePolyhedronDef:
		return &WldFragPolyhedronDef{}
	case FragCodePolyhedron:
		return &WldFragPolyhedron{}
	case FragCodeSphereListDef:
		return &WldFragSphereListDef{}
	case FragCodeSphereList:
		return &WldFragSphereList{}
	case FragCodeLight:
		return &WldFragLight{}
	case FragCodeLightDef:
		return &WldFragLightDef{}
	case FragCodePointLightOld:
		return &WldFragPointLightOld{}
	case FragCodePointLightOldDef:
		return &WldFragPointLightOldDef{}
	case FragCodeSound:
		return &WldFragSound{}
	case FragCodeSoundDef:
		return &WldFragSoundDef{}
	case FragCodeWorldTree:
		return &WldFragWorldTree{}
	case FragCodeRegion:
		return &WldFragRegion{}
	case FragCodeActiveGeoRegion:
		return &WldFragActiveGeoRegion{}
	case FragCodeSkyRegion:
		return &WldFragSkyRegion{}
	case FragCodeDirectionalLightOld:
		return &WldFragDirectionalLightOld{}
	case FragCodeBlitSpriteDef:
		return &WldFragBlitSpriteDef{}
	case FragCodeBlitSprite:
		return &WldFragBlitSprite{}
	case FragCodePointLight:
		return &WldFragPointLight{}
	case FragCodeZone:
		return &WldFragZone{}
	case FragCodeAmbientLight:
		return &WldFragAmbientLight{}
	case FragCodeDirectionalLight:
		return &WldFragDirectionalLight{}
	case FragCodeDMSpriteDef:
		return &WldFragDMSpriteDef{}
	case FragCodeDMSprite:
		return &WldFragDMSprite{}
	case FragCodeDMTrackDef:
		return &WldFragDMTrackDef{}
	case FragCodeDMTrack:
		return &WldFragDMTrack{}
	case FragCodeMaterialDef:
		return &WldFragMaterialDef{}
	case FragCodeMaterialPalette:
		return &WldFragMaterialPalette{}
	case FragCodeDmRGBTrackDef:
		return &WldFragDmRGBTrackDef{}
	case FragCodeDmRGBTrack:
		return &WldFragDmRGBTrack{}
	case FragCodeParticleCloudDef:
		return &WldFragParticleCloudDef{}
	case FragCodeGlobalAmbientLightDef:
		return &WldFragGlobalAmbientLightDef{}
	case FragCodeDmSpriteDef2:
		return &WldFragDmSpriteDef2{}
	case FragCodeDmTrackDef2:
		return &WldFragDmTrackDef2{}
	}
	return nil
}
