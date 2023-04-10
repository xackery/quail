package fragment

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

//ref: https://github.com/SCMcLaughlin/p99-iksar-anim-oneclick/blob/master/src/structs_wld_frag.h

var (
	fragmentTypes = make(map[int32](func(r io.ReadSeeker) (archive.WldFragmenter, error)))
	names         = make(map[int32]string)
)

func New(fragIndex int32, r io.ReadSeeker) (archive.WldFragmenter, error) {
	loadFunc, ok := fragmentTypes[fragIndex]
	if !ok {
		return nil, fmt.Errorf("unknown frag index: %d 0x%x", fragIndex, fragIndex)
	}
	return loadFunc(r)
}

// SetNames set a global names prop each fragment can look up from to get name based on hash index
func SetNames(in map[int32]string) {
	names = in
}

func init() {
	// 0x01 DEFAULTPALETTEFILE
	fragmentTypes[1] = LoadDefaultPalette
	// 0x02 USERDATA
	fragmentTypes[2] = LoadUserData
	// 0x03 TextureImage aka BitmapName aka TextureImage
	fragmentTypes[3] = LoadBitmapInfo
	// 0x04 SimpleSpriteDef aka SimpleSpriteDef AKA BitmapInfo
	fragmentTypes[4] = LoadSimpleSpriteDef
	// 0x05 SimpleSprite aka BitmapInfoReference aka SIMPLESPRITEINST
	fragmentTypes[5] = LoadSimpleSprite
	// 0x06 TwoDSpriteDef aka 2DSpriteDef aka TwoDimensionalObject - Only found in gequip files. Seems to represent 2d sprites in the world (coins). 2DSPRITEDEF 0x06 Only found in gequip files. Seems to represent 2d sprites in the world (coins).
	fragmentTypes[6] = LoadTwoDSpriteDef
	// 0x07, TwoDSprite only found in gequip files. This fragment can be referenced by an actor fragment. Fragment7
	fragmentTypes[7] = LoadTwoDSprite
	// 0x08 3DSPRITEDEF Camera
	fragmentTypes[8] = LoadThreeDSpriteDef
	// 0x09 3DSprite Related to 3DSPRITEDEF, Maybe the 3DSPRITETAG to map the string. When I added a second ACTION to ACTORINST test, a second 0x09 was added. maybe from the 3DSPRITE %s. bodies are both compressed to 0 1 0 CameraReference
	fragmentTypes[9] = LoadThreeDSprite
	// 0xa 4DSPRITEDEF
	fragmentTypes[10] = LoadFourDSpriteDef
	// 0xb FUN_004079a0
	fragmentTypes[11] = LoadFourDSprite
	// 0xc PARTICLESPRITEDEF
	fragmentTypes[12] = LoadParticleSpriteDef
	// 0xd Unknown
	fragmentTypes[13] = LoadParticleSprite
	// 0xe COMPOSITESPRITEDEF
	fragmentTypes[14] = LoadCompositeSpriteDef
	// 0xf Unknown
	fragmentTypes[15] = LoadCompositeSprite
	// 0x10 SkeletonHierarchy aka HIERARCHICALSPRITEDEF
	fragmentTypes[16] = LoadHierarchialSpriteDef
	// 0x11 SkeletonHierarchyReference
	fragmentTypes[17] = LoadHierarchialSprite
	// 0x12 TrackDefinition
	fragmentTypes[18] = LoadTrackDef
	// 0x13 TrackInstance
	fragmentTypes[19] = LoadTrack
	// 0x14 Model aka Static or Animated Model aka ActorDef
	fragmentTypes[20] = LoadActorDef
	// 0x15 Actor aka ActorInstance aka ObjectLocation
	fragmentTypes[21] = LoadActor
	// 0x16 Sphere
	fragmentTypes[22] = LoadSphere
	// 0x17 POLYHEDRONDEFINITION aka PolyhedronDef
	fragmentTypes[23] = LoadPolyhedronDef
	// 0x18 POLYHEDRONDEFINITIONREF
	fragmentTypes[24] = LoadPolyhedron
	// 0x19 SPHERELISTDEFINITION
	fragmentTypes[25] = LoadSphereListDef
	// 0x1A SPHERELIST
	fragmentTypes[26] = LoadSphereList
	// 0x1B LightDefinition aka LightSource
	fragmentTypes[27] = LoadLightDef
	// 0x1C LightDefinitionReference
	fragmentTypes[28] = LoadLight
	// 0x1D Pointlight aka PointLightOld
	fragmentTypes[29] = LoadPointLightOld
	// 0x1E unknown 0x1E
	// 0x1F SoundDefinition
	fragmentTypes[31] = LoadSoundDef
	// 0x20 SoundInstance
	fragmentTypes[32] = LoadSound
	// 0x21 BspTree aka WorldTree
	fragmentTypes[33] = LoadWorldTree
	// 0x22 BspRegion aka Region
	fragmentTypes[34] = LoadRegion
	// 0x23 ACTIVEGEOMETRYREGION
	fragmentTypes[35] = LoadActiveGeoRegion
	// 0x24 SKYREGION
	fragmentTypes[36] = LoadSkyRegion
	// 0x25 DIRECTIONALLIGHT
	fragmentTypes[37] = LoadDirectionalLightOld
	// 0x26 BLITSPRITEDEFINITION
	fragmentTypes[38] = LoadBlitSpriteDef
	// 0x27 BLITSPRITEDEFINITIONREF
	fragmentTypes[39] = LoadBlitSprite
	// 0x28 PointLight aka LightInstance
	fragmentTypes[40] = LoadPointLight
	// 0x29 Zone aka BspRegionType
	fragmentTypes[41] = LoadZone
	// 0x2A AmbientLight
	fragmentTypes[42] = LoadAmbientLight
	// 0x2B DirectionalLight
	// 0x2C LegacyMesh ak DmSpriteDef
	fragmentTypes[44] = LoadDmSpriteDef
	// 0x2D MeshReference aka DmSprite
	fragmentTypes[45] = LoadDmSprite
	// 0x2E DmTrackDef aka Unknown0x2E
	fragmentTypes[46] = LoadDmTrackDef
	// 0x2F DmTrack aka MeshAnimatedVerticesReference
	fragmentTypes[47] = LoadDmTrack
	// 0x30 Material aka MaterialDef
	fragmentTypes[48] = LoadMaterialDef
	// 0x31 MaterialList aka MaterialPalette
	fragmentTypes[49] = LoadMaterialPalette
	// 0x32 VertexColor aka DmRGBTrackDef
	fragmentTypes[50] = LoadDmRGBTrackDef
	// 0x33 VertexColorsReference aka DmRGBTrack
	fragmentTypes[51] = LoadDmRGBTrack
	// 0x34 ParticleCloud aka ParticleCloudDef
	fragmentTypes[52] = LoadParticleCloudDef
	// 0x35 GlobalAmbientLight
	fragmentTypes[53] = LoadGlobalAmbientLightDef
	// 0x36 Mesh
	fragmentTypes[54] = LoadDmSpriteDef2
	// 0x37 MeshAnimatedVertices
	fragmentTypes[55] = LoadDmTrackDef2
}

func nameFromHashIndex(r io.ReadSeeker) (string, error) {
	name := ""
	var value int32
	err := binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return "", fmt.Errorf("read hash index: %w", err)
	}

	name, ok := names[-value]
	if !ok {
		return "", fmt.Errorf("hash 0x%x not found in names (len %d)", -value, len(names))
	}
	//dump.Hex(value, "name=(%s)", name)
	return name, nil
}

func dumpFragment(r io.ReadSeeker) {
	data := []byte{}
	for i := 0; i < 24; i++ {
		data = append(data, byte(0))
	}
	binary.Read(r, binary.LittleEndian, &data)
	fmt.Println(hex.Dump(data))
	r.Seek(0, io.SeekStart)
}
