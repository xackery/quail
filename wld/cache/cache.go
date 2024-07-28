package cache

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
	cm.PolyhedronDefs = []*PolyhedronDef{}
	cm.PolyhedronInstances = []*PolyhedronInstance{}

}
