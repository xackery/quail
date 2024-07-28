package quail

import (
	"fmt"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wld"
)

func (q *Quail) wldRead(srcWld *raw.Wld) error {
	q.wld = &wld.Wld{}
	err := q.wld.ReadRaw(srcWld)
	if err != nil {
		return fmt.Errorf("read wld: %w", err)
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
