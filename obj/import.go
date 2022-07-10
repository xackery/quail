package obj

import (
	"fmt"

	"github.com/xackery/quail/common"
)

func Import(req *ObjRequest) error {
	var err error
	if req == nil {
		return fmt.Errorf("request is nil")
	}
	req.Obj = &ObjData{}
	err = mattxtImport(req)
	if err != nil {
		return fmt.Errorf("importMatTxt: %w", err)
	}
	err = mtlImport(req)
	if err != nil {
		return fmt.Errorf("importMatTxt: %w", err)
	}
	err = importFile(req)
	if err != nil {
		return fmt.Errorf("importObjFile: %w", err)
	}

	return nil
}

func materialByName(name string, obj *ObjData) *common.Material {
	for _, mat := range obj.Materials {
		if name == mat.Name {
			return mat
		}
	}
	return nil
}
