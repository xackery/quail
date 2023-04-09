package s3d

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestS3DDecode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	eqgFile := fmt.Sprintf("%s/shp_chr.s3d", eqPath)
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
