package pfs

import (
	"strings"

	"github.com/xackery/quail/helper"
)

// filer is an interface that file-like structs fit inside
type filer interface {
	Name() string
	Data() []byte
}

// FilerByCRC sorts a slice of Filer by CRC32
type FilerByCRC []*FileEntry

// Len returns the length of the slice
func (s FilerByCRC) Len() int {
	return len(s)
}

// Swap swaps two elements in the slice
func (s FilerByCRC) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less returns true if the CRC32 of the first element is less than the second
func (s FilerByCRC) Less(i, j int) bool {
	return helper.FilenameCRC32(s[i].Name()) < helper.FilenameCRC32(s[j].Name())
}

// FilerByName sorts a slice of Filer by name
type FilerByName []filer

// Len returns the length of the slice
func (s FilerByName) Len() int {
	return len(s)
}

// Swap swaps two elements in the slice
func (s FilerByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less returns true if the name of the first element is less than the second
func (s FilerByName) Less(i, j int) bool {
	return strings.ToLower(s[i].Name()) < strings.ToLower(s[j].Name())
}
