package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/common"
)

func TestInspectQuail(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()
	args := []string{
		fmt.Sprintf("%s/global3_chr.src.wld:1", dirTest),
		//fmt.Sprintf("%s/mim_chr.quail", dirTest),
	}

	testCmd := &cobra.Command{}
	testCmd.Flags().String("path", "", "path to inspect")
	testCmd.Flags().String("path2", "", "path to compare")

	err := runInspectE(testCmd, args)
	if err != nil {
		t.Fatal(err)
	}
}
