package eqg

import "github.com/xackery/quail/helper"

// EQG represents a modern everquest zone archive format
type EQG struct {
	Files []*fileEntry
}

type fileEntry struct {
	name string
	data []byte
}

type byCRC []*fileEntry

func (s byCRC) Len() int {
	return len(s)
}

func (s byCRC) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byCRC) Less(i, j int) bool {
	return helper.FilenameCRC32(s[i].name) < helper.FilenameCRC32(s[j].name)
}
