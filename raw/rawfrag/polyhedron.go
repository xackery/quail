package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragPolyhedron is Polyhedron in libeq, Polygon Animation Reference in openzone, POLYHEDRON (ref) in wld, Fragment18 in lantern
type WldFragPolyhedron struct {
	parents     []common.TreeLinker
	children    []common.TreeLinker
	fragID      int
	tag         string
	NameRef     int32   `yaml:"name_ref"`
	FragmentRef int32   `yaml:"fragment_ref"`
	Flags       uint32  `yaml:"flags"`
	Scale       float32 `yaml:"scale"`
}

func (e *WldFragPolyhedron) FragCode() int {
	return FragCodePolyhedron
}

func (e *WldFragPolyhedron) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.FragmentRef)
	enc.Uint32(e.Flags)
	enc.Float32(e.Scale)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPolyhedron) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.FragmentRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.Scale = dec.Float32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragPolyhedron) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragPolyhedron) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragPolyhedron) Tag() string {
	return e.tag
}

func (e *WldFragPolyhedron) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragPolyhedron) FragID() int {
	return e.fragID
}

func (e *WldFragPolyhedron) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragPolyhedron) FragType() string {
	return "PLYI"
}

func (e *WldFragPolyhedron) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
