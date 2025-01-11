package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragSimpleSprite is SimpleSprite in libeq, Texture Bitmap Info Reference in openzone, SIMPLESPRITEINST in wld, BitmapInfoReference in lantern
type WldFragSimpleSprite struct {
	parents   []common.TreeLinker
	children  []common.TreeLinker
	fragID    int
	tag       string
	NameRef   int32  `yaml:"name_ref"`
	SpriteRef uint32 `yaml:"sprite_ref"`
	Flags     uint32 `yaml:"flags"`
}

func (e *WldFragSimpleSprite) FragCode() int {
	return FragCodeSimpleSprite
}

func (e *WldFragSimpleSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.SpriteRef)
	enc.Uint32(e.Flags)
	enc.Bytes(make([]byte, 2)) // TODO: why 2 extra bytes?
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSimpleSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.SpriteRef = dec.Uint32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSimpleSprite) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragSimpleSprite) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragSimpleSprite) Tag() string {
	return e.tag
}

func (e *WldFragSimpleSprite) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragSimpleSprite) FragID() int {
	return e.fragID
}

func (e *WldFragSimpleSprite) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragSimpleSprite) FragType() string {
	return "SISI"
}

func (e *WldFragSimpleSprite) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
