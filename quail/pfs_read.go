package quail

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

// PfsRead imports the quail target file
func (e *Quail) PfsRead(path string) error {
	pfs, err := pfs.NewFile(path)
	if err != nil {
		return fmt.Errorf("eqg load: %w", err)
	}
	defer pfs.Close()

	for _, file := range pfs.Files() {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		reader, err := raw.Read(ext, bytes.NewReader(file.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", file.Name(), err)
		}
		reader.SetFileName(file.Name())
		err = e.RawRead(reader)
		if err != nil {
			return fmt.Errorf("rawRead %s: %w", file.Name(), err)
		}
	}
	return nil
}
