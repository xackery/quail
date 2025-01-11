package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragRegion is Region in libeq, Bsp WldFragRegion in openzone, REGION in wld, BspRegion in lantern
type WldFragRegion struct {
	nameRef         int32        `yaml:"name_ref"`
	Flags           uint32       `yaml:"flags"`
	AmbientLightRef int32        `yaml:"ambient_light_ref"`
	RegionVertices  [][3]float32 `yaml:"region_vertices"`
	RegionProximals [][2]float32 `yaml:"region_proximals"`
	RenderVertices  [][3]float32 `yaml:"render_vertices"`
	Walls           []Wall       `yaml:"walls"`
	Obstacles       []Obstacle   `yaml:"obstacles"`
	VisNodes        []VisNode    `yaml:"visible_nodes"`
	VisLists        []VisList    `yaml:"vis_lists"`
	Sphere          [4]float32   `yaml:"sphere"`
	ReverbVolume    float32
	ReverbOffset    int32
	UserData        string
	MeshReference   int32
}

type Wall struct {
	Flags                       uint32       `yaml:"flags"`
	VertexCount                 uint32       `yaml:"vertex_count"`
	RenderMethod                uint32       `yaml:"render_method"`
	RenderFlags                 uint32       `yaml:"render_flags"`
	RenderPen                   uint32       `yaml:"render_pen"`
	RenderBrightness            float32      `yaml:"render_brightness"`
	RenderScaledAmbient         float32      `yaml:"render_scaled_ambient"`
	RenderSimpleSpriteReference uint32       `yaml:"render_simple_sprite_reference"`
	RenderUVInfoOrigin          [3]float32   `yaml:"render_uv_info_origin"`
	RenderUVInfoUAxis           [3]float32   `yaml:"render_uv_info_u_axis"`
	RenderUVInfoVAxis           [3]float32   `yaml:"render_uv_info_v_axis"`
	RenderUVMapEntryCount       uint32       `yaml:"render_uv_map_entry_count"`
	RenderUVMapEntries          [][2]float32 `yaml:"render_uv_map_entries"`
	Normal                      [4]float32   `yaml:"normal"`
	Vertices                    []uint32     `yaml:"vertices"`
}

type Obstacle struct {
	Flags      uint32
	NextRegion int32
	Type       int32
	Vertices   []uint32
	NormalABCD [4]float32
	EdgeWall   uint32
	UserData   string
}

type VisNode struct {
	NormalABCD   [4]float32
	VisListIndex uint32
	FrontTree    uint32
	BackTree     uint32
}

type VisList struct {
	Ranges []byte
}

func (e *WldFragRegion) FragCode() int {
	return FragCodeRegion
}

