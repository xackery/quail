package model

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/mesh/ter"
)

// Encode encodes a model to a writer based on the filetype
func Encode(model *common.Model, version uint32, w io.Writer) error {
	switch model.FileType {
	case "mod":
		return mod.Encode(model, version, w)
	case "mds":
		return mds.Encode(model, version, w)
	case "ter":
		return ter.Encode(model, version, w)
	}
	return fmt.Errorf("unknown file type: %s", model.FileType)
}
