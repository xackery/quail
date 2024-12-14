package quail

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/os"
	"github.com/xackery/quail/wce"
)

// JsonRead loads a .json path
func (q *Quail) JsonRead(path string) error {

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("path %s is a directory, but should be a file", path+"/_root.wce")
	}

	baseName := filepath.Base(path)
	baseName = strings.TrimSuffix(baseName, ".json")
	q.Wld, err = wce.ReadJSON(baseName+".wld", path)
	if err != nil {
		return fmt.Errorf("json read: %w", err)
	}

	q.Textures = make(map[string][]byte)

	path = filepath.Dir(path)
	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		ext := filepath.Ext(dir.Name())
		if ext == ".wce" {
			continue
		}
		if ext == ".mod" {
			continue
		}
		if ext == ".json" {
			continue
		}
		textureData, err := os.ReadFile(path + "/" + dir.Name())
		if err != nil {
			return err
		}
		q.Textures[dir.Name()] = textureData
	}

	return nil
}
