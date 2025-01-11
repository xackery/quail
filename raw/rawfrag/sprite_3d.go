package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragSprite3D is Sprite3D in libeq, Camera Reference in openzone, 3DSPRITE (ref) in wld, CameraReference in lantern
type WldFragSprite3D struct {
	parents        []common.TreeLinker
	children       []common.TreeLinker
	fragID         int
	tag            string
	NameRef        int32
	Sprite3DDefRef int32
	Flags          uint32
}

func (e *WldFragSprite3D) FragCode() int {
	return FragCodeSprite3D
}

func (e *WldFragSprite3D) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.Sprite3DDefRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite3D) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Sprite3DDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSprite3D) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragSprite3D) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragSprite3D) Tag() string {
	return e.tag
}

func (e *WldFragSprite3D) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragSprite3D) FragID() int {
	return e.fragID
}

func (e *WldFragSprite3D) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragSprite3D) FragType() string {
	return "S3DI"
}

func (e *WldFragSprite3D) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
