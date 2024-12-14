package helper

import "github.com/xackery/quail/os"

// IsFile returns true if path is a file
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}
