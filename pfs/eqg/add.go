package eqg

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs/archive"
)

// Add adds a new entry to a eqg
func (e *EQG) Add(name string, data []byte) error {
	name = strings.ToLower(name)
	for _, f := range e.files {
		if strings.EqualFold(f.Name(), name) {
			//log.Warnf("entry %s already exists", name)
			return nil
			//return fmt.Errorf("entry %s already exists", name)
		}
	}
	log.Debugf("EQG adding %s (%d bytes)", name, len(data))
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
