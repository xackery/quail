package raw

import (
	"io"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/raw/rawfrag"
)

// FragName returns the name of a fragment
func FragName(fragCode int) string {
	return rawfrag.FragName(fragCode)
}

func NewFrag(r io.ReadSeeker) helper.FragmentReadWriter {
	return rawfrag.NewFrag(r)
}
