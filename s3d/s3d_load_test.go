package s3d

import (
	"os"
	"testing"
)

func TestS3DLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &S3D{}
	f, err := os.Open("../eq/tmp/load.s3d")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

}
