package quail

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wce"
)

// PfsRead imports the quail target file
func (q *Quail) PfsRead(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".wld" {
		r, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open %s: %w", path, err)
		}
		defer r.Close()
		rawWld := &raw.Wld{}
		err = rawWld.Read(r)
		if err != nil {
			return fmt.Errorf("wld read: %w", err)
		}

		q.Wld = wce.New(filepath.Base(path))

		err = q.Wld.ReadWldRaw(rawWld)
		if err != nil {
			return fmt.Errorf("wld read: %w", err)
		}
		return nil
	}
	pfs, err := pfs.NewFile(path)
	if err != nil {
		return fmt.Errorf("pfs load: %w", err)
	}
	defer pfs.Close()

	for _, file := range pfs.Files() {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		reader, err := raw.Read(ext, bytes.NewReader(file.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", file.Name(), err)
		}
		reader.SetFileName(file.Name())
		err = q.RawRead(reader)
		if err != nil {
			return fmt.Errorf("rawRead %s: %w", file.Name(), err)
		}
	}
	return nil
}
