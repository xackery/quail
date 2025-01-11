package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragDirectionalLight is DirectionalLight in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type WldFragDirectionalLight struct {
	children []common.TreeLinker
	fragID   int
	tag      string
	parents  []common.TreeLinker
}

func (e *WldFragDirectionalLight) FragCode() int {
	return FragCodeDirectionalLight
}

func (e *WldFragDirectionalLight) Write(w io.Writer, isNewWorld bool) error {
	return fmt.Errorf("not implemented")
}

func (e *WldFragDirectionalLight) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragDirectionalLight) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragDirectionalLight) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragDirectionalLight) Tag() string {
	return e.tag
}

func (e *WldFragDirectionalLight) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragDirectionalLight) FragID() int {
	return e.fragID
}

func (e *WldFragDirectionalLight) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragDirectionalLight) FragType() string {
	return "DLID"
}

func (e *WldFragDirectionalLight) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
