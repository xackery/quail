package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// 0x21 worldTree
type worldTree struct {
	nameRef   int32
	nodeCount uint32
	nodes     []worldTreeNode
}

type worldTreeNode struct {
	normal    geo.Vector3
	distance  float32
	regionRef int32
	frontRef  int32
	backRef   int32
}

func (e *WLD) worldTreeRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &worldTree{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.nodeCount = dec.Uint32()
	for i := uint32(0); i < def.nodeCount; i++ {
		node := worldTreeNode{}
		node.normal.X = dec.Float32()
		node.normal.Y = dec.Float32()
		node.normal.Z = dec.Float32()
		node.distance = dec.Float32()
		node.regionRef = dec.Int32()
		node.frontRef = dec.Int32()
		node.backRef = dec.Int32()
		def.nodes = append(def.nodes, node)
	}

	if dec.Error() != nil {
		return fmt.Errorf("worldTreeRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *worldTree) build(e *WLD) error {
	return nil
}

func (e *WLD) worldTreeWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
