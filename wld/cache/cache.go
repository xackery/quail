package cache

import (
	"github.com/xackery/quail/model"
)

type CacheManager struct {
	FileName              string
	GlobalAmbientLight    string
	Version               uint32
	Bitmaps               []*Bitmap
	SimpleSpriteDefs      []*SimpleSpriteDef
	SpriteInstances       []*SpriteInstance
	Particles             []*Particle
	ParticleInstances     []*ParticleInstance
	MaterialDefs          []*MaterialDef
	MaterialPalettes      []*MaterialPalette
	DmSpriteDef2s         []*DmSpriteDef2
	MeshInstances         []*MeshInstance
	AlternateMeshes       []*AlternateMesh
	ActorDefs             []*ActorDef
	ActorInsts            []*ActorInst
	Animations            []*Animation
	AnimationInstances    []*AnimationInstance
	Skeletons             []*Skeleton
	SkeletonInstances     []*SkeletonInstance
	LightDefs             []*LightDef
	LightInstances        []*LightInstance
	AmbientLightInstances []*AmbientLightInstance
	PointLights           []*PointLight
	Sprite3DDefs          []*Sprite3DDef
	CameraInstances       []*CameraInstance
	BspTrees              []*BspTree
	Regions               []*Region
	RegionInstances       []*RegionInstance
	Spheres               []*Sphere
	PolyhedronDefs        []*PolyhedronDef
	PolyhedronInstances   []*PolyhedronInstance
}

func (cm *CacheManager) Close() {
	cm.Bitmaps = []*Bitmap{}
	cm.SimpleSpriteDefs = []*SimpleSpriteDef{}
	cm.SpriteInstances = []*SpriteInstance{}
	cm.Particles = []*Particle{}
	cm.ParticleInstances = []*ParticleInstance{}
	cm.MaterialDefs = []*MaterialDef{}
	cm.MaterialPalettes = []*MaterialPalette{}
	cm.DmSpriteDef2s = []*DmSpriteDef2{}
	cm.MeshInstances = []*MeshInstance{}
	cm.AlternateMeshes = []*AlternateMesh{}
	cm.ActorDefs = []*ActorDef{}
	cm.ActorInsts = []*ActorInst{}
	cm.Animations = []*Animation{}
	cm.AnimationInstances = []*AnimationInstance{}
	cm.Skeletons = []*Skeleton{}
	cm.SkeletonInstances = []*SkeletonInstance{}
	cm.LightDefs = []*LightDef{}
	cm.LightInstances = []*LightInstance{}
	cm.AmbientLightInstances = []*AmbientLightInstance{}
	cm.PointLights = []*PointLight{}
	cm.Sprite3DDefs = []*Sprite3DDef{}
	cm.CameraInstances = []*CameraInstance{}
	cm.BspTrees = []*BspTree{}
	cm.Regions = []*Region{}
	cm.RegionInstances = []*RegionInstance{}
	cm.Spheres = []*Sphere{}

}

// Bitmap is a struct representing a material
type Bitmap struct {
	fragID          uint32
	Tag             string
	Textures        []string
	SimpleSpriteDef SimpleSpriteDef
	SimpleSprite    SpriteInstance
}

func (cm *CacheManager) bitmapByFragID(fragID uint32) *Bitmap {
	for _, bitmap := range cm.Bitmaps {
		if bitmap.fragID == fragID {
			return bitmap
		}
	}
	return nil
}

func (cm *CacheManager) bitmapByTag(tag string) *Bitmap {
	for _, bitmap := range cm.Bitmaps {
		if bitmap.Tag == tag {
			return bitmap
		}
	}
	return nil
}

type SimpleSpriteDef struct {
	fragID       uint32
	Tag          string
	Flags        uint32
	CurrentFrame int32
	Sleep        uint32
	BMInfos      [][2]string
}

func (cm *CacheManager) spriteByFragID(fragID uint32) *SimpleSpriteDef {
	for _, sprite := range cm.SimpleSpriteDefs {
		if sprite.fragID == fragID {
			return sprite
		}
	}
	return nil
}

