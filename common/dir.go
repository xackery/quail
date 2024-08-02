package common

import (
	"os"
	"path/filepath"
	"strings"
)

func IsTestExtensive() bool {
	return os.Getenv("IS_TEST_EXTENSIVE") == "1"
}

func IsTest() bool {
	return os.Getenv("IS_TEST") == "1"
}

// TestDir returns a testing directory
func DirTest() string {
	if !IsTest() {
		return "."
	}
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
	return dir
}
