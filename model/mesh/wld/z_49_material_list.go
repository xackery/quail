package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ghostiam/binstruct"
	"github.com/xackery/quail/log"
)

// materialList 0x31 49
type materialList struct {
	NameRef       int32
	Flags         uint32
	MaterialCount uint32
	MaterialRefs  []uint32 `bin:"len:MaterialCount"`
}

func (e *WLD) materialListRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &materialList{}

	dec := binstruct.NewDecoder(r, binary.LittleEndian)
	err := dec.Decode(def)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *materialList) build(e *WLD) error {
	return nil
}
