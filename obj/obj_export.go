package obj

import (
	"fmt"
)

func Export(obj *ObjData, objPath string, mtlPath string, matTxtPath string) error {
	var err error
	err = exportMatTxt(obj, matTxtPath)
	if err != nil {
		return fmt.Errorf("exportMatTxt: %w", err)
	}
	err = exportMtl(obj, mtlPath)
	if err != nil {
		return fmt.Errorf("exportMtl: %w", err)
	}
	err = exportObjFile(obj, objPath)
	if err != nil {
		return fmt.Errorf("exportObjFile: %w", err)
	}

	return nil
}
