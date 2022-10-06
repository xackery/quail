package s3d

import (
	"fmt"
	"os"
	"testing"
)

func TestS3DSingleDecode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	s3dFile := "test/eq/global17_amr.s3d"
	//dump.New(s3dFile)
	//defer dump.WriteFileClose(s3dFile)
	e, err := NewFile(s3dFile)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	for _, file := range e.Files() {
		fmt.Println(file.Name())
	}

}
