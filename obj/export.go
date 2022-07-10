package obj

import (
	"fmt"
)

func Export(req *ObjRequest) error {
	var err error
	err = mattxtExport(req)
	if err != nil {
		return fmt.Errorf("exportMatTxt: %w", err)
	}
	err = mtlExport(req)
	if err != nil {
		return fmt.Errorf("exportMtl: %w", err)
	}
	err = objExport(req)
	if err != nil {
		return fmt.Errorf("exportObjFile: %w", err)
	}

	return nil
}