func (cm *CacheManager) spriteByTag(tag string) *SimpleSpriteDef {
	for _, sprite := range cm.SimpleSpriteDefs {
		if sprite.Tag == tag {
			return sprite
		}
	}
	return nil
}

type SpriteInstance struct {
	fragID uint32
	Tag    string
	Flags  uint32
	Sprite string
}

func (cm *CacheManager) spriteInstanceByFragID(fragID uint32) *SpriteInstance {
	for _, spriteInstance := range cm.SpriteInstances {
		if spriteInstance.fragID == fragID {
			return spriteInstance
		}
	}
	return nil
}

func (cm *CacheManager) spriteInstanceByTag(tag string) *SpriteInstance {
	for _, spriteInstance := range cm.SpriteInstances {
		if spriteInstance.Tag == tag {
			return spriteInstance
		}
	}
	return nil
}

// Particle is also known as BlitSpriteDef
type Particle struct {
	fragID         uint32
	Tag            string
	Flags          uint32
	SpriteTag      string
	Unknown        int32
	spriteInstance *SpriteInstance
}

func (cm *CacheManager) particleByFragID(fragID uint32) *Particle {
	for _, particle := range cm.Particles {
		if particle.fragID == fragID {
			return particle
		}
	}
	return nil
}

func (cm *CacheManager) particleByTag(tag string) *Particle {
	for _, particle := range cm.Particles {
		if particle.Tag == tag {
			return particle
		}
	}
	return nil
}

type ParticleInstance struct {
	fragID                uint32
	Tag                   string
	Unk1                  uint32
	Unk2                  uint32
	ParticleMovement      uint32
	Flags                 uint32 //Flag 1, High Opacity, Flag 3, Follows Item
	SimultaneousParticles uint32
	Unk6                  uint32
	Unk7                  uint32
	Unk8                  uint32
	Unk9                  uint32
	Unk10                 uint32
	SpawnRadius           float32
	SpawnAngle            float32
	SpawnLifespan         uint32
	SpawnVelocity         float32
	SpawnNormalZ          float32
	SpawnNormalX          float32
	SpawnNormalY          float32
	SpawnRate             uint32
	SpawnScale            float32
	Color                 model.RGBA
	particle              *Particle
}

func (cm *CacheManager) particleInstanceByFragID(fragID uint32) *ParticleInstance {
	for _, particleInstance := range cm.ParticleInstances {
		if particleInstance.fragID == fragID {
			return particleInstance
		}
	}
	return nil
}

func (cm *CacheManager) particleInstanceByTag(tag string) *ParticleInstance {
	for _, particleInstance := range cm.ParticleInstances {
		if particleInstance.Tag == tag {
			return particleInstance
		}
	}
	return nil
}

// MaterialDef is a struct representing a material
type MaterialDef struct {
	fragID               uint32
	Tag                  string
	Flags                uint32
	RenderMethod         string
	RGBPen               [4]uint8
	Brightness           float32
	ScaledAmbient        float32
	SimpleSpriteInstTag  string
	SimpleSpriteInstFlag uint32
	spriteInstance       *SpriteInstance
	Pair1                uint32
	Pair2                float32
	Palette              MaterialPalette
}

func (cm *CacheManager) materialByFragID(fragID uint32) *MaterialDef {
	for _, material := range cm.MaterialDefs {
		if material.fragID == fragID {
			return material
		}
	}
	return nil
}

func (cm *CacheManager) materialByTag(tag string) *MaterialDef {
	for _, material := range cm.MaterialDefs {
		if material.Tag == tag {
			return material
		}
	}
	return nil
}

// MaterialPalette is a struct representing a material palette
type MaterialPalette struct {
	fragID    uint32
	Tag       string
	Flags     uint32
	Materials []string
}

func (cm *CacheManager) materialInstanceByFragID(fragID uint32) *MaterialPalette {
	for _, materialInstance := range cm.MaterialPalettes {
		if materialInstance.fragID == fragID {
			return materialInstance
		}
	}
	return nil
}

func (cm *CacheManager) materialInstanceByTag(tag string) *MaterialPalette {
	for _, materialInstance := range cm.MaterialPalettes {
		if materialInstance.Tag == tag {
			return materialInstance
		}
	}
	return nil
}

