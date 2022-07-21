package mds

import (
	"fmt"
	"os"

	"github.com/xackery/quail/obj"
)

func (e *MDS) ObjImport(objPath string, mtlPath string, matPath string) error {
	var err error
	rm, err := os.Open(mtlPath)
	if err != nil {
		return err
	}
	defer rm.Close()
	req := &obj.ObjRequest{
		ObjPath:    objPath,
		MtlPath:    mtlPath,
		MattxtPath: matPath,
	}
	err = obj.Import(req)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}
	e.materials = req.Data.Materials
	e.triangles = req.Data.Triangles
	e.vertices = req.Data.Vertices

	return nil
}
