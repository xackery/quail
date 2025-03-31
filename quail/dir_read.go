package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/qfs"
	"github.com/xackery/quail/wce"
)

// DirRead loads a .quail directory
func (q *Quail) DirRead(path string) error {

	if q.FileSystem == nil {
		q.FileSystem = &qfs.OSFS{}
	}

	fi, err := q.FileSystem.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("path %s is not a directory", path)
	}

	fi, err = q.FileSystem.Stat(path + "/_root.wce")
	if err != nil {
		// there's a fallback case where only assets are being packed
		dirs, err := q.FileSystem.ReadDir(path + "/assets")
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		for _, dir := range dirs {
			if dir.IsDir() {
				continue
			}
			textureData, err := q.FileSystem.ReadFile(path + "/assets/" + dir.Name())
			if err != nil {
				return err
			}
			q.assetAdd(dir.Name(), textureData)
		}
		return nil

	}
	if fi.IsDir() {
		return fmt.Errorf("path %s is a directory, but should be a file", path+"/_root.wce")
	}

	baseName := filepath.Base(path)
	baseName = strings.TrimSuffix(baseName, ".quail")
	q.Wld = wce.New(baseName + ".wld")
	q.Wld.IsStripped = q.IsStripped
	err = q.Wld.ReadAscii(path + "/_root.wce")
	if err != nil {
		return err
	}

	fi, err = q.FileSystem.Stat(path + "/_objects/_root.wce")
	if err == nil && !fi.IsDir() {
		q.WldObject = wce.New(baseName + "objects.wld")
		err = q.WldObject.ReadAscii(path + "/_objects/_root.wce")
		if err != nil {
			return err
		}
	}

	fi, err = q.FileSystem.Stat(path + "/_lights/_root.wce")
	if err == nil && !fi.IsDir() {
		q.WldLights = wce.New(baseName + "lights.wld")
		err = q.WldLights.ReadAscii(path + "/_lights/_root.wce")
		if err != nil {
			return err
		}
	}

	dirs, err := q.FileSystem.ReadDir(path + "/assets")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		textureData, err := q.FileSystem.ReadFile(path + "/assets/" + dir.Name())
		if err != nil {
			return err
		}
		q.assetAdd(dir.Name(), textureData)
	}

	return nil
}