type DmSpriteDef2 struct {
	fragID               uint32
	Tag                  string
	Flags                uint32
	MaterialPaletteTag   string
	DmTrackTag           string
	Fragment3Ref         int32
	Fragment4Ref         int32
	CenterOffset         [3]float32
	Params2              [3]uint32
	MaxDistance          float32
	Min                  [3]float32
	Max                  [3]float32
	Scale                uint16
	Vertices             [][3]int16
	UVs                  [][2]int16
	VertexNormals        [][3]int8
	Colors               [][4]uint8
	Faces                []Face
	FaceMaterialGroups   [][2]uint16
	SkinAssignmentGroups [][2]uint16
	VertexMaterialGroups [][2]int16
	MeshOps              []MeshOp
}

type MeshOp struct {
	Index1    uint16
	Index2    uint16
	Offset    float32
	Param1    uint8
	TypeField uint8
}

func (cm *CacheManager) meshByFragID(fragID uint32) *DmSpriteDef2 {
	for _, mesh := range cm.DmSpriteDef2s {
		if mesh.fragID == fragID {
			return mesh
		}
	}
	return nil
}

func (cm *CacheManager) meshByTag(tag string) *DmSpriteDef2 {
	for _, mesh := range cm.DmSpriteDef2s {
		if mesh.Tag == tag {
			return mesh
		}
	}
	return nil
}

type MeshInstance struct {
	fragID uint32
	Tag    string
	Mesh   string
	Params uint32
}

func (cm *CacheManager) meshInstanceByFragID(fragID uint32) *MeshInstance {
	for _, meshInstance := range cm.MeshInstances {
		if meshInstance.fragID == fragID {
			return meshInstance
		}
	}
	return nil
}

func (cm *CacheManager) meshInstanceByTag(tag string) *MeshInstance {
	for _, meshInstance := range cm.MeshInstances {
		if meshInstance.Tag == tag {
			return meshInstance
		}
	}
	return nil
}

type AlternateMesh struct {
	fragID         uint32
	Tag            string
	Flags          uint32
	Fragment1Maybe int16
	Material       string
	Fragment3      uint32
	CenterPosition model.Vector3
	Params2        uint32
	Something2     uint32
	Something3     uint32
	Verticies      []model.Vector3
	TexCoords      []model.Vector3
	Normals        []model.Vector3
	Colors         []int32
	Polygons       []*AlternateMeshSpritePolygon
	VertexPieces   []*AlternateMeshVertexPiece
	PostVertexFlag uint32
	RenderGroups   []*AlternateMeshRenderGroup
	VertexTex      []model.Vector2
	Size6Pieces    []*AlternateMeshSize6Entry
}

type AlternateMeshSpritePolygon struct {
	Flag int16
	Unk1 int16
	Unk2 int16
	Unk3 int16
	Unk4 int16
	I1   int16
	I2   int16
	I3   int16
}

type AlternateMeshVertexPiece struct {
	Count  int16
	Offset int16
}

type AlternateMeshRenderGroup struct {
	PolygonCount int16
	MaterialId   int16
}

type AlternateMeshSize6Entry struct {
	Unk1 uint32
	Unk2 uint32
	Unk3 uint32
	Unk4 uint32
	Unk5 uint32
}

func (cm *CacheManager) alternateMeshByFragID(fragID uint32) *AlternateMesh {
	for _, alternateMesh := range cm.AlternateMeshes {
		if alternateMesh.fragID == fragID {
			return alternateMesh
		}
	}
	return nil
}

func (cm *CacheManager) alternateMeshByTag(tag string) *AlternateMesh {
	for _, alternateMesh := range cm.AlternateMeshes {
		if alternateMesh.Tag == tag {
			return alternateMesh
		}
	}
	return nil
}

type Animation struct {
	fragID     uint32
	Tag        string
	Flags      uint32
	Transforms []*AnimationTransform
}
type AnimationTransform struct {
	RotateDenominator int16
	RotateX           int16
	RotateY           int16
	RotateZ           int16
	ShiftX            int16
	ShiftY            int16
	ShiftZ            int16
	ShiftDenominator  int16
}

