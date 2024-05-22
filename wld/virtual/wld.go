// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package virtual

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/raw"
)

// Wld is a struct representing a Wld file
type Wld struct {
	FileName              string
	GlobalAmbientLight    string
	Version               uint32
	Bitmaps               []*Bitmap
	Sprites               []*Sprite
	SpriteInstances       []*SpriteInstance
	Particles             []*Particle
	ParticleInstances     []*ParticleInstance
	Materials             []*Material
	MaterialInstances     []*MaterialInstance
	Meshes                []*Mesh
	MeshInstances         []*MeshInstance
	AlternateMeshes       []*AlternateMesh
	Actors                []*Actor
	ActorInstances        []*ActorInstance
	Animations            []*Animation
	AnimationInstances    []*AnimationInstance
	Skeletons             []*Skeleton
	SkeletonInstances     []*SkeletonInstance
	Lights                []*Light
	LightInstances        []*LightInstance
	AmbientLightInstances []*AmbientLightInstance
	PointLightInstances   []*PointLightInstance
	Cameras               []*Camera
	CameraInstances       []*CameraInstance
	BspTrees              []*BspTree
	Regions               []*Region
	RegionInstances       []*RegionInstance
	Spheres               []*Sphere
}

// Bitmap is a struct representing a material
type Bitmap struct {
	fragID          uint32
	Tag             string
	Textures        []string
	SimpleSpriteDef Sprite
	SimpleSprite    SpriteInstance
}

func (wld *Wld) bitmapByFragID(fragID uint32) *Bitmap {
	for _, bitmap := range wld.Bitmaps {
		if bitmap.fragID == fragID {
			return bitmap
		}
	}
	return nil
}

func (wld *Wld) bitmapByTag(tag string) *Bitmap {
	for _, bitmap := range wld.Bitmaps {
		if bitmap.Tag == tag {
			return bitmap
		}
	}
	return nil
}

type Sprite struct {
	fragID       uint32
	Tag          string
	Flags        uint32
	CurrentFrame int32
	Sleep        uint32
	Bitmaps      []string
}

func (e *Sprite) Ascii() string {
	// spk usually

	out := "SIMPLESPRITEDEF\n"
	out += fmt.Sprintf("\tSIMPLESPRITETAG \"%s\"\n", e.Tag)
	out += fmt.Sprintf("\tNUMFRAMES %d\n", 1)
	out += fmt.Sprintf("\tBMINFO %s\n", strings.Join(e.Bitmaps, ` `))
	out += "ENDSIMPLESPRITEDEF\n\n"
	return out
}

func (wld *Wld) spriteByFragID(fragID uint32) *Sprite {
	for _, sprite := range wld.Sprites {
		if sprite.fragID == fragID {
			return sprite
		}
	}
	return nil
}

