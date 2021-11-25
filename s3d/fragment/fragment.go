package fragment

import (
	"fmt"
	"io"
)

// Fragment is what every fragment object type adheres to
type Fragment interface {
	// FragmentType identifies the fragment type
	FragmentType() string
}

var fragmentTypes = make(map[int32](func(r io.ReadSeeker) (Fragment, error)))

func Load(fragIndex int32, r io.ReadSeeker) (Fragment, error) {
	loadFunc, ok := fragmentTypes[fragIndex]
	if !ok {
		return nil, fmt.Errorf("unknown frag index: 0x%x", fragIndex)
	}
	return loadFunc(r)
}

func init() {
	//0x01 DEFAULTPALETTEFILE
	//0x02 USERDATA
	//0x03 SimpleSprite aka BitmapName
	fragmentTypes[0x03] = loadSimpleSprite
	//0x04 SimpleSpriteReference aka SimpleSpriteDef AKA BitmapInfo
	fragmentTypes[0x04] = loadSimpleSpriteReference
	//0x05 BitmapInfoReference aka SIMPLESPRITEINST
	fragmentTypes[0x05] = loadSimpleSpriteInstance
	//0x06 Only found in gequip files. Seems to represent 2d sprites in the world (coins). 2DSPRITEDEF 0x06 Only found in gequip files. Seems to represent 2d sprites in the world (coins).
	//0x07, only found in gequip files. This fragment can be referenced by an actor fragment. Fragment7
	//0x08 3DSPRITEDEF Camera
	fragmentTypes[0x08] = loadCamera
	//0x09 Related to 3DSPRITEDEF, Maybe the 3DSPRITETAG to map the string. When I added a second ACTION to ACTORINST test, a second 0x09 was added. maybe from the 3DSPRITE %s. bodies are both compressed to 0 1 0 CameraReference
	//0xa 4DSPRITEDEF
	//0xb FUN_004079a0
	//0xc PARTICLESPRITEDEF
	//0xd Unknown
	//0xe COMPOSITESPRITEDEF
	//0xf Unknown
	//0x10 SkeletonHierarchy aka HIERARCHICALSPRITEDEF
	//0x11 SkeletonHierarchyReference
	fragmentTypes[0x11] = loadSkeletonReference
	//0x12 TrackDefinition
	fragmentTypes[0x12] = loadTrack
	//0x13 TrackInstance
	fragmentTypes[0x13] = loadTrackReference
	//0x14 Actor aka ActorDef
	//0x15 ActorInstance aka ObjectInstance
	fragmentTypes[0x15] = loadObjectInstance
	//0x16 ACTORINSTANCETEST
	//0x17 POLYHEDRONDEFINITION
	//0x18 unknown
	//0x19 SPHERELISTDEFINITION
	//0x1A unknown
	//0x1B LightDefinition aka LightSource
	fragmentTypes[0x1B] = loadLightSource
	//0x1C LightDefinitionReference
	fragmentTypes[0x1C] = loadLightSourceReference
	//0x1D Pointlight
	//0x1E unknown 0x1E
	//0x1F SoundDefinition
	//0x20 SoundInstance
	//0x21 BspTree aka WorldTree
	//0x22 BspRegion aka Region
	fragmentTypes[0x22] = loadRegion
	//0x23 ACTIVEGEOMETRYREGION
	//0x24 SKYREGION
	//0x25 DIRECTIONALLIGHT
	//0x26 BLITSPRITEDEFINITION
	fragmentTypes[0x26] = loadParticleSprite
	fragmentTypes[0x27] = loadParticleSpriteReference
	//0x28 PointLight aka LightInstance
	fragmentTypes[0x28] = loadLightInstance
	//0x29 Zone aka BspRegionType
	//0x2A AmbientLight
	//0x2B DirectionalLight
	//0x2C LegacyMesh
	fragmentTypes[0x2C] = loadLegacyMesh
	//0x2D MeshReference
	fragmentTypes[0x2D] = loadMeshReference
	//0x30 Material
	fragmentTypes[0x30] = loadMaterial
	//0x31 MaterialList
	fragmentTypes[0x31] = loadMaterialList
	//0x32 VertexColors
	fragmentTypes[0x32] = loadVertexColor
	//0x33 VertexColorsReference
	fragmentTypes[0x33] = loadVertexColorReference
	//0x34 ParticleCloud
	fragmentTypes[0x34] = loadParticleCloud
	//0x35 GlobalAmbientLight
	fragmentTypes[0x35] = loadGlobalAmbientLight
	//0x36 Mesh
	fragmentTypes[0x35] = loadMesh
	//0x37 MeshAnimatedVertices
	fragmentTypes[0x36] = loadMeshAnimatedVertices
}
