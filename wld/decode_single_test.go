package wld

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/s3d"
)

func TestDecode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/eq/crushbone.s3d"
	file := "crushbone.wld"

	archive, err := s3d.NewFile(path)
	if err != nil {
		t.Fatalf("s3d new: %s", err)
	}

	dump.New(file)
	defer dump.WriteFileClose(path + "_" + file)
	e, err := NewFile("crushbone", archive, file)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	fmt.Println(e.name)

}
