package quail

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wce"
)

func (q *Quail) wldRead(srcWld *raw.Wld, filename string) error {
	wld := wce.New(filename)
	err := wld.ReadWldRaw(srcWld)
	if err != nil {
		return fmt.Errorf("read wld: %w", err)
	}

	if strings.ToLower(filename) == "objects.wld" {
		q.WldObject = wld
	} else if strings.ToLower(filename) == "lights.wld" {
		q.WldLights = wld
	} else {
		q.Wld = wld
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
