package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragPointLightOldDef is empty in libeq, empty in openzone, empty in wld
type WldFragPointLightOldDef struct {
	parents       []common.TreeLinker
	children      []common.TreeLinker
	fragID        int
	tag           string
	NameRef       int32 `yaml:"name_ref"`
	PointLightRef int32 `yaml:"point_light_ref"`
}

func (e *WldFragPointLightOldDef) FragCode() int {
	return FragCodePointLightOldDef
}

func (e *WldFragPointLightOldDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.PointLightRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPointLightOldDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.PointLightRef = dec.Int32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragPointLightOldDef) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragPointLightOldDef) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragPointLightOldDef) Tag() string {
	return e.tag
}

func (e *WldFragPointLightOldDef) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragPointLightOldDef) FragID() int {
	return e.fragID
}

func (e *WldFragPointLightOldDef) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragPointLightOldDef) FragType() string {
	return "PLOD"
}

func (e *WldFragPointLightOldDef) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
