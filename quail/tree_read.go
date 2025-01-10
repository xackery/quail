package quail

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/tree"
)

// TreeRead imports the quail target file
func (q *Quail) TreeRead(path string) error {
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

		err = tree.Dump(rawWld, os.Stdout)
		if err != nil {
			return fmt.Errorf("tree dump: %w", err)
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
		fmt.Fprintf(os.Stdout, "file: %s\n", file.Name())
		err = tree.Dump(reader, os.Stdout)
		if err != nil {
			return fmt.Errorf("rawDump %s: %w", file.Name(), err)
		}

	}
	return nil
}
