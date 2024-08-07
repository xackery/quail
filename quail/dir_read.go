package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/wld"
)

// DirRead loads a .quail directory
func (q *Quail) DirRead(path string) error {

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("path %s is not a directory", path)
	}

	fi, err = os.Stat(path + "/_root.wce")
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("path %s is a directory, but should be a file", path+"/_root.wce")
	}

	baseName := filepath.Base(path)
	baseName = strings.TrimSuffix(baseName, ".quail")
	q.wld = &wld.Wld{
		FileName: baseName + ".wld",
	}
	err = q.wld.ReadAscii(path + "/_root.wce")
	if err != nil {
		return err
	}

	fi, err = os.Stat(path + "/_objects/_root.wce")
	if err == nil && !fi.IsDir() {
		q.wldObject = &wld.Wld{
			FileName: baseName + "objects.wld",
		}
		err = q.wldObject.ReadAscii(path + "/_objects/_root.wce")
		if err != nil {
			return err
		}
	}

	fi, err = os.Stat(path + "/_lights/_root.wce")
	if err == nil && !fi.IsDir() {
		q.wldLights = &wld.Wld{
			FileName: baseName + "lights.wld",
		}
		err = q.wldLights.ReadAscii(path + "/_lights/_root.wce")
		if err != nil {
			return err
		}
	}

	q.Textures = make(map[string][]byte)

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
		textureData, err := os.ReadFile(path + "/" + dir.Name())
		if err != nil {
			return err
		}
		q.Textures[dir.Name()] = textureData
	}

	return nil
}
