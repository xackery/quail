package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/helper"
)

func TestConvertQuail(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()
	args := []string{
		//fmt.Sprintf("%s/dbx.eqg", eqPath),
		//fmt.Sprintf("%s/dbx.quail", dirTest),
		fmt.Sprintf("%s/dbx.quail", dirTest), //eqPath),
		fmt.Sprintf("%s/dbx.eqg", dirTest),
	}
	err := runConvertE(args)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoubleConvertQuail(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()

	err := runConvertE([]string{
		fmt.Sprintf("%s/dbx.eqg", eqPath),
		fmt.Sprintf("%s/dbx.quail", dirTest),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = runConvertE([]string{
		fmt.Sprintf("%s/dbx.quail", dirTest),
		fmt.Sprintf("%s/dbx.eqg", dirTest),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = runConvertE([]string{
		fmt.Sprintf("%s/dbx.eqg", dirTest),
		fmt.Sprintf("%s/dbx2.quail", dirTest),
	})
	if err != nil {
		t.Fatal(err)
	}
}
