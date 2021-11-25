package s3d

import (
	"os"
	"testing"
)

func TestS3DSave(t *testing.T) {
	e := &S3D{}
	f, err := os.Create("test/save.s3d")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	err = e.Save(f)
	if err != nil {
		t.Fatalf("save: %v", err)
	}

}
