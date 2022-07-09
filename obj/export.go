package obj

import (
	"fmt"
)

func Export(obj *ObjData, objPath string, mtlPath string, matTxtPath string) error {
	var err error
	err = mattxtExport(obj, matTxtPath)
	if err != nil {
		return fmt.Errorf("exportMatTxt: %w", err)
	}
	err = mtlExport(obj, mtlPath)
	if err != nil {
		return fmt.Errorf("exportMtl: %w", err)
	}
	err = exportFile(obj, objPath)
	if err != nil {
		return fmt.Errorf("exportObjFile: %w", err)
	}

	return nil
}
