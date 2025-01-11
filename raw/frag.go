package raw

import (
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/raw/rawfrag"
)

// FragName returns the name of a fragment
func FragName(fragCode int) string {
	return rawfrag.FragName(fragCode)
}

func NewFrag(r io.ReadSeeker) common.FragmentReadWriter {
	return rawfrag.NewFrag(r)
}
