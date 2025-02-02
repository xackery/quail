package helper

import (
	"os"
	"path/filepath"
	"strings"
)

// TestDir returns a testing directory
func DirTest() string {
	path := "go.mod"

	// look for go.mod
	_, err := os.Stat(path)
	if err != nil {
		path = "../go.mod"
		_, err = os.Stat(path)
		if err != nil {
			path = "../../go.mod"
			_, err = os.Stat(path)
			if err != nil {
				path = "../../../go.mod"
				_, err = os.Stat(path)
				if err != nil {
					return "."
				}
			}
		}
	}
	dir, err := filepath.Abs(strings.ReplaceAll(path, "go.mod", "") + "/test")
	if err != nil {
		return "."
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "."
	}
	os.Chdir(dir)
	return dir
}
