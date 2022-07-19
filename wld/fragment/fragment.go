package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
)

//ref: https://github.com/SCMcLaughlin/p99-iksar-anim-oneclick/blob/master/src/structs_wld_frag.h

var (
	fragmentTypes = make(map[int32](func(r io.ReadSeeker) (common.WldFragmenter, error)))
	names         = make(map[uint32]string)
)

func New(fragIndex int32, r io.ReadSeeker) (common.WldFragmenter, error) {
	loadFunc, ok := fragmentTypes[fragIndex]
	if !ok {
		return nil, fmt.Errorf("unknown frag index: %d 0x%x", fragIndex, fragIndex)
	}
	return loadFunc(r)
}

// SetNames set a global names prop each fragment can look up from to get name based on hash index
func SetNames(in map[uint32]string) {
	names = in
}

func init() {
	//0x01 DEFAULTPALETTEFILE
	//0x02 USERDATA
	//0x03 SimpleSprite aka BitmapName
	fragmentTypes[3] = LoadSimpleSprite
	//0x04 SimpleSpriteReference aka SimpleSpriteDef AKA BitmapInfo
	fragmentTypes[4] = LoadSimpleSpriteReference
	//0x05 BitmapInfoReference aka SIMPLESPRITEINST
	fragmentTypes[5] = LoadSimpleSpriteInstance
	//0x06 Only found in gequip files. Seems to represent 2d sprites in the world (coins). 2DSPRITEDEF 0x06 Only found in gequip files. Seems to represent 2d sprites in the world (coins).
	//0x07, only found in gequip files. This fragment can be referenced by an actor fragment. Fragment7
	//0x08 3DSPRITEDEF Camera
	fragmentTypes[8] = LoadCamera
	//0x09 Related to 3DSPRITEDEF, Maybe the 3DSPRITETAG to map the string. When I added a second ACTION to ACTORINST test, a second 0x09 was added. maybe from the 3DSPRITE %s. bodies are both compressed to 0 1 0 CameraReference
	fragmentTypes[9] = LoadCameraReference
	//0xa 4DSPRITEDEF
	//0xb FUN_004079a0
	//0xc PARTICLESPRITEDEF
	//0xd Unknown
	//0xe COMPOSITESPRITEDEF
	//0xf Unknown
	//0x10 SkeletonHierarchy aka HIERARCHICALSPRITEDEF
	fragmentTypes[16] = LoadSkeletonHierarchy
	//0x11 SkeletonHierarchyReference
	fragmentTypes[17] = LoadSkeletonReference
	//0x12 TrackDefinition
	fragmentTypes[18] = LoadTrack
	//0x13 TrackInstance
	fragmentTypes[19] = LoadTrackReference
	//0x14 Actor aka ActorDef
	fragmentTypes[20] = LoadActor
	//0x15 ActorInstance aka ObjectInstance
	fragmentTypes[21] = LoadActorInstance
	//0x16 ACTORINSTANCETEST
	fragmentTypes[22] = LoadActorInstanceTest
	//0x17 POLYHEDRONDEFINITION
	fragmentTypes[23] = LoadPolygonAnimation
	//0x18 POLYHEDRONDEFINITIONREF
	fragmentTypes[24] = LoadPolygonAnimationReference
	//0x19 SPHERELISTDEFINITION
	//0x1A unknown
	//0x1B LightDefinition aka LightSource
	fragmentTypes[27] = LoadLightSource
	//0x1C LightDefinitionReference
	fragmentTypes[28] = LoadLightSourceReference
	//0x1D Pointlight
	//0x1E unknown 0x1E
	//0x1F SoundDefinition
	//0x20 SoundInstance
	//0x21 BspTree aka WorldTree
	fragmentTypes[33] = LoadWorldTree
	//0x22 BspRegion aka Region
	fragmentTypes[34] = LoadRegion
	//0x23 ACTIVEGEOMETRYREGION
	//0x24 SKYREGION
	//0x25 DIRECTIONALLIGHT
	//0x26 BLITSPRITEDEFINITION
	fragmentTypes[38] = LoadParticleSprite
	fragmentTypes[39] = LoadParticleSpriteReference
	//0x28 PointLight aka LightInstance
	fragmentTypes[40] = LoadLightInstance
	//0x29 Zone aka BspRegionType
	fragmentTypes[41] = LoadRegionType
	//0x2A AmbientLight
	fragmentTypes[42] = LoadAmbientLight
	//0x2B DirectionalLight
	//0x2C LegacyMesh
	fragmentTypes[44] = LoadLegacyMesh
	//0x2D MeshReference
	fragmentTypes[45] = LoadMeshReference
	//0x30 Material
	fragmentTypes[48] = LoadMaterial
	//0x31 MaterialList
	fragmentTypes[49] = LoadMaterialList
	//0x32 VertexColors
	fragmentTypes[50] = LoadVertexColor
	//0x33 VertexColorsReference
	fragmentTypes[51] = LoadVertexColorReference
	//0x34 ParticleCloud
	fragmentTypes[52] = LoadParticleCloud
	//0x35 GlobalAmbientLight
	fragmentTypes[53] = LoadGlobalAmbientLight
	//0x36 Mesh
	fragmentTypes[54] = LoadMesh
	//0x37 MeshAnimatedVertices
	fragmentTypes[55] = LoadMeshAnimatedVertices
}

func nameFromHashIndex(r io.ReadSeeker) (string, error) {
	name := ""
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return "", fmt.Errorf("read hash index: %w", err)
	}
	name, ok := names[value]
	if !ok {
		return "", fmt.Errorf("hash 0x%x not found in names (%d)", value, len(names))
	}
	dump.Hex(value, "name=(%s)", name)
	return name, nil
}
