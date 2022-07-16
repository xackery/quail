package cmd

import (
	"os"
	"testing"
)

func TestExtractSingleTest(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	var err error
	err = os.WriteFile("./test/1dirtfloor.bmp", []byte{}, 0644)
	if err != nil {
		t.Fatalf("writeFile: %s", err)
	}
	err = extractS3D("../s3d/test/eq/crushbone.s3d", "./test/eq/_crushbone.s3d", false)
	if err != nil {
		t.Fatalf("extract %s", err)
	}
}
