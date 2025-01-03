package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DirWrite exports the quail target to a directory
func (q *Quail) JsonWrite(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	if q.Wld != nil {
		err = q.Wld.WriteJSON(path)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("No Wld for", path)
	}
	if q.WldObject != nil {
		path := strings.TrimSuffix(path, ".json") + "_objects.json"
		err = q.WldObject.WriteJSON(path)
		if err != nil {
			return fmt.Errorf("write object: %w", err)
		}
	} else {
		fmt.Println("No WldObject for", path)
	}
	if q.WldLights != nil {
		path := strings.TrimSuffix(path, ".json") + "_lights.json"
		err = q.WldLights.WriteJSON(path)
		if err != nil {
			return fmt.Errorf("write lights: %w", err)
		}
	} else {
		fmt.Println("No Lights for", path)
	}
	os.MkdirAll(strings.TrimSuffix(path, ".json"), 0755)

	for name, data := range q.Textures {
		err = os.WriteFile(strings.TrimSuffix(path, ".json")+"/"+name, data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
