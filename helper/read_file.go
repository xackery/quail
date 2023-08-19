package helper

import (
	"bufio"
	"fmt"
	"os"
)

// ReadFile loads a file and splits it into a string slice
func ReadFile(path string) ([]string, error) {
	lines := []string{}
	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
