package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragPointLight is PointLight in libeq, Light Info in openzone, POINTLIGHT in wld, LightInstance in lantern
type WldFragPointLight struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
	NameRef  int32      `yaml:"name_ref"`
	LightRef int32      `yaml:"light_ref"`
	Flags    uint32     `yaml:"flags"`
	Location [3]float32 `yaml:"location"`
	Radius   float32    `yaml:"radius"`
}

func (e *WldFragPointLight) FragCode() int {
	return FragCodePointLight
}

func (e *WldFragPointLight) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.LightRef)
	enc.Uint32(e.Flags)

	enc.Float32(e.Location[0])
	enc.Float32(e.Location[1])
	enc.Float32(e.Location[2])
	enc.Float32(e.Radius)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPointLight) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.LightRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.Location[0] = dec.Float32()
	e.Location[1] = dec.Float32()
	e.Location[2] = dec.Float32()
	e.Radius = dec.Float32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragPointLight) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragPointLight) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragPointLight) Tag() string {
	return e.tag
}

func (e *WldFragPointLight) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragPointLight) FragID() int {
	return e.fragID
}

func (e *WldFragPointLight) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragPointLight) FragType() string {
	return "PLTI"
}

func (e *WldFragPointLight) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
