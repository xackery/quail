//go:build test_all
// +build test_all

package wce_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAllWldFiles(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Fatalf("EQ_PATH not set")
	}

	start := time.Now()

	// walk eqPath, find .s3d's, add them to tests
	tests = []testEntry{}

	files, err := os.ReadDir(eqPath)
	if err != nil {
		t.Fatalf("os.ReadDir(%q) failed: %v", eqPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if len(file.Name()) < 4 {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext != ".s3d" {
			continue
		}

		baseName := file.Name()[:len(file.Name())-4]
		tests = append(tests, testEntry{baseName: baseName})
	}

	fmt.Printf("Found %d .s3d files\n", len(tests))

	TestStep4(t)
	if t.Failed() {
		t.Fatalf("TestStep4 failed")
	}

	fmt.Printf("Took %0.2f seconds for %d total tests\n", time.Since(start).Seconds(), len(tests))
}
