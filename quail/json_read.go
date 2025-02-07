package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	baseFilePath := filepath.Base(path)
	baseName := strings.TrimSuffix(baseFilePath, ".json")

	q.Wld, err = wce.ReadJSON(baseName+".wld", baseFilePath)
	if err != nil {
		return fmt.Errorf("json read: %w", err)
	}
	lightsPath := baseName + "_lights.json"
	lights, err := os.Stat(lightsPath)
	if err == nil {
		q.WldLights, err = wce.ReadJSON("lights.wld", lights.Name())
		if err != nil {
			return fmt.Errorf("json lights read: %w", err)
		}
	}

	objectsPath := baseName + "_objects.json"
	objects, err := os.Stat(objectsPath)
	if err == nil {
		q.WldObject, err = wce.ReadJSON("objects.wld", objects.Name())
		if err != nil {
			return fmt.Errorf("json objects read: %w", err)
		}
	}

	dirs, err := os.ReadDir(baseName)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(dir.Name()))
		if ext == ".bmp" || ext == ".dds" || ext == ".png" || ext == ".jpg" {
			textureData, err := os.ReadFile(baseName + "/" + dir.Name())
			if err != nil {
				fmt.Println("Text", err)
				continue
			}
			q.assetAdd(dir.Name(), textureData)
		}
	}

	return nil
}
