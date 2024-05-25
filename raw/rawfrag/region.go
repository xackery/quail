package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragRegion is Region in libeq, Bsp WldFragRegion in openzone, REGION in wld, BspRegion in lantern
type WldFragRegion struct {
	NameRef              int32           `yaml:"name_ref"`
	Flags                uint32          `yaml:"flags"`
	AmbientLightRef      int32           `yaml:"ambient_light_ref"`
	CuttingObstacleCount uint32          `yaml:"cutting_obstacle_count"`
	RegionVertices       []model.Vector3 `yaml:"region_vertices"`
	RegionProximals      []model.Vector2 `yaml:"region_proximals"`
	RenderVertices       []model.Vector3 `yaml:"render_vertices"`
	Walls                []Wall          `yaml:"walls"`
	Obstacles            []Obstacle      `yaml:"obstacles"`
	VisNodes             []VisNode       `yaml:"visible_nodes"`
	VisLists             []VisList       `yaml:"vis_lists"`
}

type Wall struct {
	Flags                       uint32          `yaml:"flags"`
	VertexCount                 uint32          `yaml:"vertex_count"`
	RenderMethod                uint32          `yaml:"render_method"`
	RenderFlags                 uint32          `yaml:"render_flags"`
	RenderPen                   uint32          `yaml:"render_pen"`
	RenderBrightness            float32         `yaml:"render_brightness"`
	RenderScaledAmbient         float32         `yaml:"render_scaled_ambient"`
	RenderSimpleSpriteReference uint32          `yaml:"render_simple_sprite_reference"`
	RenderUVInfoOrigin          model.Vector3   `yaml:"render_uv_info_origin"`
	RenderUVInfoUAxis           model.Vector3   `yaml:"render_uv_info_u_axis"`
	RenderUVInfoVAxis           model.Vector3   `yaml:"render_uv_info_v_axis"`
	RenderUVMapEntryCount       uint32          `yaml:"render_uv_map_entry_count"`
	RenderUVMapEntries          []model.Vector2 `yaml:"render_uv_map_entries"`
	Normal                      model.Quad4     `yaml:"normal"`
	Vertices                    []uint32        `yaml:"vertices"`
}

type Obstacle struct {
	Flags      uint32
	NextRegion int32
	Type       int32
	Vertices   []uint32
	NormalABCD model.Quad4 // NORMALABCD %f %f %f %f
	EdgeWall   uint32      // EDGEWALL 0 %d
	UserData   string      // USERDATA %s
}

type VisNode struct {
	NormalABCD   model.Quad4 // NORMALABCD %f %f %f %f
	VisListIndex uint32      // VISLISTINDEX %d
	FrontTree    uint32      // FRONTTREE %d
	BackTree     uint32      // BACKTREE %d
}

type VisList struct {
	Ranges []byte
}

func (e *WldFragRegion) FragCode() int {
	return FragCodeRegion
}

func (e *WldFragRegion) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Int32(e.AmbientLightRef)
	enc.Uint32(uint32(len(e.RegionVertices)))
	enc.Uint32(uint32(len(e.RegionProximals)))
	enc.Uint32(uint32(len(e.RenderVertices)))
	enc.Uint32(uint32(len(e.Walls)))
	enc.Uint32(uint32(len(e.Obstacles)))
	enc.Uint32(e.CuttingObstacleCount)
	enc.Uint32(uint32(len(e.VisNodes)))
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
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragRegion) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.AmbientLightRef = dec.Int32()
	regionVertexCount := dec.Uint32()
	regionProximalCount := dec.Uint32()
	renderVertexCount := dec.Uint32()
	wallCount := dec.Uint32()
	obstacleCount := dec.Uint32()
	e.CuttingObstacleCount = dec.Uint32()
	visibleNodeCount := dec.Uint32()
	visListCount := dec.Uint32()
	e.RegionVertices = make([]model.Vector3, regionVertexCount)
	for i := uint32(0); i < regionVertexCount; i++ {
		e.RegionVertices[i] = model.Vector3{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
		}
	}
	e.RegionProximals = make([]model.Vector2, regionProximalCount)
	for i := uint32(0); i < regionProximalCount; i++ {
		e.RegionProximals[i] = model.Vector2{
			X: dec.Float32(),
			Y: dec.Float32(),
		}
	}

	e.RenderVertices = make([]model.Vector3, renderVertexCount)
	for i := uint32(0); i < renderVertexCount; i++ {
		e.RenderVertices[i] = model.Vector3{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
		}
	}

	e.Walls = make([]Wall, wallCount)
	for i := uint32(0); i < wallCount; i++ {
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
			wall.RenderUVMapEntries = append(wall.RenderUVMapEntries, model.Vector2{
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
		e.Walls[i] = wall
	}

	for i := uint32(0); i < obstacleCount; i++ {
		obstacle := Obstacle{
			Flags:      dec.Uint32(),
			NextRegion: dec.Int32(),
			Type:       dec.Int32(),
		}
		//if obstacle.Type == 14 || obstacle.Type == -15 {
		vertexCount := dec.Uint32()

		obstacle.Vertices = make([]uint32, vertexCount)
		for i := uint32(0); i < vertexCount; i++ {
			obstacle.Vertices[i] = dec.Uint32()
		}
		if obstacle.Type == -15 { // edgepolygonnormalabcd
			obstacle.NormalABCD.X = dec.Float32()
			obstacle.NormalABCD.Y = dec.Float32()
			obstacle.NormalABCD.Z = dec.Float32()
			obstacle.NormalABCD.W = dec.Float32()
		}
		if obstacle.Type == 18 { // edgewall
			obstacle.EdgeWall = dec.Uint32()
		}
		if obstacle.Flags&0x04 != 0 { // userdata
			obstacle.UserData = dec.StringLenPrefixUint32()
		}
	}

	for i := uint32(0); i < visibleNodeCount; i++ {
		visNode := VisNode{
			NormalABCD: model.Quad4{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
				W: dec.Float32(),
			},
			VisListIndex: dec.Uint32(),
			FrontTree:    dec.Uint32(),
			BackTree:     dec.Uint32(),
		}
		e.VisNodes = append(e.VisNodes, visNode)
	}

	for i := uint32(0); i < visListCount; i++ {
		visList := VisList{}
		visList.Ranges = dec.Bytes(8)
		e.VisLists = append(e.VisLists, visList)
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
