package cmd

import (
	"os"
	"testing"
)

func TestImportArenaEQG(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	err := importPath("test/_arena.eqg", "test/eq/arena.eqg")
	if err != nil {
		t.Fatalf("import %s", err)
	}
}
