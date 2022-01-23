package ter

import (
	"fmt"
	"os"

	"github.com/xackery/quail/obj"
)

func (e *TER) ImportObj(objPath string, mtlPath string, matPath string) error {
	var err error
	rm, err := os.Open(mtlPath)
	if err != nil {
		return err
	}
	defer rm.Close()
	objData, err := obj.Import(objPath, mtlPath, matPath)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}
	e.materials = objData.Materials
	e.triangles = objData.Triangles
	e.vertices = objData.Vertices

	return nil
}
