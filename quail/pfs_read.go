package quail

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wce"
)

// PfsRead imports the quail target file
func (q *Quail) PfsRead(path string) error {
	ext := strings.ToLower(filepath.Ext(path))

	if ext == ".eqg" {
		archive, err := pfs.NewFile(path)
		if err != nil {
			return fmt.Errorf("open %s: %w", path, err)
		}
		defer archive.Close()

		baseName := filepath.Base(path)
		baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

		q.Wld = wce.New(baseName)
		err = q.Wld.ReadEqgRaw(archive)
		if err != nil {
			return fmt.Errorf("wld read: %w", err)
		}
	}
	pfs, err := pfs.NewFile(path)
	if err != nil {
		return fmt.Errorf("pfs load: %w", err)
	}
	defer pfs.Close()

	for _, file := range pfs.Files() {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".lit" {
			q.assetAdd(file.Name(), file.Data())
			continue
		}
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

func (q *Quail) assetAdd(name string, data []byte) error {
	if q.Assets == nil {
		q.Assets = make(map[string][]byte)
	}
	q.Assets[name] = data
	return nil
}
