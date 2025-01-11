package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragParticleSprite is ParticleSprite in libeq, empty in openzone, PARTICLESPRITE (ref) in wld
type WldFragParticleSprite struct {
	parents              []common.TreeLinker
	children             []common.TreeLinker
	fragID               int
	tag                  string
	NameRef              int32  `yaml:"name_ref"`
	ParticleSpriteDefRef int32  `yaml:"particle_sprite_def_ref"`
	Flags                uint32 `yaml:"flags"`
}

func (e *WldFragParticleSprite) FragCode() int {
	return FragCodeParticleSprite
}

func (e *WldFragParticleSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.ParticleSpriteDefRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragParticleSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.ParticleSpriteDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragParticleSprite) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragParticleSprite) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragParticleSprite) Tag() string {
	return e.tag
}

func (e *WldFragParticleSprite) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragParticleSprite) FragID() int {
	return e.fragID
}

func (e *WldFragParticleSprite) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragParticleSprite) FragType() string {
	return "PASI"
}

func (e *WldFragParticleSprite) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
