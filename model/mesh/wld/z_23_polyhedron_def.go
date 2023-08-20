package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
)

// 0x17 polyhedronDef
type polyhedronDef struct {
	nameRef  int32
	flags    uint32
	size1    uint32
	size2    uint32
	params1  float32
	params2  float32
	entries1 []common.Vector3
	entries2 []entries2
}

type entries2 struct {
	unk1 uint32
	unk2 []uint32
}

func (e *WLD) polyhedronDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &polyhedronDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	def.size1 = dec.Uint32()
	def.size2 = dec.Uint32()
	def.params1 = dec.Float32()
	def.params2 = dec.Float32()
	for i := 0; i < int(def.size1); i++ {
		var entry common.Vector3
		entry.X = dec.Float32()
		entry.Y = dec.Float32()
		entry.Z = dec.Float32()
		def.entries1 = append(def.entries1, entry)
	}

	for i := 0; i < int(def.size2); i++ {
		var entry entries2
		entry.unk1 = dec.Uint32()

		for j := 0; j < int(entry.unk1); j++ {
			entry.unk2 = append(entry.unk2, dec.Uint32())
		}
		def.entries2 = append(def.entries2, entry)
	}

	if dec.Error() != nil {
		return fmt.Errorf("polyhedronDefRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *polyhedronDef) build(e *WLD) error {
	return nil
}

func (e *WLD) polyhedronDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
