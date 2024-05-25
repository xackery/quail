package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragPolyhedronDef is PolyhedronDef in libeq, Polygon animation in openzone, POLYHEDRONDEFINITION in wld, Fragment17 in lantern
type WldFragPolyhedronDef struct {
	NameRef  int32                       `yaml:"name_ref"`
	Flags    uint32                      `yaml:"flags"`
	Size1    uint32                      `yaml:"size_1"`
	Size2    uint32                      `yaml:"size_2"`
	Params1  float32                     `yaml:"params_1"`
	Params2  float32                     `yaml:"params_2"`
	Entries1 []model.Vector3             `yaml:"entries_1"`
	Entries2 []WldFragPolyhedronEntries2 `yaml:"entries_2"`
}

type WldFragPolyhedronEntries2 struct {
	Unk1 uint32   `yaml:"unk_1"`
	Unk2 []uint32 `yaml:"unk_2"`
}

func (e *WldFragPolyhedronDef) FragCode() int {
	return FragCodePolyhedronDef
}

func (e *WldFragPolyhedronDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.Size1)
	enc.Uint32(e.Size2)
	enc.Float32(e.Params1)
	enc.Float32(e.Params2)
	for _, entry := range e.Entries1 {
		enc.Float32(entry.X)
		enc.Float32(entry.Y)
		enc.Float32(entry.Z)
	}
	for _, entry := range e.Entries2 {
		enc.Uint32(entry.Unk1)
		for _, unk2 := range entry.Unk2 {
			enc.Uint32(unk2)
		}
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPolyhedronDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.Size1 = dec.Uint32()
	e.Size2 = dec.Uint32()
	e.Params1 = dec.Float32()
	e.Params2 = dec.Float32()
	for i := uint32(0); i < e.Size1; i++ {
		v := model.Vector3{}
		v.X = dec.Float32()
		v.Y = dec.Float32()
		v.Z = dec.Float32()
		e.Entries1 = append(e.Entries1, v)
	}
	for i := uint32(0); i < e.Size2; i++ {
		entry := WldFragPolyhedronEntries2{}
		entry.Unk1 = dec.Uint32()
		for j := uint32(0); j < e.Size1; j++ {
			entry.Unk2 = append(entry.Unk2, dec.Uint32())
		}
		e.Entries2 = append(e.Entries2, entry)
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
