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

const (
	ErrorInvalidExt = "invalid extension"
)

// TreeRead imports the quail target file
func (q *Quail) TreeRead(path string, file string) error {
	isValidExt := false
	exts := []string{".eqg", ".s3d", ".pfs", ".pak"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, ext := range exts {
		if strings.HasSuffix(path, ext) {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return q.treeReadFile(nil, path, file)
	}

	pfs, err := pfs.NewFile(path)
	if err != nil {
		return fmt.Errorf("%s load: %w", ext, err)
	}

	return q.treeReadFile(pfs, path, file)
}

func (q *Quail) treeReadFile(pfs *pfs.Pfs, path string, file string) error {
	if pfs == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return q.treeInspectContent(filepath.Base(path), bytes.NewReader(data))
	}

	isFound := false
	for _, fe := range pfs.Files() {
		if len(file) > 1 && !strings.EqualFold(fe.Name(), file) {
			continue
		}

		err := q.treeInspectContent(fe.Name(), bytes.NewReader(fe.Data()))
		if err != nil && err.Error() != ErrorInvalidExt {
			return fmt.Errorf("inspect content: %w", err)
		}
		isFound = true
	}
	if isFound {
		return nil
	}
	if len(file) < 2 {
		return fmt.Errorf("no files found to tree")
	}

	return fmt.Errorf("%s not found in %s", file, filepath.Base(path))
}

func (q *Quail) treeInspectContent(file string, r *bytes.Reader) error {
	var err error
	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".wld":
		fmt.Printf("Tree: %s\n", file)
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
	return fmt.Errorf("%s", ErrorInvalidExt)
}
