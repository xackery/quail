package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragBlitSpriteDef is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSprite in lantern
type WldFragBlitSpriteDef struct {
	parents           []common.TreeLinker
	children          []common.TreeLinker
	fragID            int
	tag               string
	NameRef           int32  `yaml:"name_ref"`
	Flags             uint32 `yaml:"flags"`
	SpriteInstanceRef uint32 `yaml:"sprite_instance_ref"`
	Unknown           int32  `yaml:"unknown"`
}

func (e *WldFragBlitSpriteDef) FragCode() int {
	return FragCodeBlitSpriteDef
}

func (e *WldFragBlitSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SpriteInstanceRef)
	enc.Int32(e.Unknown)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragBlitSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.SpriteInstanceRef = dec.Uint32()
	e.Unknown = dec.Int32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragBlitSpriteDef) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragBlitSpriteDef) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragBlitSpriteDef) Tag() string {
	return e.tag
}

func (e *WldFragBlitSpriteDef) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragBlitSpriteDef) FragID() int {
	return e.fragID
}

func (e *WldFragBlitSpriteDef) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragBlitSpriteDef) FragType() string {
	return "BLSD"
}

func (e *WldFragBlitSpriteDef) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
