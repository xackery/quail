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

type FilerByCRC []Filer

func (s FilerByCRC) Len() int {
	return len(s)
}

func (s FilerByCRC) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FilerByCRC) Less(i, j int) bool {
	return helper.FilenameCRC32(s[i].Name()) < helper.FilenameCRC32(s[j].Name())
}

type FilerByName []Filer

func (s FilerByName) Len() int {
	return len(s)
}

func (s FilerByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FilerByName) Less(i, j int) bool {
	return strings.ToLower(s[i].Name()) < strings.ToLower(s[j].Name())
}
