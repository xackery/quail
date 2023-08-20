package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
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
	regionVertices       []common.Vector3
	regionProximals      []common.Vector2
	renderVertices       []common.Vector3
	walls                []wall
}

type wall struct {
	flags        uint32
	vertexCount  uint32
	renderMethod uint32
	renderInfo   renderInfo
	normal       common.Quad4
	vertices     []uint32
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
	/*common.visibles = make([]uint32, def.visibleNodeCount)
	for i := uint32(0); i < def.visibleNodeCount; i++ {
		def.visibles[i] = dec.Uint32()
	}*/
	def.regionVertices = make([]common.Vector3, def.regionVertexCount)
	for i := uint32(0); i < def.regionVertexCount; i++ {
		def.regionVertices[i] = common.Vector3{
			X: dec.Float32(),
			Y: dec.Float32(),
			Z: dec.Float32(),
		}
	}
	def.regionProximals = make([]common.Vector2, def.regionProximalCount)
	for i := uint32(0); i < def.regionProximalCount; i++ {
		def.regionProximals[i] = common.Vector2{
			X: dec.Float32(),
			Y: dec.Float32(),
		}
	}
	if def.wallCount != 0 {
		def.renderVertexCount = 0
	}

	def.renderVertices = make([]common.Vector3, def.renderVertexCount)
	for i := uint32(0); i < def.renderVertexCount; i++ {
		def.renderVertices[i] = common.Vector3{
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

		wall.normal = common.Quad4{
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
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *region) build(e *WLD) error {
	return nil
}

func (e *WLD) regionWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
