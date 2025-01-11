package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragSprite4D is Sprite4D in libeq, empty in openzone, 4DSPRITE (ref) in wld
type WldFragSprite4D struct {
	parents  []common.TreeLinker
	children []common.TreeLinker
	fragID   int
	tag      string
	NameRef  int32  `yaml:"name_ref"`
	FourDRef int32  `yaml:"four_d_ref"`
	Params1  uint32 `yaml:"params_1"`
}

func (e *WldFragSprite4D) FragCode() int {
	return FragCodeSprite4D
}

func (e *WldFragSprite4D) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.FourDRef)
	enc.Uint32(e.Params1)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite4D) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.FourDRef = dec.Int32()
	e.Params1 = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSprite4D) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragSprite4D) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragSprite4D) Tag() string {
	return e.tag
}

func (e *WldFragSprite4D) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragSprite4D) FragID() int {
	return e.fragID
}

func (e *WldFragSprite4D) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragSprite4D) FragType() string {
	return "S4DI"
}

func (e *WldFragSprite4D) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
