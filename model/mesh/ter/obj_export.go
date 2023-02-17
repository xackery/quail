package ter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/model/plugin/obj"
)

func (e *TER) ObjExport(objPath string, mtlPath string, matPath string) error {
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

	req := &obj.ObjRequest{
		Data:       objData,
		ObjPath:    objPath,
		MtlPath:    mtlPath,
		MattxtPath: matPath,
	}
	err = obj.Export(req)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	return nil
}
