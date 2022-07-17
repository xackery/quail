package eqg

import (
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/common"
)

// File returns data of a file
func (e *EQG) File(name string) ([]byte, error) {
	for _, f := range e.files {
		if f.Name() == name {
			return f.Data(), nil
		}
		if strings.EqualFold(f.Name(), strings.ToLower(name)) {
			return f.Data(), nil
		}
	}
	return nil, fmt.Errorf("read %s: %w", name, os.ErrNotExist)
}

// Files returns a string array of every file inside an EQG
func (e *EQG) Files() []common.Filer {
	return e.files
}
