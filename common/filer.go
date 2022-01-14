package common

import "github.com/xackery/quail/helper"

// Filer is an interface that file-like structs fit inside
type Filer interface {
	Name() string
	Data() []byte
}

type ByCRC []Filer

func (s ByCRC) Len() int {
	return len(s)
}

func (s ByCRC) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCRC) Less(i, j int) bool {
	return helper.FilenameCRC32(s[i].Name()) < helper.FilenameCRC32(s[j].Name())
}
