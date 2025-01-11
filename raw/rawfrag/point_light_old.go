package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragPointLightOld is empty in libeq, empty in openzone, POINTLIGHT?? in wld
type WldFragPointLightOld struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
	NameRef  int32  `yaml:"name_ref"`
	Flags    uint32 `yaml:"flags"`
}

func (e *WldFragPointLightOld) FragCode() int {
	return FragCodePointLightOld
}

func (e *WldFragPointLightOld) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPointLightOld) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragPointLightOld) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragPointLightOld) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragPointLightOld) Tag() string {
	return e.tag
}

func (e *WldFragPointLightOld) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragPointLightOld) FragID() int {
	return e.fragID
}

func (e *WldFragPointLightOld) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragPointLightOld) FragType() string {
	return "PLOI"
}

func (e *WldFragPointLightOld) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
