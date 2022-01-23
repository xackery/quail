package wld

import "github.com/xackery/quail/common"

// WLD is a wld file struct
type WLD struct {
	ShortName      string
	BspRegionCount uint32
	Hash           map[int]string
	fragments      []*fragmentInfo
}

type fragmentInfo struct {
	name string
	data common.WldFragmenter
}
