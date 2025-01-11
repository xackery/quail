package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragSimpleSpriteDef is SimpleSpriteDef in libeq, WldFragSimpleSpriteDef Bitmap Info in openzone, SIMPLESPRITEDEF in wld, BitmapInfo in lantern
type WldFragSimpleSpriteDef struct {
	parents      []common.TreeLinker
	children     []common.TreeLinker
	fragID       int
	tag          string
	NameRef      int32    `yaml:"name_ref"`
	Flags        uint32   `yaml:"flags"`
	CurrentFrame int32    `yaml:"current_frame"`
	Sleep        uint32   `yaml:"sleep"`
	BitmapRefs   []uint32 `yaml:"bitmap_refs"`
}

func (e *WldFragSimpleSpriteDef) FragCode() int {
	return FragCodeSimpleSpriteDef
}

func (e *WldFragSimpleSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.BitmapRefs)))
	if e.Flags&0x20 != 0 {
		enc.Int32(e.CurrentFrame)
	}
	if e.Flags&0x08 != 0 && e.Flags&0x10 != 0 {
		enc.Uint32(e.Sleep)
	}
	for _, textureRef := range e.BitmapRefs {
		enc.Uint32(textureRef)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSimpleSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	textureRefCount := dec.Uint32()
	if e.Flags&0x20 != 0 {
		e.CurrentFrame = dec.Int32()
	}
	if e.Flags&0x08 != 0 && e.Flags&0x10 != 0 {
		e.Sleep = dec.Uint32()
	}
	for i := 0; i < int(textureRefCount); i++ {
		e.BitmapRefs = append(e.BitmapRefs, dec.Uint32())
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSimpleSpriteDef) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragSimpleSpriteDef) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragSimpleSpriteDef) Tag() string {
	return e.tag
}

func (e *WldFragSimpleSpriteDef) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragSimpleSpriteDef) FragID() int {
	return e.fragID
}

func (e *WldFragSimpleSpriteDef) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragSimpleSpriteDef) FragType() string {
	return "SISD"
}

func (e *WldFragSimpleSpriteDef) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
