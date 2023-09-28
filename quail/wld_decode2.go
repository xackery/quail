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
	e, err := common.WldOpen(r)
	if err != nil {
		return nil, fmt.Errorf("wld open: %w", err)
	}
	defer e.Close()

	err = wld.Decode(e)
	if err != nil {
		return nil, fmt.Errorf("wld decode: %w", err)
	}

	return e.Models, nil
}
