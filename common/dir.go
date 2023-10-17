package common

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDir returns a testing directory
func DirTest(t *testing.T) string {
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
					t.Fatalf("failed to find go.mod: %s", err.Error())
				}
			}
		}
	}
	dir, err := filepath.Abs(strings.ReplaceAll(path, "go.mod", "") + "/test")
	if err != nil {
		t.Fatalf("failed to get test path: %s", err.Error())
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("failed to create test path: %s", err.Error())
	}
	return dir
}
