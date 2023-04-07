package obj

import (
	"fmt"

	"github.com/xackery/quail/model/geo"
)

func Import(req *ObjRequest) error {
	var err error
	if req == nil {
		return fmt.Errorf("request is nil")
	}
	req.Data = &ObjData{}
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

func materialByName(name string, obj *ObjData) *geo.Material {
	for _, mat := range obj.Materials {
		if name == mat.Name {
			return mat
		}
	}
	return nil
}