func (wld *Wld) spriteByTag(tag string) *Sprite {
	for _, sprite := range wld.Sprites {
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

func (wld *Wld) spriteInstanceByFragID(fragID uint32) *SpriteInstance {
	for _, spriteInstance := range wld.SpriteInstances {
		if spriteInstance.fragID == fragID {
			return spriteInstance
		}
	}
	return nil
}

func (wld *Wld) spriteInstanceByTag(tag string) *SpriteInstance {
	for _, spriteInstance := range wld.SpriteInstances {
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

func (e *Particle) Ascii() string {
	// SPK usually
	out := "PARTICLEDEFINITION\n"
	out += fmt.Sprintf("\tTAG \"%s\"\n", e.Tag)
	if e.spriteInstance != nil {
		out += "\tSIMPLESPRITEINST\n"
		out += fmt.Sprintf("\t\tTAG %s\n", e.spriteInstance.Tag)
		if e.spriteInstance.Flags != 0 {
			out += fmt.Sprintf("\t\tFLAGS %d\n", e.spriteInstance.Flags)
		}
		out += "\tENDSIMPLESPRITEINST\n"
	}
	out += "ENDPARTICLEDEFINITION\n\n"
	return out

}

func (wld *Wld) particleByFragID(fragID uint32) *Particle {
	for _, particle := range wld.Particles {
		if particle.fragID == fragID {
			return particle
		}
	}
	return nil
}

func (wld *Wld) particleByTag(tag string) *Particle {
	for _, particle := range wld.Particles {
		if particle.Tag == tag {
			return particle
		}
	}
	return nil
}

type ParticleInstance struct {
	fragID                uint32
	Tag                   string
	Unk1                  uint32   `yaml:"unk1"`
	Unk2                  uint32   `yaml:"unk2"`
	ParticleMovement      uint32   `yaml:"particle_movement"` // 0x01 sphere, 0x02 plane, 0x03 stream, 0x04 none
	Flags                 uint32   //Flag 1, High Opacity, Flag 3, Follows Item
	SimultaneousParticles uint32   `yaml:"simultaneous_particles"`
	Unk6                  uint32   `yaml:"unk6"`
	Unk7                  uint32   `yaml:"unk7"`
	Unk8                  uint32   `yaml:"unk8"`
	Unk9                  uint32   `yaml:"unk9"`
	Unk10                 uint32   `yaml:"unk10"`
	SpawnRadius           float32  `yaml:"spawn_radius"` // sphere radius
	SpawnAngle            float32  `yaml:"spawn_angle"`  // cone angle
	SpawnLifespan         uint32   `yaml:"spawn_lifespan"`
	SpawnVelocity         float32  `yaml:"spawn_velocity"`
	SpawnNormalZ          float32  `yaml:"spawn_normal_z"`
	SpawnNormalX          float32  `yaml:"spawn_normal_x"`
	SpawnNormalY          float32  `yaml:"spawn_normal_y"`
	SpawnRate             uint32   `yaml:"spawn_rate"`
	SpawnScale            float32  `yaml:"spawn_scale"`
	Color                 raw.RGBA `yaml:"color"`
	particle              *Particle
}

func (e *ParticleInstance) Ascii() string {
	return ""
}

func (wld *Wld) particleInstanceByFragID(fragID uint32) *ParticleInstance {
	for _, particleInstance := range wld.ParticleInstances {
		if particleInstance.fragID == fragID {
			return particleInstance
		}
	}
	return nil
}

func (wld *Wld) particleInstanceByTag(tag string) *ParticleInstance {
	for _, particleInstance := range wld.ParticleInstances {
		if particleInstance.Tag == tag {
			return particleInstance
		}
	}
	return nil
}

// Material is a struct representing a material
type Material struct {
	fragID         uint32
	Tag            string
	Flags          uint32  `yaml:"flags"`
	RenderMethod   uint32  `yaml:"render_method"`
	RGBPen         uint32  `yaml:"rgb_pen"`
	Brightness     float32 `yaml:"brightness"`
	ScaledAmbient  float32 `yaml:"scaled_ambient"`
	Texture        string  `yaml:"texture"`
	spriteInstance *SpriteInstance
	Pairs          [2]uint32 `yaml:"pairs"`
	Palette        MaterialInstance
}

func (e *Material) Ascii() string {
	// SPK usually
	out := "MATERIALDEFINITION\n"
	out += fmt.Sprintf("\tTAG \"%s\"\n", e.Tag)
	out += fmt.Sprintf("\tRENDERMETHOD %s\n", renderMethod(e.RenderMethod))
	out += fmt.Sprintf("\tRGBPEN %d %d %d\n", (e.RGBPen>>16)&0xFF, (e.RGBPen>>8)&0xFF, e.RGBPen&0xFF)
	out += fmt.Sprintf("\tBRIGHTNESS %0.6f\n", e.Brightness)
	out += fmt.Sprintf("\tSCALEDAMBIENT %0.6f\n", e.ScaledAmbient)

	if e.spriteInstance != nil {
		out += "\tSIMPLESPRITEINST\n"
		out += fmt.Sprintf("\t\tTAG %s\n", e.spriteInstance.Tag)
		if e.spriteInstance.Flags != 0 {
			out += fmt.Sprintf("\t\tFLAGS %d\n", e.spriteInstance.Flags)
		}
		out += "\tENDSIMPLESPRITEINST\n"
	}
	out += fmt.Sprintf("\tTEXTURE \"%s\" // likely simplespriteref?\n", e.Texture)
	out += fmt.Sprintf("\tPAIRS %d %d\n", e.Pairs[0], e.Pairs[1])
	out += "ENDMATERIALDEFINITION\n\n"
	return out
}

func (wld *Wld) materialByFragID(fragID uint32) *Material {
	for _, material := range wld.Materials {
		if material.fragID == fragID {
			return material
		}
	}
	return nil
}

func (wld *Wld) materialByTag(tag string) *Material {
	for _, material := range wld.Materials {
		if material.Tag == tag {
			return material
		}
	}
	return nil
}

// MaterialInstance is a struct representing a material palette
type MaterialInstance struct {
	fragID    uint32
	Tag       string
	Flags     uint32
	Materials []string
}

func (wld *Wld) materialInstanceByFragID(fragID uint32) *MaterialInstance {
	for _, materialInstance := range wld.MaterialInstances {
		if materialInstance.fragID == fragID {
			return materialInstance
		}
	}
	return nil
}

func (wld *Wld) materialInstanceByTag(tag string) *MaterialInstance {
	for _, materialInstance := range wld.MaterialInstances {
		if materialInstance.Tag == tag {
			return materialInstance
		}
	}
	return nil
}

type Mesh struct {
	fragID            uint32
	Tag               string
	Flags             uint32         `yaml:"flags"`
	MaterialInstance  string         `yaml:"material_instance"`
	AnimationInstance string         `yaml:"animation_instance"`
	Fragment3Ref      int32          `yaml:"fragment_3_ref"`
	Fragment4Ref      int32          `yaml:"fragment_4_ref"` // unknown, usually ref to first texture
	Center            raw.Vector3    `yaml:"center"`
	Params2           raw.UIndex3    `yaml:"params_2"`
	MaxDistance       float32        `yaml:"max_distance"`
	Min               raw.Vector3    `yaml:"min"`
	Max               raw.Vector3    `yaml:"max"`
	RawScale          uint16         `yaml:"raw_scale"`
	MeshopCount       uint16         `yaml:"meshop_count"`
	Scale             float32        `yaml:"scale"`
	Vertices          [][3]int16     `yaml:"vertices"`
	UVs               [][2]int16     `yaml:"uvs"`
	Normals           [][3]int8      `yaml:"normals"`
	Colors            []raw.RGBA     `yaml:"colors"`
	Triangles         []raw.Triangle `yaml:"triangles"`
}

func (wld *Wld) meshByFragID(fragID uint32) *Mesh {
	for _, mesh := range wld.Meshes {
		if mesh.fragID == fragID {
			return mesh
		}
	}
	return nil
}

func (wld *Wld) meshByTag(tag string) *Mesh {
	for _, mesh := range wld.Meshes {
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

func (wld *Wld) meshInstanceByFragID(fragID uint32) *MeshInstance {
	for _, meshInstance := range wld.MeshInstances {
		if meshInstance.fragID == fragID {
			return meshInstance
		}
	}
	return nil
}

func (wld *Wld) meshInstanceByTag(tag string) *MeshInstance {
	for _, meshInstance := range wld.MeshInstances {
		if meshInstance.Tag == tag {
			return meshInstance
		}
	}
	return nil
}

type AlternateMesh struct {
	fragID         uint32
	Tag            string
	Flags          uint32                        `yaml:"flags"`
	Fragment1Maybe int16                         `yaml:"fragment_1_maybe"`
	Material       string                        `yaml:"material"`
	Fragment3      uint32                        `yaml:"fragment_3"`
	CenterPosition raw.Vector3                   `yaml:"center_position"`
	Params2        uint32                        `yaml:"params_2"`
	Something2     uint32                        `yaml:"something_2"`
	Something3     uint32                        `yaml:"something_3"`
	Verticies      []raw.Vector3                 `yaml:"verticies"`
	TexCoords      []raw.Vector3                 `yaml:"tex_coords"`
	Normals        []raw.Vector3                 `yaml:"normals"`
	Colors         []int32                       `yaml:"colors"`
	Polygons       []*AlternateMeshSpritePolygon `yaml:"polygons"`
	VertexPieces   []*AlternateMeshVertexPiece   `yaml:"vertex_pieces"`
	PostVertexFlag uint32                        `yaml:"post_vertex_flag"`
	RenderGroups   []*AlternateMeshRenderGroup   `yaml:"render_groups"`
	VertexTex      []raw.Vector2                 `yaml:"vertex_tex"`
	Size6Pieces    []*AlternateMeshSize6Entry    `yaml:"size_6_pieces"`
}

type AlternateMeshSpritePolygon struct {
	Flag int16 `yaml:"flag"`
	Unk1 int16 `yaml:"unk_1"`
	Unk2 int16 `yaml:"unk_2"`
	Unk3 int16 `yaml:"unk_3"`
	Unk4 int16 `yaml:"unk_4"`
	I1   int16 `yaml:"i_1"`
	I2   int16 `yaml:"i_2"`
	I3   int16 `yaml:"i_3"`
}

type AlternateMeshVertexPiece struct {
	Count  int16 `yaml:"count"`
	Offset int16 `yaml:"offset"`
}

type AlternateMeshRenderGroup struct {
	PolygonCount int16 `yaml:"polygon_count"`
	MaterialId   int16 `yaml:"material_id"`
}

type AlternateMeshSize6Entry struct {
	Unk1 uint32 `yaml:"unk_1"`
	Unk2 uint32 `yaml:"unk_2"`
	Unk3 uint32 `yaml:"unk_3"`
	Unk4 uint32 `yaml:"unk_4"`
	Unk5 uint32 `yaml:"unk_5"`
}

func (wld *Wld) alternateMeshByFragID(fragID uint32) *AlternateMesh {
	for _, alternateMesh := range wld.AlternateMeshes {
		if alternateMesh.fragID == fragID {
			return alternateMesh
		}
	}
	return nil
}

func (wld *Wld) alternateMeshByTag(tag string) *AlternateMesh {
	for _, alternateMesh := range wld.AlternateMeshes {
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

func (wld *Wld) animationByFragID(fragID uint32) *Animation {
	for _, animation := range wld.Animations {
		if animation.fragID == fragID {
			return animation
		}
	}
	return nil
}

func (wld *Wld) animationByTag(tag string) *Animation {
	for _, animation := range wld.Animations {
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

func (wld *Wld) animationInstanceByFragID(fragID uint32) *AnimationInstance {
	for _, animationInstance := range wld.AnimationInstances {
		if animationInstance.fragID == fragID {
			return animationInstance
		}
	}
	return nil
}

func (wld *Wld) animationInstanceByTag(tag string) *AnimationInstance {
	for _, animationInstance := range wld.AnimationInstances {
		if animationInstance.Tag == tag {
			return animationInstance
		}
	}
	return nil
}

type Actor struct {
	fragID           uint32
	Tag              string
	Flags            uint32
	CallbackTag      string      `yaml:"callback_tag"`
	ActionCount      uint32      `yaml:"action_count"`
	FragmentRefCount uint32      `yaml:"fragment_ref_count"`
	BoundsRef        int32       `yaml:"bounds_ref"`
	CurrentAction    uint32      `yaml:"current_action"`
	Offset           raw.Vector3 `yaml:"offset"`
	Rotation         raw.Vector3 `yaml:"rotation"`
	Unk1             uint32      `yaml:"unk1"`
	Actions          []ActorAction
	FragmentRefs     []uint32 `yaml:"fragment_refs"`
	Unk2             uint32   `yaml:"unk2"`
}

func (e *Actor) Ascii() string {
	// spk usually

	out := "ACTORDEF\n"
	out += fmt.Sprintf("\tACTORTAG \"%s\"\n", e.Tag)
	out += fmt.Sprintf("\tCALLBACK %s\n", e.CallbackTag)
	//out += fmt.Sprintf("\tLODCOUNT %d\n", e.LodCount)
	//out += fmt.Sprintf("\tUNK1 %d\n", e.Unk1)
	out += "ENDACTORDEF\n\n"
	return out

}

type ActorAction struct {
	LodCount uint32    `yaml:"lod_count"`
	Unk1     uint32    `yaml:"unk1"`
	Lods     []float32 `yaml:"lods"`
}

func (wld *Wld) actorByFragID(fragID uint32) *Actor {
	for _, actor := range wld.Actors {
		if actor.fragID == fragID {
			return actor
		}
	}
	return nil
}

func (wld *Wld) actorByTag(tag string) *Actor {
	for _, actor := range wld.Actors {
		if actor.Tag == tag {
			return actor
		}
	}
	return nil
}

type ActorInstance struct {
	fragID         uint32
	Tag            string
	ActorTag       string      `yaml:"actor"`
	Flags          uint32      `yaml:"flags"`
	Sphere         string      `yaml:"sphere"`
	CurrentAction  uint32      `yaml:"current_action"`
	Offset         raw.Vector3 `yaml:"offset"`
	Rotation       raw.Vector3 `yaml:"rotation"`
	Unk1           uint32      `yaml:"unk1"`
	BoundingRadius float32     `yaml:"bounding_radius"`
	Scale          float32     `yaml:"scale"`
	Sound          string      `yaml:"sound"`
	Unk2           int32       `yaml:"unk2"`
}

func (wld *Wld) actorInstanceByFragID(fragID uint32) *ActorInstance {
	for _, actorInstance := range wld.ActorInstances {
		if actorInstance.fragID == fragID {
			return actorInstance
		}
	}
	return nil
}

func (wld *Wld) actorInstanceByTag(tag string) *ActorInstance {
	for _, actorInstance := range wld.ActorInstances {
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
	BoundingRadius     float32          `yaml:"bounding_radius"`
	Bones              []*SkeletonEntry `yaml:"bones"`
	Skins              []uint32         `yaml:"skins"`
	SkinLinks          []uint32         `yaml:"skin_links"`
}

type SkeletonEntry struct {
	Tag          string
	Flags        uint32 `yaml:"flags"`
	Track        string
	MeshOrSprite string
	SubBones     []uint32 `yaml:"sub_bones"`
}

func (wld *Wld) skeletonByFragID(fragID uint32) *Skeleton {
	for _, skeleton := range wld.Skeletons {
		if skeleton.fragID == fragID {
			return skeleton
		}
	}
	return nil
}

func (wld *Wld) skeletonByTag(tag string) *Skeleton {
	for _, skeleton := range wld.Skeletons {
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

func (wld *Wld) skeletonInstanceByFragID(fragID uint32) *SkeletonInstance {
	for _, skeletonInstance := range wld.SkeletonInstances {
		if skeletonInstance.fragID == fragID {
			return skeletonInstance
		}
	}
	return nil
}

func (wld *Wld) skeletonInstanceByTag(tag string) *SkeletonInstance {
	for _, skeletonInstance := range wld.SkeletonInstances {
		if skeletonInstance.Tag == tag {
			return skeletonInstance
		}
	}
	return nil
}

type Light struct {
	fragID          uint32
	Tag             string
	Flags           uint32
	FrameCurrentRef uint32
	Levels          []float32
	Colors          []raw.Vector3
}

func (wld *Wld) lightByFragID(fragID uint32) *Light {
	for _, light := range wld.Lights {
		if light.fragID == fragID {
			return light
		}
	}
	return nil
}

func (wld *Wld) lightByTag(tag string) *Light {
	for _, light := range wld.Lights {
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

func (wld *Wld) lightInstanceByFragID(fragID uint32) *LightInstance {
	for _, lightInstance := range wld.LightInstances {
		if lightInstance.fragID == fragID {
			return lightInstance
		}
	}
	return nil
}

func (wld *Wld) lightInstanceByTag(tag string) *LightInstance {
	for _, lightInstance := range wld.LightInstances {
		if lightInstance.Tag == tag {
			return lightInstance
		}
	}
	return nil
}

type Camera struct {
	fragID        uint32
	Tag           string
	Flags         uint32
	SphereListRef uint32           `yaml:"sphere_list_ref"`
	CenterOffset  raw.Vector3      `yaml:"center_offset"`
	Radius        float32          `yaml:"radius"`
	Vertices      []raw.Vector3    `yaml:"vertices"`
	BspNodes      []*CameraBspNode `yaml:"bsp_nodes"`
}

type CameraBspNode struct {
	FrontTree                   uint32        `yaml:"front_tree"`
	BackTree                    uint32        `yaml:"back_tree"`
	VertexIndexes               []uint32      `yaml:"vertex_indexes"`
	RenderMethod                uint8         `yaml:"render_method"`
	RenderFlags                 uint8         `yaml:"render_flags"`
	RenderPen                   uint8         `yaml:"render_pen"`
	RenderBrightness            uint8         `yaml:"render_brightness"`
	RenderScaledAmbient         uint8         `yaml:"render_scaled_ambient"`
	RenderSimpleSpriteReference uint8         `yaml:"render_simple_sprite_reference"`
	RenderUVInfoOrigin          raw.Vector3   `yaml:"render_uv_info_origin"`
	RenderUVInfoUAxis           raw.Vector3   `yaml:"render_uv_info_u_axis"`
	RenderUVInfoVAxis           raw.Vector3   `yaml:"render_uv_info_v_axis"`
	RenderUVMapEntries          []raw.Vector2 `yaml:"render_uv_map_entries"`
}

func (wld *Wld) cameraByFragID(fragID uint32) *Camera {
	for _, camera := range wld.Cameras {
		if camera.fragID == fragID {
			return camera
		}
	}
	return nil
}

func (wld *Wld) cameraByTag(tag string) *Camera {
	for _, camera := range wld.Cameras {
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

func (wld *Wld) cameraInstanceByFragID(fragID uint32) *CameraInstance {
	for _, cameraInstance := range wld.CameraInstances {
		if cameraInstance.fragID == fragID {
			return cameraInstance
		}
	}
	return nil
}

func (wld *Wld) cameraInstanceByTag(tag string) *CameraInstance {
	for _, cameraInstance := range wld.CameraInstances {
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
	Normal    raw.Vector3
	Distance  float32
	RegionTag string
	Front     *BspTreeNode
	Back      *BspTreeNode
}

func (wld *Wld) bspTreeByFragID(fragID uint32) *BspTree {
	for _, bspTree := range wld.BspTrees {
		if bspTree.fragID == fragID {
			return bspTree
		}
	}
	return nil
}

func (wld *Wld) bspTreeByTag(tag string) *BspTree {
	for _, bspTree := range wld.BspTrees {
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
	AmbientLightRef      int32         `yaml:"ambient_light_ref"`
	RegionVertexCount    uint32        `yaml:"region_vertex_count"`
	RegionProximalCount  uint32        `yaml:"region_proximal_count"`
	RenderVertexCount    uint32        `yaml:"render_vertex_count"`
	WallCount            uint32        `yaml:"wall_count"`
	ObstacleCount        uint32        `yaml:"obstacle_count"`
	CuttingObstacleCount uint32        `yaml:"cutting_obstacle_count"`
	VisibleNodeCount     uint32        `yaml:"visible_node_count"`
	RegionVertices       []raw.Vector3 `yaml:"region_vertices"`
	RegionProximals      []raw.Vector2 `yaml:"region_proximals"`
	RenderVertices       []raw.Vector3 `yaml:"render_vertices"`
	Walls                []*RegionWall
}

type RegionWall struct {
	Flags                       uint32        `yaml:"flags"`
	VertexCount                 uint32        `yaml:"vertex_count"`
	RenderMethod                uint32        `yaml:"render_method"`
	RenderFlags                 uint32        `yaml:"render_flags"`
	RenderPen                   uint32        `yaml:"render_pen"`
	RenderBrightness            float32       `yaml:"render_brightness"`
	RenderScaledAmbient         float32       `yaml:"render_scaled_ambient"`
	RenderSimpleSpriteReference uint32        `yaml:"render_simple_sprite_reference"`
	RenderUVInfoOrigin          raw.Vector3   `yaml:"render_uv_info_origin"`
	RenderUVInfoUAxis           raw.Vector3   `yaml:"render_uv_info_u_axis"`
	RenderUVInfoVAxis           raw.Vector3   `yaml:"render_uv_info_v_axis"`
	RenderUVMapEntryCount       uint32        `yaml:"render_uv_map_entry_count"`
	RenderUVMapEntries          []raw.Vector2 `yaml:"render_uv_map_entries"`
	Normal                      raw.Quad4     `yaml:"normal"`
	Vertices                    []uint32      `yaml:"vertices"`
}

func (wld *Wld) regionByFragID(fragID uint32) *Region {
	for _, region := range wld.Regions {
		if region.fragID == fragID {
			return region
		}
	}
	return nil
}

func (wld *Wld) regionByTag(tag string) *Region {
	for _, region := range wld.Regions {
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

func (wld *Wld) regionInstanceByFragID(fragID uint32) *RegionInstance {
	for _, regionInstance := range wld.RegionInstances {
		if regionInstance.fragID == fragID {
			return regionInstance
		}
	}
	return nil
}

func (wld *Wld) regionInstanceByTag(tag string) *RegionInstance {
	for _, regionInstance := range wld.RegionInstances {
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

func (wld *Wld) ambientLightInstanceByFragID(fragID uint32) *AmbientLightInstance {
	for _, ambientLightInstance := range wld.AmbientLightInstances {
		if ambientLightInstance.fragID == fragID {
			return ambientLightInstance
		}
	}
	return nil
}

func (wld *Wld) ambientLightInstanceByTag(tag string) *AmbientLightInstance {
	for _, ambientLightInstance := range wld.AmbientLightInstances {
		if ambientLightInstance.Tag == tag {
			return ambientLightInstance
		}
	}
	return nil
}

type PointLightInstance struct {
	fragID           uint32
	Tag              string
	LightInstanceTag string
	Flags            uint32
	X                float32
	Y                float32
	Z                float32
	Radius           float32
}

func (wld *Wld) pointLightInstanceByFragID(fragID uint32) *PointLightInstance {
	for _, pointLightInstance := range wld.PointLightInstances {
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

func (wld *Wld) sphereByFragID(fragID uint32) *Sphere {
	for _, sphere := range wld.Spheres {
		if sphere.fragID == fragID {
			return sphere
		}
	}
	return nil
}

func (wld *Wld) sphereByTag(tag string) *Sphere {
	for _, sphere := range wld.Spheres {
		if sphere.Tag == tag {
			return sphere
		}
	}
	return nil
}
