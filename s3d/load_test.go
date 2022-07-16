package s3d

import (
	"fmt"
	"os"
	"testing"
)

func TestS3DLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	eqgFile := "test/eq/crushbone.s3d"
	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	f, err := os.Open(eqgFile)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	for _, file := range e.Files() {
		fmt.Println(file.Name())
	}

}
