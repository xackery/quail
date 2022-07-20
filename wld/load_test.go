package wld

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/s3d"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/eq/crushbone.s3d"
	file := "crushbone.wld"

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	defer f.Close()

	archive, err := s3d.NewFile(path)
	if err != nil {
		t.Fatalf("s3d new: %s", err)
	}

	for _, fe := range archive.Files() {
		fmt.Println(fe.Name())
	}

	dump.New(path)
	defer dump.WriteFileClose(path)
	e, err := NewFile("crushbone", archive, file)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	fmt.Println(e.name)

}
