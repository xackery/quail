package quail

import (
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/wld"
	"github.com/xackery/quail/pfs"
)

// Decode decodes a WLD file
func (q *Quail) WLDDecode(r io.ReadSeeker, pfs *pfs.PFS) (*common.Wld, error) {
	world := common.NewWld("")
	err := wld.Decode(world, r)
	if err != nil {
		return nil, fmt.Errorf("wld decode: %w", err)
	}

	err = q.WldUnmarshal(world)
	if err != nil {
		return nil, fmt.Errorf("wld import: %w", err)
	}

	return world, nil
}
