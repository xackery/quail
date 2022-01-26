package ter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/obj"
)

func (e *TER) ExportObj(objPath string, mtlPath string, matPath string) error {
	var err error
	wm, err := os.Create(mtlPath)
	if err != nil {
		return err
	}
	defer wm.Close()
	objData := &obj.ObjData{
		Name:      e.name,
		Materials: e.materials,
		Triangles: e.triangles,
		Vertices:  e.vertices,
	}
	if objData.Name == "" {

		objData.Name = filepath.Base(objPath)
		objData.Name = strings.TrimSuffix(objData.Name, filepath.Ext(objData.Name))
	}

	err = obj.Export(objData, objPath, mtlPath, matPath)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	return nil
}
