package wld

import "github.com/xackery/quail/common"

// WLD is a wld file struct
type WLD struct {
	IsOldWorld     bool
	ShortName      string
	FragmentCount  uint32
	BspRegionCount uint32
	Hash           map[int]string
	Fragments      []common.WldFragmenter
}
