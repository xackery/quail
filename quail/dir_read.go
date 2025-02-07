package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/wce"
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
	q.Wld = wce.New(baseName + ".wld")
	err = q.Wld.ReadAscii(path + "/_root.wce")
	if err != nil {
		return err
	}

	fi, err = os.Stat(path + "/_objects/_root.wce")
	if err == nil && !fi.IsDir() {
		q.WldObject = wce.New(baseName + "objects.wld")
		err = q.WldObject.ReadAscii(path + "/_objects/_root.wce")
		if err != nil {
			return err
		}
	}

	fi, err = os.Stat(path + "/_lights/_root.wce")
	if err == nil && !fi.IsDir() {
		q.WldLights = wce.New(baseName + "lights.wld")
		err = q.WldLights.ReadAscii(path + "/_lights/_root.wce")
		if err != nil {
			return err
		}
	}

	dirs, err := os.ReadDir(path + "/assets")
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		textureData, err := os.ReadFile(path + "/assets/" + dir.Name())
		if err != nil {
			return err
		}
		q.assetAdd(dir.Name(), textureData)
	}

	return nil
}