func (e *WldFragRegion) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	padStart := enc.Pos()
	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)
	enc.Int32(e.AmbientLightRef)
	enc.Uint32(uint32(len(e.RegionVertices)))
	enc.Uint32(uint32(len(e.RegionProximals)))
	enc.Uint32(uint32(len(e.RenderVertices)))
	enc.Uint32(uint32(len(e.Walls)))
	enc.Uint32(uint32(len(e.Obstacles)))
	enc.Uint32(0) // cuttingobstaclecount
	enc.Uint32(uint32(len(e.VisNodes)))
	enc.Uint32(uint32(len(e.VisLists)))
	for _, regionVertex := range e.RegionVertices {
		enc.Float32(regionVertex[0])
		enc.Float32(regionVertex[1])
		enc.Float32(regionVertex[2])
	}
	for _, regionProximal := range e.RegionProximals {
		enc.Float32(regionProximal[0])
		enc.Float32(regionProximal[1])
	}
	for _, renderVertex := range e.RenderVertices {
		enc.Float32(renderVertex[0])
		enc.Float32(renderVertex[1])
		enc.Float32(renderVertex[2])
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
		enc.Float32(wall.RenderUVInfoOrigin[0])
		enc.Float32(wall.RenderUVInfoOrigin[1])
		enc.Float32(wall.RenderUVInfoOrigin[2])
		enc.Float32(wall.RenderUVInfoUAxis[0])
		enc.Float32(wall.RenderUVInfoUAxis[1])
		enc.Float32(wall.RenderUVInfoUAxis[2])
		enc.Float32(wall.RenderUVInfoVAxis[0])
		enc.Float32(wall.RenderUVInfoVAxis[1])
		enc.Float32(wall.RenderUVInfoVAxis[2])
		enc.Uint32(wall.RenderUVMapEntryCount)
		for _, renderUVMapEntry := range wall.RenderUVMapEntries {
			enc.Float32(renderUVMapEntry[0])
			enc.Float32(renderUVMapEntry[1])
		}
		enc.Float32(wall.Normal[0])
		enc.Float32(wall.Normal[1])
		enc.Float32(wall.Normal[2])
		enc.Float32(wall.Normal[3])
		for _, vertex := range wall.Vertices {
			enc.Uint32(vertex)
		}
	}
	for _, obstacle := range e.Obstacles {
		enc.Uint32(obstacle.Flags)
		enc.Int32(obstacle.NextRegion)
		enc.Int32(obstacle.Type)
		vertexCount := uint32(len(obstacle.Vertices))
		enc.Uint32(vertexCount)
		for _, vertex := range obstacle.Vertices {
			enc.Uint32(vertex)
		}
		if obstacle.Type == -15 { // edgepolygonnormalabcd
			enc.Float32(obstacle.NormalABCD[0])
			enc.Float32(obstacle.NormalABCD[1])
			enc.Float32(obstacle.NormalABCD[2])
			enc.Float32(obstacle.NormalABCD[3])
		}
		if obstacle.Type == 18 { // edgewall
			enc.Uint32(obstacle.EdgeWall)
		}
		if obstacle.Flags&0x04 != 0 { // userdata
			enc.StringLenPrefixUint32(obstacle.UserData)
		}
	}
	for _, visNode := range e.VisNodes {
		enc.Float32(visNode.NormalABCD[0])
		enc.Float32(visNode.NormalABCD[1])
		enc.Float32(visNode.NormalABCD[2])
		enc.Float32(visNode.NormalABCD[3])
		enc.Uint32(visNode.VisListIndex)
		enc.Uint32(visNode.FrontTree)
		enc.Uint32(visNode.BackTree)
	}
	for _, visList := range e.VisLists {
		if e.Flags&0x80 == 0x80 {
			enc.Uint16(uint16(len(visList.Ranges)))
			for _, val := range visList.Ranges {
				enc.Byte(val)
			}
		} else {

			enc.Uint16(uint16(len(visList.Ranges) / 2))
			for i := 0; i < len(visList.Ranges); i += 2 {
				enc.Byte(visList.Ranges[i])
				enc.Byte(visList.Ranges[i+1])
			}
		}
	}
	if e.Flags&0x01 != 0 { // has sphere
		enc.Float32(e.Sphere[0])
		enc.Float32(e.Sphere[1])
		enc.Float32(e.Sphere[2])
		enc.Float32(e.Sphere[3])
	}

	if e.Flags&0x02 != 0 { // has reverb volume
		enc.Float32(e.ReverbVolume)
	}

	if e.Flags&0x04 != 0 { // has reverb offset
		enc.Int32(e.ReverbOffset)
	}

	if e.UserData != "" {
		enc.StringLenPrefixUint32(e.UserData)
	} else {
		enc.Uint32(0)
	}

	if e.Flags&0x100 != 0 || e.Flags&0x40 != 0 { // has mesh reference
		enc.Int32(e.MeshReference)
	}

	diff := enc.Pos() - padStart
	paddingSize := (4 - diff%4) % 4
	if paddingSize > 0 {
		enc.Bytes(make([]byte, paddingSize))
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragRegion) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.AmbientLightRef = dec.Int32()
	regionVertexCount := dec.Uint32()
	regionProximalCount := dec.Uint32()
	renderVertexCount := dec.Uint32()
	wallCount := dec.Uint32()
	obstacleCount := dec.Uint32()
	cuttingObstacleCount := dec.Uint32()
	if cuttingObstacleCount > 0 {
		return fmt.Errorf("you found an unknown cutting obstacle count! Let xack know")
	}
	visibleNodeCount := dec.Uint32()
	visListCount := dec.Uint32()
	for i := uint32(0); i < regionVertexCount; i++ {
		e.RegionVertices[i] = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
	}
	for i := uint32(0); i < regionProximalCount; i++ {
		e.RegionProximals[i] = [2]float32{dec.Float32(), dec.Float32()}
	}

	for i := uint32(0); i < renderVertexCount; i++ {
		e.RenderVertices[i] = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
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
		wall.RenderUVInfoOrigin[0] = dec.Float32()
		wall.RenderUVInfoOrigin[1] = dec.Float32()
		wall.RenderUVInfoOrigin[2] = dec.Float32()
		wall.RenderUVInfoUAxis[0] = dec.Float32()
		wall.RenderUVInfoUAxis[1] = dec.Float32()
		wall.RenderUVInfoUAxis[2] = dec.Float32()
		wall.RenderUVInfoVAxis[0] = dec.Float32()
		wall.RenderUVInfoVAxis[1] = dec.Float32()
		wall.RenderUVInfoVAxis[2] = dec.Float32()
		wall.RenderUVMapEntryCount = dec.Uint32()
		for i := uint32(0); i < wall.RenderUVMapEntryCount; i++ {
			wall.RenderUVMapEntries = append(wall.RenderUVMapEntries, [2]float32{dec.Float32(), dec.Float32()})
		}
		wall.Normal[0] = dec.Float32()
		wall.Normal[1] = dec.Float32()
		wall.Normal[2] = dec.Float32()
		wall.Normal[3] = dec.Float32()
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
			obstacle.NormalABCD[0] = dec.Float32()
			obstacle.NormalABCD[1] = dec.Float32()
			obstacle.NormalABCD[2] = dec.Float32()
			obstacle.NormalABCD[3] = dec.Float32()
		}
		if obstacle.Type == 18 { // edgewall
			obstacle.EdgeWall = dec.Uint32()
		}
		if obstacle.Flags&0x04 != 0 { // userdata
			obstacle.UserData = dec.StringLenPrefixUint32()
		}

		e.Obstacles = append(e.Obstacles, obstacle)
	}

	for i := uint32(0); i < visibleNodeCount; i++ {
		visNode := VisNode{
			NormalABCD:   [4]float32{dec.Float32(), dec.Float32(), dec.Float32(), dec.Float32()},
			VisListIndex: dec.Uint32(),
			FrontTree:    dec.Uint32(),
			BackTree:     dec.Uint32(),
		}
		e.VisNodes = append(e.VisNodes, visNode)
	}

	for i := uint32(0); i < visListCount; i++ {
		visList := VisList{}
		rangeCount := dec.Uint16()
		for i := uint16(0); i < rangeCount; i++ {
			if e.Flags&0x80 != 0 {
				visList.Ranges = append(visList.Ranges, dec.Byte())
			} else {
				visList.Ranges = append(visList.Ranges, dec.Byte())
				visList.Ranges = append(visList.Ranges, dec.Byte())
			}
		}
		e.VisLists = append(e.VisLists, visList)
	}

	if e.Flags&0x01 != 0 { // has sphere
		e.Sphere[0] = dec.Float32()
		e.Sphere[1] = dec.Float32()
		e.Sphere[2] = dec.Float32()
		e.Sphere[3] = dec.Float32()
	}

	if e.Flags&0x02 != 0 { // has reverb volume
		e.ReverbVolume = dec.Float32()
	}

	if e.Flags&0x04 != 0 { // has reverb offset
		e.ReverbOffset = dec.Int32()
	}

	e.UserData = dec.StringLenPrefixUint32()

	if e.Flags&0x100 != 0 || e.Flags&0x40 != 0 { // has mesh reference
		e.MeshReference = dec.Int32()
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragRegion) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragRegion) SetNameRef(nameRef int32) {
	e.nameRef = nameRef
}