func (cm *CacheManager) animationByFragID(fragID uint32) *Animation {
	for _, animation := range cm.Animations {
		if animation.fragID == fragID {
			return animation
		}
	}
	return nil
}

func (cm *CacheManager) animationByTag(tag string) *Animation {
	for _, animation := range cm.Animations {
		if animation.Tag == tag {
			return animation
		}
	}
	return nil
}

type AnimationInstance struct {
	fragID    uint32
	Tag       string
	Animation string
	Flags     uint32
	Sleep     uint32
}

func (cm *CacheManager) animationInstanceByFragID(fragID uint32) *AnimationInstance {
	for _, animationInstance := range cm.AnimationInstances {
		if animationInstance.fragID == fragID {
			return animationInstance
		}
	}
	return nil
}

func (cm *CacheManager) animationInstanceByTag(tag string) *AnimationInstance {
	for _, animationInstance := range cm.AnimationInstances {
		if animationInstance.Tag == tag {
			return animationInstance
		}
	}
	return nil
}

type ActorDef struct {
	fragID           uint32
	Tag              string
	Flags            uint32
	Callback         string
	ActionCount      uint32
	FragmentRefCount uint32
	BoundsRef        int32
	CurrentAction    uint32
	Location         [6]float32
	Unk1             uint32
	Actions          []ActorAction
	FragmentRefs     []uint32
	Unk2             uint32
}

type ActorAction struct {
	Unk1 uint32
	Lods []float32
}

func (cm *CacheManager) actorByFragID(fragID uint32) *ActorDef {
	for _, actor := range cm.ActorDefs {
		if actor.fragID == fragID {
			return actor
		}
	}
	return nil
}

func (cm *CacheManager) actorByTag(tag string) *ActorDef {
	for _, actor := range cm.ActorDefs {
		if actor.Tag == tag {
			return actor
		}
	}
	return nil
}

type ActorInst struct {
	fragID         uint32
	Tag            string
	ActorDefTag    string
	Flags          uint32
	SphereTag      string
	CurrentAction  uint32
	Location       [6]float32
	Unk1           uint32
	BoundingRadius float32
	Scale          float32
	Sound          string
	Unk2           int32
}

func (cm *CacheManager) actorInstanceByFragID(fragID uint32) *ActorInst {
	for _, actorInstance := range cm.ActorInsts {
		if actorInstance.fragID == fragID {
			return actorInstance
		}
	}
	return nil
}

func (cm *CacheManager) actorInstanceByTag(tag string) *ActorInst {
	for _, actorInstance := range cm.ActorInsts {
		if actorInstance.Tag == tag {
			return actorInstance
		}
	}
	return nil
}

type Skeleton struct {
	fragID             uint32
	Tag                string
	Flags              uint32
	CollisionVolumeRef uint32
	CenterOffset       [3]uint32
	BoundingRadius     float32
	Bones              []*SkeletonEntry
	Skins              []uint32
	SkinLinks          []uint32
}

type SkeletonEntry struct {
	Tag          string
	Flags        uint32
	Track        string
	MeshOrSprite string
	SubBones     []uint32
}

func (cm *CacheManager) skeletonByFragID(fragID uint32) *Skeleton {
	for _, skeleton := range cm.Skeletons {
		if skeleton.fragID == fragID {
			return skeleton
		}
	}
	return nil
}

func (cm *CacheManager) skeletonByTag(tag string) *Skeleton {
	for _, skeleton := range cm.Skeletons {
		if skeleton.Tag == tag {
			return skeleton
		}
	}
	return nil
}

type SkeletonInstance struct {
	fragID   uint32
	Tag      string
	Skeleton string
	Flags    uint32
}

func (cm *CacheManager) skeletonInstanceByFragID(fragID uint32) *SkeletonInstance {
	for _, skeletonInstance := range cm.SkeletonInstances {
		if skeletonInstance.fragID == fragID {
			return skeletonInstance
		}
	}
	return nil
}

