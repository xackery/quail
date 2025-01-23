package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragWorldTree is WorldTree in libeq, BSP Tree in openzone, WORLDTREE in wld, BspTree in lantern
// For serialization, refer to here: https://github.com/knervous/LanternExtractor2/blob/knervous/merged/LanternExtractor/EQ/Wld/DataTypes/BspNode.cs
// For constructing, refer to here: https://github.com/knervous/LanternExtractor2/blob/920541d15958e90aa91f7446a74226cbf26b829a/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs#L304
type WldFragWorldTree struct {
	nameRef int32           `yaml:"name_ref"`
	Nodes   []WorldTreeNode `yaml:"nodes"`
}

type WorldTreeNode struct {
	Normal    [4]float32 `yaml:"normal"`
	RegionRef int32      `yaml:"region_ref"`
	FrontRef  int32      `yaml:"front_ref"`
	BackRef   int32      `yaml:"back_ref"`
}

func (e *WldFragWorldTree) FragCode() int {
	return FragCodeWorldTree
}

func (e *WldFragWorldTree) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint32(uint32(len(e.Nodes)))
	for _, node := range e.Nodes {
		enc.Float32(node.Normal[0])
		enc.Float32(node.Normal[1])
		enc.Float32(node.Normal[2])
		enc.Float32(node.Normal[3])
		enc.Int32(node.RegionRef)
		enc.Int32(node.FrontRef)
		enc.Int32(node.BackRef)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragWorldTree) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	nodeCount := dec.Uint32()
	for i := uint32(0); i < nodeCount; i++ {
		node := WorldTreeNode{}
		node.Normal[0] = dec.Float32()
		node.Normal[1] = dec.Float32()
		node.Normal[2] = dec.Float32()
		node.Normal[3] = dec.Float32()
		node.RegionRef = dec.Int32()
		node.FrontRef = dec.Int32()
		node.BackRef = dec.Int32()
		e.Nodes = append(e.Nodes, node)
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragWorldTree) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragWorldTree) SetNameRef(nameRef int32) {
	e.nameRef = nameRef
}
