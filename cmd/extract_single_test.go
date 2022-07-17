package cmd

import (
	"os"
	"testing"
)

func TestExtractS3D(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	err := extractS3D("../s3d/test/eq/crushbone.s3d", "./test/eq/_crushbone.s3d", false)
	if err != nil {
		t.Fatalf("extract %s", err)
	}
}

func TestExtractEQG(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	err := extractEQG("test/eq/bazaar.eqg", "test/eq/_bazaar.eqg", false)
	if err != nil {
		t.Fatalf("extract %s", err)
	}
}
