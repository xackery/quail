package eqg

import (
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// EQG represents a modern everquest zone archive format
type EQG struct {
	files []*common.FileEntry
}

type byCRC []*common.FileEntry

func (s byCRC) Len() int {
	return len(s)
}

func (s byCRC) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byCRC) Less(i, j int) bool {
	return helper.FilenameCRC32(s[i].Name) < helper.FilenameCRC32(s[j].Name)
}