func (cm *CacheManager) skeletonInstanceByTag(tag string) *SkeletonInstance {
	for _, skeletonInstance := range cm.SkeletonInstances {
		if skeletonInstance.Tag == tag {
			return skeletonInstance
		}
	}
	return nil
}

type LightDef struct {
	fragID          uint32
	Tag             string
	Flags           uint32
	FrameCurrentRef uint32
	Sleep           uint32
	LightLevels     []float32
	Colors          [][3]float32
}

func (cm *CacheManager) lightByFragID(fragID uint32) *LightDef {
	for _, light := range cm.LightDefs {
		if light.fragID == fragID {
			return light
		}
	}
	return nil
}

func (cm *CacheManager) lightByTag(tag string) *LightDef {
	for _, light := range cm.LightDefs {
		if light.Tag == tag {
			return light
		}
	}
	return nil
}

type LightInstance struct {
	fragID uint32
	Tag    string
	Light  string
	Flags  uint32
}

func (cm *CacheManager) lightInstanceByFragID(fragID uint32) *LightInstance {
	for _, lightInstance := range cm.LightInstances {
		if lightInstance.fragID == fragID {
			return lightInstance
		}
	}
	return nil
}

func (cm *CacheManager) lightInstanceByTag(tag string) *LightInstance {
	for _, lightInstance := range cm.LightInstances {
		if lightInstance.Tag == tag {
			return lightInstance
		}
	}
	return nil
}

type Sprite3DDef struct {
	fragID        uint32
	Tag           string
	Flags         uint32
	SphereListRef uint32
	CenterOffset  [3]float32
	Radius        float32
	Vertices      [][3]float32
	BspNodes      []*CameraBspNode
}

type CameraBspNode struct {
	FrontTree                   uint32
	BackTree                    uint32
	VertexIndexes               []uint32
	RenderMethod                string
	RenderFlags                 uint8
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          [3]float32
	RenderUVInfoUAxis           [3]float32
	RenderUVInfoVAxis           [3]float32
	RenderUVMapEntries          []CameraBspNodeUVMapEntry
}

type CameraBspNodeUVMapEntry struct {
	UvOrigin [3]float32
	UAxis    [3]float32
	VAxis    [3]float32
}

func (cm *CacheManager) cameraByFragID(fragID uint32) *Sprite3DDef {
	for _, camera := range cm.Sprite3DDefs {
		if camera.fragID == fragID {
			return camera
		}
	}
	return nil
}

func (cm *CacheManager) cameraByTag(tag string) *Sprite3DDef {
	for _, camera := range cm.Sprite3DDefs {
		if camera.Tag == tag {
			return camera
		}
	}
	return nil
}

type CameraInstance struct {
	fragID    uint32
	Tag       string
	CameraTag string
	Flags     uint32
}

func (cm *CacheManager) cameraInstanceByFragID(fragID uint32) *CameraInstance {
	for _, cameraInstance := range cm.CameraInstances {
		if cameraInstance.fragID == fragID {
			return cameraInstance
		}
	}
	return nil
}

func (cm *CacheManager) cameraInstanceByTag(tag string) *CameraInstance {
	for _, cameraInstance := range cm.CameraInstances {
		if cameraInstance.Tag == tag {
			return cameraInstance
		}
	}
	return nil
}

type BspTree struct {
	fragID uint32
	Tag    string
	Nodes  []*BspTreeNode
}

type BspTreeNode struct {
	Normal    model.Vector3
	Distance  float32
	RegionTag string
	Front     *BspTreeNode
	Back      *BspTreeNode
}

func (cm *CacheManager) bspTreeByFragID(fragID uint32) *BspTree {
	for _, bspTree := range cm.BspTrees {
		if bspTree.fragID == fragID {
			return bspTree
		}
	}
	return nil
}

func (cm *CacheManager) bspTreeByTag(tag string) *BspTree {
	for _, bspTree := range cm.BspTrees {
		if bspTree.Tag == tag {
			return bspTree
		}
	}
	return nil
}

