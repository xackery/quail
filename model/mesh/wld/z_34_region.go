package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// 0x22 region
type region struct {
	nameRef              int32
	flags                uint32
	ambientLightRef      int32
	regionVertexCount    uint32
	regionProximalCount  uint32
	renderVertexCount    uint32
	wallCount            uint32
	obstacleCount        uint32
	cuttingObstacleCount uint32
	visibleNodeCount     uint32
	regionVertices       []geo.Vector3
	regionProximals      []geo.Vector2
	renderVertices       []geo.Vector3
	walls                []wall
	obstacles            []obstacle
	visNodes             []visNode
	visibles             []visible
	sphere               geo.Quad4
	reverbVolume         float32
	reverbOffset         int32
	userData             string
	meshRef              int32
}

type wall struct {
	flags        uint32
	vertexCount  uint32
	renderMethod uint32
	renderInfo   renderInfo
	normal       geo.Quad4
	vertices     []uint32
}

type obstacle struct {
	flags        uint32
	nextRegion   int32
	obstacleType uint8
	vertexCount  uint32
	vertices     []uint32
	normal       geo.Quad4
	edgeWall     uint32
	userData     string
}

type visNode struct {
	normal    geo.Quad4
	visIndex  uint32
	frontTree uint32
	backTree  uint32
}

type visible struct {
	rangeCount uint16
	ranges     []uint8
}

func (e *WLD) regionRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &region{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	def.ambientLightRef = dec.Int32()
	def.regionVertexCount = dec.Uint32()
	def.regionProximalCount = dec.Uint32()
	def.renderVertexCount = dec.Uint32()
	def.wallCount = dec.Uint32()
	def.obstacleCount = dec.Uint32()
	def.cuttingObstacleCount = dec.Uint32()
	def.visibleNodeCount = dec.Uint32()
	/*def.visibles = make([]uint32, def.visibleNodeCount)
	for i := uint32(0); i < def.visibleNodeCount; i++ {
		def.visibles[i] = dec.Uint32()
	}*/
	def.regionVertices = make([]geo.Vector3, def.regionVertexCount)
	for i := uint32(0); i < def.regionVertexCount; i++ {
		def.regionVertices[i] = geo.Vector3{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
		}
	}
	def.regionProximals = make([]geo.Vector2, def.regionProximalCount)
	for i := uint32(0); i < def.regionProximalCount; i++ {
		def.regionProximals[i] = geo.Vector2{
			X: dec.Float32(),
			Y: dec.Float32(),
		}
	}
	if def.wallCount != 0 {
		def.renderVertexCount = 0
	}

	def.renderVertices = make([]geo.Vector3, def.renderVertexCount)
	for i := uint32(0); i < def.renderVertexCount; i++ {
		def.renderVertices[i] = geo.Vector3{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
		}
	}
	def.walls = make([]wall, def.wallCount)
	for i := uint32(0); i < def.wallCount; i++ {
		wall := wall{}
		wall.flags = dec.Uint32()
		wall.vertexCount = dec.Uint32()
		wall.renderMethod = dec.Uint32()
		wall.renderInfo.uvInfo.origin.X = dec.Float32()
		wall.renderInfo.uvInfo.origin.Y = dec.Float32()
		wall.renderInfo.uvInfo.origin.Z = dec.Float32()
		wall.renderInfo.uvInfo.uAxis.X = dec.Float32()
		wall.renderInfo.uvInfo.uAxis.Y = dec.Float32()
		wall.renderInfo.uvInfo.uAxis.Z = dec.Float32()
		wall.renderInfo.uvInfo.vAxis.X = dec.Float32()
		wall.renderInfo.uvInfo.vAxis.Y = dec.Float32()
		wall.renderInfo.uvInfo.vAxis.Z = dec.Float32()

		wall.normal = geo.Quad4{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
			W: dec.Float32(),
		}
		wall.vertices = make([]uint32, wall.vertexCount)
		for j := uint32(0); j < wall.vertexCount; j++ {
			wall.vertices[j] = dec.Uint32()
		}

		def.walls[i] = wall
	}

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *region) build(e *WLD) error {
	return nil
}

func (e *WLD) regionWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
