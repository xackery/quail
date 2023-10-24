package quail

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/wld"
	"github.com/xackery/quail/pfs"
)

// Decode decodes a WLD file
func WLDDecode(r io.ReadSeeker, pfs *pfs.PFS) (*common.Wld, error) {
	world := common.NewWld("")
	err := wld.Decode(world, r)
	if err != nil {
		return nil, fmt.Errorf("wld decode: %w", err)
	}

	err = wld.Convert(world)
	if err != nil {
		return nil, fmt.Errorf("wld convert: %w", err)
	}

	return world, nil
}
