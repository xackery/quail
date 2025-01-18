package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
)

func TestConvertQuail(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()
	args := []string{
		fmt.Sprintf("%s/mim_chr.s3d", eqPath),
		fmt.Sprintf("%s/mim_chr.quail", dirTest),
	}
	err := runConvertE(args)
	if err != nil {
		t.Fatal(err)
	}
}
