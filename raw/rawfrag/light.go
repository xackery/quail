package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragLight is Light in libeq, Light Source Reference in openzone, POINTLIGHTT ?? in wld, LightSourceReference in lantern
type WldFragLight struct {
	parents     []common.TreeLinker
	children    []common.TreeLinker
	fragID      int
	tag         string
	NameRef     int32  `yaml:"name_ref"`
	LightDefRef int32  `yaml:"light_def_ref"`
	Flags       uint32 `yaml:"flags"`
}

func (e *WldFragLight) FragCode() int {
	return FragCodeLight
}

func (e *WldFragLight) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.LightDefRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragLight) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.LightDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragLight) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragLight) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragLight) Tag() string {
	return e.tag
}

func (e *WldFragLight) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragLight) FragID() int {
	return e.fragID
}

func (e *WldFragLight) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragLight) FragType() string {
	return "LITI"
}

func (e *WldFragLight) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
