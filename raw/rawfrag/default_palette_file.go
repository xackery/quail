package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragDefaultPaletteFile is DefaultPaletteFile in libeq, empty in openzone, DEFAULTPALETTEFILE in wld
type WldFragDefaultPaletteFile struct {
	parents    []common.TreeLinker
	children   []common.TreeLinker
	fragID     int
	tag        string
	NameRef    int32  `yaml:"name_ref"`
	NameLength uint16 `yaml:"name_length"`
	FileName   string `yaml:"file_name"`
}

func (e *WldFragDefaultPaletteFile) FragCode() int {
	return FragCodeDefaultPaletteFile
}

func (e *WldFragDefaultPaletteFile) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint16(e.NameLength)
	enc.String(e.FileName)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDefaultPaletteFile) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.NameLength = dec.Uint16()
	e.FileName = dec.StringFixed(int(e.NameLength))
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragDefaultPaletteFile) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragDefaultPaletteFile) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragDefaultPaletteFile) Tag() string {
	return e.tag
}

func (e *WldFragDefaultPaletteFile) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragDefaultPaletteFile) FragID() int {
	return e.fragID
}

func (e *WldFragDefaultPaletteFile) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragDefaultPaletteFile) FragType() string {
	return "DPLF"
}

func (e *WldFragDefaultPaletteFile) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
