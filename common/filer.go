package common

import (
	"strings"

	"github.com/xackery/quail/helper"
)

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

type ByName []Filer

func (s ByName) Len() int {
	return len(s)
}

func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByName) Less(i, j int) bool {
	return strings.ToLower(s[i].Name()) < strings.ToLower(s[j].Name())
}
