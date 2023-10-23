package quail

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/wld"
	"github.com/xackery/quail/pfs"
)

// Decode decodes a WLD file
func WLDDecode2(r io.ReadSeeker, pfs *pfs.PFS) ([]*common.Model, error) {
	world := &common.Wld{}
	err := wld.Decode(world, r)
	if err != nil {
		return nil, fmt.Errorf("wld decode: %w", err)
	}

	return nil, fmt.Errorf("not implemented")
}
