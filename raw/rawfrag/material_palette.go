package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragMaterialPalette is MaterialPalette in libeq, TextureList in openzone, MATERIALPALETTE in wld, WldFragMaterialPalette in lantern
type WldFragMaterialPalette struct {
	parents      []common.TreeLinker
	children     []common.TreeLinker
	fragID       int
	tag          string
	NameRef      int32
	Flags        uint32
	MaterialRefs []uint32
}

func (e *WldFragMaterialPalette) FragCode() int {
	return FragCodeMaterialPalette
}

func (e *WldFragMaterialPalette) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.MaterialRefs)))
	for _, materialRef := range e.MaterialRefs {
		enc.Uint32(materialRef)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragMaterialPalette) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	materialCount := dec.Uint32()
	for i := 0; i < int(materialCount); i++ {
		e.MaterialRefs = append(e.MaterialRefs, dec.Uint32())
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragMaterialPalette) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragMaterialPalette) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragMaterialPalette) Tag() string {
	return e.tag
}

func (e *WldFragMaterialPalette) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragMaterialPalette) FragID() int {
	return e.fragID
}

func (e *WldFragMaterialPalette) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragMaterialPalette) FragType() string {
	return "MATP"
}

func (e *WldFragMaterialPalette) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
