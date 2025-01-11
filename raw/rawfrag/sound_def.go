package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragSoundDef is empty in libeq, empty in openzone, SOUNDINSTANCE in wld
type WldFragSoundDef struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
	NameRef  int32  `yaml:"name_ref"`
	Flags    uint32 `yaml:"flags"`
}

func (e *WldFragSoundDef) FragCode() int {
	return FragCodeSoundDef
}

func (e *WldFragSoundDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSoundDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSoundDef) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragSoundDef) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragSoundDef) Tag() string {
	return e.tag
}

func (e *WldFragSoundDef) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragSoundDef) FragID() int {
	return e.fragID
}

func (e *WldFragSoundDef) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragSoundDef) FragType() string {
	return "SNDD"
}

func (e *WldFragSoundDef) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