type Region struct {
	fragID               uint32
	Tag                  string
	Flags                uint32
	AmbientLightRef      int32
	RegionVertexCount    uint32
	RegionProximalCount  uint32
	RenderVertexCount    uint32
	WallCount            uint32
	ObstacleCount        uint32
	CuttingObstacleCount uint32
	VisibleNodeCount     uint32
	RegionVertices       []model.Vector3
	RegionProximals      []model.Vector2
	RenderVertices       []model.Vector3
	Walls                []*RegionWall
}

type RegionWall struct {
	Flags                       uint32
	VertexCount                 uint32
	RenderMethod                uint32
	RenderFlags                 uint32
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          model.Vector3
	RenderUVInfoUAxis           model.Vector3
	RenderUVInfoVAxis           model.Vector3
	RenderUVMapEntryCount       uint32
	RenderUVMapEntries          []model.Vector2
	Normal                      model.Quad4
	Vertices                    []uint32
}

func (cm *CacheManager) regionByFragID(fragID uint32) *Region {
	for _, region := range cm.Regions {
		if region.fragID == fragID {
			return region
		}
	}
	return nil
}

func (cm *CacheManager) regionByTag(tag string) *Region {
	for _, region := range cm.Regions {
		if region.Tag == tag {
			return region
		}
	}
	return nil
}

type RegionInstance struct {
	fragID     uint32
	Tag        string
	Flags      uint32
	RegionTags []string
	UserData   string
}

func (cm *CacheManager) regionInstanceByFragID(fragID uint32) *RegionInstance {
	for _, regionInstance := range cm.RegionInstances {
		if regionInstance.fragID == fragID {
			return regionInstance
		}
	}
	return nil
}

func (cm *CacheManager) regionInstanceByTag(tag string) *RegionInstance {
	for _, regionInstance := range cm.RegionInstances {
		if regionInstance.Tag == tag {
			return regionInstance
		}
	}
	return nil
}

type AmbientLightInstance struct {
	fragID     uint32
	Tag        string
	LightTag   string
	Flags      uint32
	RegionTags []string
}

func (cm *CacheManager) ambientLightInstanceByFragID(fragID uint32) *AmbientLightInstance {
	for _, ambientLightInstance := range cm.AmbientLightInstances {
		if ambientLightInstance.fragID == fragID {
			return ambientLightInstance
		}
	}
	return nil
}

func (cm *CacheManager) ambientLightInstanceByTag(tag string) *AmbientLightInstance {
	for _, ambientLightInstance := range cm.AmbientLightInstances {
		if ambientLightInstance.Tag == tag {
			return ambientLightInstance
		}
	}
	return nil
}

type PointLight struct {
	fragID      uint32
	Tag         string
	LightDefTag string
	Flags       uint32
	Location    [3]float32
	Radius      float32
}

func (cm *CacheManager) pointLightInstanceByFragID(fragID uint32) *PointLight {
	for _, pointLightInstance := range cm.PointLights {
		if pointLightInstance.fragID == fragID {
			return pointLightInstance
		}
	}
	return nil
}

type Sphere struct {
	fragID uint32
	Tag    string
	Radius float32
}

func (cm *CacheManager) sphereByFragID(fragID uint32) *Sphere {
	for _, sphere := range cm.Spheres {
		if sphere.fragID == fragID {
			return sphere
		}
	}
	return nil
}

func (cm *CacheManager) sphereByTag(tag string) *Sphere {
	for _, sphere := range cm.Spheres {
		if sphere.Tag == tag {
			return sphere
		}
	}
	return nil
}

type Face struct {
	Index [3]uint16
	Flags uint16
}

type PolyhedronDef struct {
	fragID   uint32
	Tag      string
	Flags    uint32
	Size1    uint32
	Size2    uint32
	Params1  float32
	Params2  float32
	Entries1 []model.Vector3
	Entries2 []PolyhedronEntries2
}

type PolyhedronEntries2 struct {
	Unk1 uint32
	Unk2 []uint32
}

func (cm *CacheManager) polyhedronByFragID(fragID uint32) *PolyhedronDef {
	for _, polyhedron := range cm.PolyhedronDefs {
		if polyhedron.fragID == fragID {
			return polyhedron
		}
	}
	return nil
}

type PolyhedronInstance struct {
	fragID uint32
	Tag    string
	Flags  uint32
	Scale  float32
}
