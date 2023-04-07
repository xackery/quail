package eqg

import (
	"fmt"

	"github.com/xackery/quail/pfs/archive"
)

// Add adds a new entry to a eqg
func (e *EQG) Add(name string, data []byte) error {
	for _, f := range e.files {
		if f.Name() == name {
			return fmt.Errorf("entry %s already exists", name)
		}
	}
	fe := &archive.FileEntry{}
	err := fe.SetName(name)
	if err != nil {
		return fmt.Errorf("setname: %w", err)
	}
	err = fe.SetData(data)
	if err != nil {
		return fmt.Errorf("setdata: %w", err)
	}
	e.files = append(e.files, fe)
	return nil
}
