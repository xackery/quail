package s3d

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestS3DDecode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	eqgFile := "test/eq/crushbone.s3d"
	dump.New(eqgFile)
	defer dump.WriteFileClose(eqgFile)
	e, err := NewFile(eqgFile)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	for _, file := range e.Files() {
		fmt.Println(file.Name())
	}

}
