package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestExportS3D(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	err := exportS3D("../s3d/test/eq/crushbone.s3d", "./test/eq/_crushbone.s3d")
	if err != nil {
		t.Fatalf("extract %s", err)
	}
}

func TestExportEQG(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	out := "test/eq/_bazaar.eqg"
	fi, err := os.Stat(out)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("creating directory %s/\n", out)
			err = os.MkdirAll(out, 0766)
			if err != nil {
				t.Fatalf("mkdirall: %s", err)
			}
		}
		fi, err = os.Stat(out)
		if err != nil {
			t.Fatalf("stat after mkdirall: %s", err)
		}
	}
	if !fi.IsDir() {
		t.Fatalf("%s is not a directory", out)
	}

	err = exportEQG("test/eq/bazaar.eqg", out)
	if err != nil {
		t.Fatalf("extract %s", err)
	}
}
