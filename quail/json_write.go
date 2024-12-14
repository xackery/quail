package quail

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/os"
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
	}
	if q.WldObject != nil {
		err = q.WldObject.WriteJSON(strings.TrimSuffix(path, ".json") + "_objects.json")
		if err != nil {
			return fmt.Errorf("write object: %w", err)
		}
	}
	if q.WldLights != nil {
		err = q.WldLights.WriteJSON(strings.TrimSuffix(path, ".json") + "_lights.json")
		if err != nil {
			return fmt.Errorf("write lights: %w", err)
		}
	}

	// for name, data := range q.Textures {

	// 	/* data, err := fixWonkyDDS(name, texture)
	// 	if err != nil {
	// 		return err
	// 	} */
	// 	err = os.WriteFile(path+"/"+name, data, 0644)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
