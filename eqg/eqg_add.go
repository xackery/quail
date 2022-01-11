package eqg

import (
	"fmt"

	"github.com/xackery/quail/common"
)

// Add adds a new entry to a eqg
func (e *EQG) Add(name string, data []byte) error {
	for _, f := range e.files {
		if f.Name == name {
			return fmt.Errorf("entry %s already exists", name)
		}
	}
	e.files = append(e.files, &common.FileEntry{Name: name, Data: data})
	return nil
}
