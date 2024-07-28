package quail

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wld"
)

func (q *Quail) wldRead(srcWld *raw.Wld, filename string) error {

	wld := &wld.Wld{
		FileName: filename,
	}
	err := wld.ReadRaw(srcWld)
	if err != nil {
		return fmt.Errorf("read wld: %w", err)
	}

	if strings.ToLower(filename) == "objects.wld" {
		q.wldObject = wld
	} else if strings.ToLower(filename) == "lights.wld" {
		q.wldLights = wld
	} else {
		q.wld = wld
	}

	return nil
}

func (q *Quail) ModelByName(name string) *common.Model {
	for _, model := range q.Models {
		if model.Header.Name == name {
			return model
		}
	}
	return nil
}
