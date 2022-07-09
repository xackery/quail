package obj

import (
	"fmt"

	"github.com/xackery/quail/common"
)

func Import(objPath string, mtlPath string, matTxtPath string) (*ObjData, error) {
	obj := &ObjData{}
	err := objImport(obj, objPath, mtlPath, matTxtPath)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func objImport(obj *ObjData, objPath string, mtlPath string, matTxtPath string) error {
	var err error
	err = mattxtImport(obj, matTxtPath)
	if err != nil {
		return fmt.Errorf("importMatTxt: %w", err)
	}
	err = mtlImport(obj, mtlPath)
	if err != nil {
		return fmt.Errorf("importMatTxt: %w", err)
	}
	err = importFile(obj, objPath)
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
