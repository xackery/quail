package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragAmbientLight is AmbientLight in libeq, Ambient Light in openzone, AMBIENTLIGHT in wld, WldFragAmbientLight in lantern
type WldFragAmbientLight struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
	NameRef  int32    `yaml:"name_ref"`
	LightRef int32    `yaml:"light_ref"`
	Flags    uint32   `yaml:"flags"`
	Regions  []uint32 `yaml:"regions"`
}

func (e *WldFragAmbientLight) FragCode() int {
	return FragCodeAmbientLight
}

func (e *WldFragAmbientLight) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.LightRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Regions)))
	for _, region := range e.Regions {
		enc.Uint32(region)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragAmbientLight) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.LightRef = dec.Int32()
	e.Flags = dec.Uint32()
	regionCount := dec.Uint32()
	for i := uint32(0); i < regionCount; i++ {
		e.Regions = append(e.Regions, dec.Uint32())
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragAmbientLight) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragAmbientLight) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragAmbientLight) Tag() string {
	return e.tag
}

func (e *WldFragAmbientLight) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragAmbientLight) FragID() int {
	return e.fragID
}

func (e *WldFragAmbientLight) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragAmbientLight) FragType() string {
	return "ALTI"
}

func (e *WldFragAmbientLight) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
