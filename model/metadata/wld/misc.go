package wld

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// Default is empty in libeq, empty in openzone, DEFAULT?? in wld
type Default struct {
}

func (e *Default) FragCode() int {
	return 0x00
}

func (e *Default) Encode(w io.Writer) error {
	return nil
}

func decodeDefault(r io.ReadSeeker) (common.FragmentReader, error) {
	return &Default{}, nil
}

// First is GlobalAmbientLightDef in libeq, First Fragment in openzone, empty in wld, GlobalAmbientLight in lantern
type First struct {
	NameRef int32
}

func (e *First) FragCode() int {
	return 0x35
}

// Encode writes the fragment to the writer
func (e *First) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeFirst(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &First{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// UserData is empty in libeq, empty in openzone, USERDATA in wld
type UserData struct {
}

func (e *UserData) FragCode() int {
	return 0x02
}

func (e *UserData) Encode(w io.Writer) error {
	return nil
}

func decodeUserData(r io.ReadSeeker) (common.FragmentReader, error) {
	return &UserData{}, nil
}

// Sound is empty in libeq, empty in openzone, SOUNDDEFINITION in wld
type Sound struct {
	NameRef int32
	Flags   uint32
}

func (e *Sound) FragCode() int {
	return 0x1F
}

func (e *Sound) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeSound(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Sound{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// SoundRef is empty in libeq, empty in openzone, SOUNDINSTANCE in wld
type SoundRef struct {
	NameRef int32
	Flags   uint32
}

func (e *SoundRef) FragCode() int {
	return 0x20
}

func (e *SoundRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeSoundRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &SoundRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WorldTree is WorldTree in libeq, BSP Tree in openzone, WORLDTREE in wld, BspTree in lantern
type WorldTree struct {
	NameRef   int32
	NodeCount uint32
	Nodes     []WorldTreeNode
}

type WorldTreeNode struct {
	Normal    common.Vector3
	Distance  float32
	RegionRef int32
	FrontRef  int32
	BackRef   int32
}

func (e *WorldTree) FragCode() int {
	return 0x21
}

func (e *WorldTree) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.NodeCount)
	for _, node := range e.Nodes {
		enc.Float32(node.Normal.X)
		enc.Float32(node.Normal.Y)
		enc.Float32(node.Normal.Z)
		enc.Float32(node.Distance)
		enc.Int32(node.RegionRef)
		enc.Int32(node.FrontRef)
		enc.Int32(node.BackRef)
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeWorldTree(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &WorldTree{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.NodeCount = dec.Uint32()
	for i := uint32(0); i < d.NodeCount; i++ {
		node := WorldTreeNode{}
		node.Normal.X = dec.Float32()
		node.Normal.Y = dec.Float32()
		node.Normal.Z = dec.Float32()
		node.Distance = dec.Float32()
		node.RegionRef = dec.Int32()
		node.FrontRef = dec.Int32()
		node.BackRef = dec.Int32()
		d.Nodes = append(d.Nodes, node)
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// Region is Region in libeq, Bsp Region in openzone, REGION in wld, BspRegion in lantern
type Region struct {
	NameRef              int32
	Flags                uint32
	AmbientLightRef      int32
	RegionVertexCount    uint32
	RegionProximalCount  uint32
	RenderVertexCount    uint32
	WallCount            uint32
	ObstacleCount        uint32
	CuttingObstacleCount uint32
	VisibleNodeCount     uint32
	RegionVertices       []common.Vector3
	RegionProximals      []common.Vector2
	RenderVertices       []common.Vector3
	Walls                []Wall
}

type Wall struct {
	Flags                       uint32
	VertexCount                 uint32
	RenderMethod                uint32
	RenderFlags                 uint32
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          common.Vector3
	RenderUVInfoUAxis           common.Vector3
	RenderUVInfoVAxis           common.Vector3
	RenderUVMapEntryCount       uint32
	RenderUVMapEntries          []common.Vector2
	Normal                      common.Quad4
	Vertices                    []uint32
}

func (e *Region) FragCode() int {
	return 0x22
}

func (e *Region) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Int32(e.AmbientLightRef)
	enc.Uint32(e.RegionVertexCount)
	enc.Uint32(e.RegionProximalCount)
	enc.Uint32(e.RenderVertexCount)
	enc.Uint32(e.WallCount)
	enc.Uint32(e.ObstacleCount)
	enc.Uint32(e.CuttingObstacleCount)
	enc.Uint32(e.VisibleNodeCount)
	for _, regionVertex := range e.RegionVertices {
		enc.Float32(regionVertex.X)
		enc.Float32(regionVertex.Y)
		enc.Float32(regionVertex.Z)
	}
	for _, regionProximal := range e.RegionProximals {
		enc.Float32(regionProximal.X)
		enc.Float32(regionProximal.Y)
	}
	for _, renderVertex := range e.RenderVertices {
		enc.Float32(renderVertex.X)
		enc.Float32(renderVertex.Y)
		enc.Float32(renderVertex.Z)
	}
	for _, wall := range e.Walls {
		enc.Uint32(wall.Flags)
		enc.Uint32(wall.VertexCount)
		enc.Uint32(wall.RenderMethod)
		enc.Uint32(wall.RenderFlags)
		enc.Uint32(wall.RenderPen)
		enc.Float32(wall.RenderBrightness)
		enc.Float32(wall.RenderScaledAmbient)
		enc.Uint32(wall.RenderSimpleSpriteReference)
		enc.Float32(wall.RenderUVInfoOrigin.X)
		enc.Float32(wall.RenderUVInfoOrigin.Y)
		enc.Float32(wall.RenderUVInfoOrigin.Z)
		enc.Float32(wall.RenderUVInfoUAxis.X)
		enc.Float32(wall.RenderUVInfoUAxis.Y)
		enc.Float32(wall.RenderUVInfoUAxis.Z)
		enc.Float32(wall.RenderUVInfoVAxis.X)
		enc.Float32(wall.RenderUVInfoVAxis.Y)
		enc.Float32(wall.RenderUVInfoVAxis.Z)
		enc.Uint32(wall.RenderUVMapEntryCount)
		for _, renderUVMapEntry := range wall.RenderUVMapEntries {
			enc.Float32(renderUVMapEntry.X)
			enc.Float32(renderUVMapEntry.Y)
		}
		enc.Float32(wall.Normal.X)
		enc.Float32(wall.Normal.Y)
		enc.Float32(wall.Normal.Z)
		enc.Float32(wall.Normal.W)
		for _, vertex := range wall.Vertices {
			enc.Uint32(vertex)
		}
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeRegion(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Region{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.AmbientLightRef = dec.Int32()
	d.RegionVertexCount = dec.Uint32()
	d.RegionProximalCount = dec.Uint32()
	d.RenderVertexCount = dec.Uint32()
	d.WallCount = dec.Uint32()
	d.ObstacleCount = dec.Uint32()
	d.CuttingObstacleCount = dec.Uint32()
	d.VisibleNodeCount = dec.Uint32()
	d.RegionVertices = make([]common.Vector3, d.RegionVertexCount)
	for i := uint32(0); i < d.RegionVertexCount; i++ {
		d.RegionVertices[i] = common.Vector3{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
		}
	}
	d.RegionProximals = make([]common.Vector2, d.RegionProximalCount)
	for i := uint32(0); i < d.RegionProximalCount; i++ {
		d.RegionProximals[i] = common.Vector2{
			X: dec.Float32(),
			Y: dec.Float32(),
		}
	}
	if d.WallCount != 0 {
		d.RenderVertexCount = 0
	}

	d.RenderVertices = make([]common.Vector3, d.RenderVertexCount)
	for i := uint32(0); i < d.RenderVertexCount; i++ {
		d.RenderVertices[i] = common.Vector3{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
		}
	}

	d.Walls = make([]Wall, d.WallCount)
	for i := uint32(0); i < d.WallCount; i++ {
		wall := Wall{}
		wall.Flags = dec.Uint32()
		wall.VertexCount = dec.Uint32()
		wall.RenderMethod = dec.Uint32()
		wall.RenderFlags = dec.Uint32()
		wall.RenderPen = dec.Uint32()
		wall.RenderBrightness = dec.Float32()
		wall.RenderScaledAmbient = dec.Float32()
		wall.RenderSimpleSpriteReference = dec.Uint32()
		wall.RenderUVInfoOrigin.X = dec.Float32()
		wall.RenderUVInfoOrigin.Y = dec.Float32()
		wall.RenderUVInfoOrigin.Z = dec.Float32()
		wall.RenderUVInfoUAxis.X = dec.Float32()
		wall.RenderUVInfoUAxis.Y = dec.Float32()
		wall.RenderUVInfoUAxis.Z = dec.Float32()
		wall.RenderUVInfoVAxis.X = dec.Float32()
		wall.RenderUVInfoVAxis.Y = dec.Float32()
		wall.RenderUVInfoVAxis.Z = dec.Float32()
		wall.RenderUVMapEntryCount = dec.Uint32()
		for i := uint32(0); i < wall.RenderUVMapEntryCount; i++ {
			wall.RenderUVMapEntries = append(wall.RenderUVMapEntries, common.Vector2{
				X: dec.Float32(),
				Y: dec.Float32(),
			})
		}
		wall.Normal.X = dec.Float32()
		wall.Normal.Y = dec.Float32()
		wall.Normal.Z = dec.Float32()
		wall.Normal.W = dec.Float32()
		wall.Vertices = make([]uint32, wall.VertexCount)
		for i := uint32(0); i < wall.VertexCount; i++ {
			wall.Vertices[i] = dec.Uint32()
		}
		d.Walls[i] = wall
	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// ActiveGeoRegion is empty in libeq, empty in openzone, ACTIVEGEOMETRYREGION in wld
type ActiveGeoRegion struct {
}

func (e *ActiveGeoRegion) FragCode() int {
	return 0x23
}

func (e *ActiveGeoRegion) Encode(w io.Writer) error {
	return nil
}

func decodeActiveGeoRegion(r io.ReadSeeker) (common.FragmentReader, error) {
	return &ActiveGeoRegion{}, nil
}

// SkyRegion is empty in libeq, empty in openzone, SKYREGION in wld
type SkyRegion struct {
}

func (e *SkyRegion) FragCode() int {
	return 0x24
}

func (e *SkyRegion) Encode(w io.Writer) error {
	return nil
}

func decodeSkyRegion(r io.ReadSeeker) (common.FragmentReader, error) {
	return &SkyRegion{}, nil
}

// Zone is Zone in libeq, Region Flag in openzone, ZONE in wld, BspRegionType in lantern
type Zone struct {
}

func (e *Zone) FragCode() int {
	return 0x29
}

func (e *Zone) Encode(w io.Writer) error {
	return nil
}

func decodeZone(r io.ReadSeeker) (common.FragmentReader, error) {
	return &Zone{}, nil
}

// DMTrack is DmTrackDef in libeq, empty in openzone, empty in wld
type DMTrack struct {
}

func (e *DMTrack) FragCode() int {
	return 0x2E
}

func (e *DMTrack) Encode(w io.Writer) error {
	return nil
}

func decodeDMTrack(r io.ReadSeeker) (common.FragmentReader, error) {
	return &DMTrack{}, nil
}

// DMTrackRef is DmTrack in libeq, Mesh Animated Vertices Reference in openzone, empty in wld, MeshAnimatedVerticesReference in lantern
type DMTrackRef struct {
}

func (e *DMTrackRef) FragCode() int {
	return 0x2F
}

func (e *DMTrackRef) Encode(w io.Writer) error {
	return nil
}

func decodeDMTrackRef(r io.ReadSeeker) (common.FragmentReader, error) {
	return &DMTrackRef{}, nil
}

// DMRGBTrack is a list of colors, one per vertex, for baked lighting. It is DmRGBTrackDef in libeq, Vertex Color in openzone, empty in wld, VertexColors in lantern
type DMRGBTrack struct {
}

func (e *DMRGBTrack) FragCode() int {
	return 0x32
}

func (e *DMRGBTrack) Encode(w io.Writer) error {
	return nil
}

func decodeDMRGBTrack(r io.ReadSeeker) (common.FragmentReader, error) {
	return &DMRGBTrack{}, nil
}

// DMRGBTrackRef is DmRGBTrack in libeq, Vertex Color Reference in openzone, empty in wld, VertexColorsReference in lantern
type DMRGBTrackRef struct {
}

func (e *DMRGBTrackRef) FragCode() int {
	return 0x33
}

func (e *DMRGBTrackRef) Encode(w io.Writer) error {
	return nil
}

func decodeDMRGBTrackRef(r io.ReadSeeker) (common.FragmentReader, error) {
	return &DMRGBTrackRef{}, nil
}

// ParticleCloud is ParticleCloudDef in libeq, empty in openzone, empty in wld, ParticleCloud in lantern
type ParticleCloud struct {
	NameRef               int32       `yaml:"name_ref"`
	Unk1                  uint32      `yaml:"unk1"`
	Unk2                  uint32      `yaml:"unk2"`
	ParticleMovement      uint32      `yaml:"particle_movement"` // 0x01 sphere, 0x02 plane, 0x03 stream, 0x04 none
	Flags                 uint32      //Flag 1, High Opacity, Flag 3, Follows Item
	SimultaneousParticles uint32      `yaml:"simultaneous_particles"`
	Unk6                  uint32      `yaml:"unk6"`
	Unk7                  uint32      `yaml:"unk7"`
	Unk8                  uint32      `yaml:"unk8"`
	Unk9                  uint32      `yaml:"unk9"`
	Unk10                 uint32      `yaml:"unk10"`
	SpawnRadius           float32     `yaml:"spawn_radius"` // sphere radius
	SpawnAngle            float32     `yaml:"spawn_angle"`  // cone angle
	SpawnLifespan         uint32      `yaml:"spawn_lifespan"`
	SpawnVelocity         float32     `yaml:"spawn_velocity"`
	SpawnNormalZ          float32     `yaml:"spawn_normal_z"`
	SpawnNormalX          float32     `yaml:"spawn_normal_x"`
	SpawnNormalY          float32     `yaml:"spawn_normal_y"`
	SpawnRate             uint32      `yaml:"spawn_rate"`
	SpawnScale            float32     `yaml:"spawn_scale"`
	Color                 common.RGBA `yaml:"color"`
	SpriteRef             uint32      `yaml:"sprite_ref"`
}

func (e *ParticleCloud) FragCode() int {
	return 0x34
}

func (e *ParticleCloud) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Unk1)
	enc.Uint32(e.Unk2)
	enc.Uint32(e.ParticleMovement)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SimultaneousParticles)
	enc.Uint32(e.Unk6)
	enc.Uint32(e.Unk7)
	enc.Uint32(e.Unk8)
	enc.Uint32(e.Unk9)
	enc.Uint32(e.Unk10)
	enc.Float32(e.SpawnRadius)
	enc.Float32(e.SpawnAngle)
	enc.Uint32(e.SpawnLifespan)
	enc.Float32(e.SpawnVelocity)
	enc.Float32(e.SpawnNormalZ)
	enc.Float32(e.SpawnNormalX)
	enc.Float32(e.SpawnNormalY)
	enc.Uint32(e.SpawnRate)
	enc.Float32(e.SpawnScale)
	enc.Uint8(e.Color.R)
	enc.Uint8(e.Color.G)
	enc.Uint8(e.Color.B)
	enc.Uint8(e.Color.A)

	enc.Uint32(e.SpriteRef)
	if enc.Error() != nil {
		return enc.Error()
	}

	return nil
}

func decodeParticleCloud(r io.ReadSeeker) (common.FragmentReader, error) {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d := &ParticleCloud{}
	d.NameRef = dec.Int32()
	d.Unk1 = dec.Uint32()
	d.Unk2 = dec.Uint32()
	d.ParticleMovement = dec.Uint32()
	d.Flags = dec.Uint32()
	d.SimultaneousParticles = dec.Uint32()
	d.Unk6 = dec.Uint32()
	d.Unk7 = dec.Uint32()
	d.Unk8 = dec.Uint32()
	d.Unk9 = dec.Uint32()
	d.Unk10 = dec.Uint32()
	d.SpawnRadius = dec.Float32()
	d.SpawnAngle = dec.Float32()
	d.SpawnLifespan = dec.Uint32()
	d.SpawnVelocity = dec.Float32()
	d.SpawnNormalZ = dec.Float32()
	d.SpawnNormalX = dec.Float32()
	d.SpawnNormalY = dec.Float32()
	d.SpawnRate = dec.Uint32()
	d.SpawnScale = dec.Float32()
	d.Color = common.RGBA{R: dec.Uint8(), G: dec.Uint8(), B: dec.Uint8(), A: dec.Uint8()}
	d.SpriteRef = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}
