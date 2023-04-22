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

	name := "shp_chr.s3d"
	eqgFile := fmt.Sprintf("%s/%s", eqPath, name)
	dump.New(eqgFile)
	defer dump.WriteFileClose("test/" + name)
	e, err := NewFile(eqgFile)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	for _, file := range e.Files() {
		fmt.Println(file.Name())
	}

}
