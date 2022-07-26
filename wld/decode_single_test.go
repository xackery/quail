package wld

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/s3d"
)

func TestDecode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	category := "goo_chr"
	path := "test/eq/" + category + ".s3d"
	file := category + ".wld"

	archive, err := s3d.NewFile(path)
	if err != nil {
		t.Fatalf("s3d new: %s", err)
	}

	dump.New(file)
	defer dump.WriteFileClose(path + "_" + file)
	e, err := NewFile(category, archive, file)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	if len(e.materials) != 42 {
		t.Fatalf("wanted 42 materials, got %d", len(e.materials))
	}

	if len(e.meshes) != 2694 {
		t.Fatalf("wanted 2694 meshes, got %d", len(e.meshes))
	}

}
