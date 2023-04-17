package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ghostiam/binstruct"
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
	visibles             []uint32
	regionVertices       []geo.Vector3
	regionProximals      []geo.Vector2
	renderVertices       []geo.Vector3
	walls                []wall
}

type wall struct {
}

func (e *WLD) regionRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &region{}

	dec := binstruct.NewDecoder(r, binary.LittleEndian)
	err := dec.Decode(def)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *region) build(e *WLD) error {
	return nil
}
