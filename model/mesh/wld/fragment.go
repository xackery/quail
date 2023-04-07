package wld

import (
	"fmt"
	"io"

	"github.com/xackery/quail/model/mesh/wld/fragment"
	"github.com/xackery/quail/pfs/archive"
)

var fragmentTypes = make(map[int32](func(r io.ReadSeeker) (archive.WldFragmenter, error)))

func (e *WLD) ParseFragment(fragIndex int32, r io.ReadSeeker) (archive.WldFragmenter, error) {
	loadFunc, ok := fragmentTypes[fragIndex]
	if !ok {
		return nil, fmt.Errorf("unknown frag index: %d 0x%x", fragIndex, fragIndex)
	}
	return loadFunc(r)
}

func init() {
	//0x01 DEFAULTPALETTEFILE
	//0x02 USERDATA
	//0x03 SimpleSprite aka BitmapName
	fragmentTypes[3] = fragment.LoadSimpleSprite
	//0x04 SimpleSpriteReference aka SimpleSpriteDef AKA BitmapInfo
	fragmentTypes[4] = fragment.LoadSimpleSpriteReference
	//0x05 BitmapInfoReference aka SIMPLESPRITEINST
	fragmentTypes[5] = fragment.LoadSimpleSpriteInstance
	//0x06 Only found in gequip files. Seems to represent 2d sprites in the world (coins). 2DSPRITEDEF 0x06 Only found in gequip files. Seems to represent 2d sprites in the world (coins).
	//0x07, only found in gequip files. This fragment can be referenced by an actor fragment. Fragment7
	//0x08 3DSPRITEDEF Camera
	fragmentTypes[8] = fragment.LoadCamera
	//0x09 Related to 3DSPRITEDEF, Maybe the 3DSPRITETAG to map the string. When I added a second ACTION to ACTORINST test, a second 0x09 was added. maybe from the 3DSPRITE %s. bodies are both compressed to 0 1 0 CameraReference
	fragmentTypes[9] = fragment.LoadCameraReference
	//0xa 4DSPRITEDEF
	//0xb FUN_004079a0
	//0xc PARTICLESPRITEDEF
	//0xd Unknown
	//0xe COMPOSITESPRITEDEF
	//0xf Unknown
	//0x10 SkeletonHierarchy aka HIERARCHICALSPRITEDEF
	//0x11 SkeletonHierarchyReference
	fragmentTypes[17] = fragment.LoadSkeletonReference
	//0x12 TrackDefinition
	fragmentTypes[18] = fragment.LoadTrack
	//0x13 TrackInstance
	fragmentTypes[19] = fragment.LoadTrackReference
	//0x14 Actor aka ActorDef
	fragmentTypes[20] = fragment.LoadActor
	//0x15 ActorInstance aka ObjectInstance
	fragmentTypes[21] = fragment.LoadActorInstance
	//0x16 ACTORINSTANCETEST
	fragmentTypes[22] = fragment.LoadActorInstanceTest
	//0x17 POLYHEDRONDEFINITION
	//0x18 unknown
	//0x19 SPHERELISTDEFINITION
	//0x1A unknown
	//0x1B LightDefinition aka LightSource
	fragmentTypes[27] = fragment.LoadLightSource
	//0x1C LightDefinitionReference
	fragmentTypes[28] = fragment.LoadLightSourceReference
	//0x1D Pointlight
	//0x1E unknown 0x1E
	//0x1F SoundDefinition
	//0x20 SoundInstance
	//0x21 BspTree aka WorldTree
	fragmentTypes[33] = fragment.LoadWorldTree
	//0x22 BspRegion aka Region
	fragmentTypes[34] = fragment.LoadRegion
	//0x23 ACTIVEGEOMETRYREGION
	//0x24 SKYREGION
	//0x25 DIRECTIONALLIGHT
	//0x26 BLITSPRITEDEFINITION
	fragmentTypes[38] = fragment.LoadParticleSprite
	fragmentTypes[39] = fragment.LoadParticleSpriteReference
	//0x28 PointLight aka LightInstance
	fragmentTypes[40] = fragment.LoadLightInstance
	//0x29 Zone aka BspRegionType
	fragmentTypes[41] = fragment.LoadRegionType
	//0x2A AmbientLight
	fragmentTypes[42] = fragment.LoadAmbientLight
	//0x2B DirectionalLight
	//0x2C LegacyMesh
	fragmentTypes[44] = fragment.LoadLegacyMesh
	//0x2D MeshReference
	fragmentTypes[45] = fragment.LoadMeshReference
	//0x30 Material
	fragmentTypes[48] = fragment.LoadMaterial
	//0x31 MaterialList
	fragmentTypes[49] = fragment.LoadMaterialList
	//0x32 VertexColors
	fragmentTypes[50] = fragment.LoadVertexColor
	//0x33 VertexColorsReference
	fragmentTypes[51] = fragment.LoadVertexColorReference
	//0x34 ParticleCloud
	fragmentTypes[52] = fragment.LoadParticleCloud
	//0x35 GlobalAmbientLight
	fragmentTypes[53] = fragment.LoadGlobalAmbientLight
	//0x36 Mesh
	fragmentTypes[54] = fragment.LoadMesh
	//0x37 MeshAnimatedVertices
	fragmentTypes[55] = fragment.LoadMeshAnimatedVertices
}
