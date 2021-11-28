package s3d

import "github.com/xackery/quail/common"

// Wld is a collection of fragments representing a world file
type Wld struct {
	IsOldWorld     bool
	ShortName      string
	FragmentCount  uint32
	BspRegionCount uint32
	Hash           map[int]string
	Fragments      []common.WldFragmenter
}
