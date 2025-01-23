package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
)

func TestInspectQuail(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	baseName := "global3_chr"

	wldName := fmt.Sprintf("%s.wld", baseName)

	archive, err := pfs.NewFile(fmt.Sprintf("%s/%s.s3d", eqPath, baseName))
	if err != nil {
		t.Fatalf("failed to open s3d %s: %s", baseName, err.Error())
	}
	defer archive.Close()

	// get wld
	data, err := archive.File(wldName)
	if err != nil {
		t.Fatalf("failed to open wld %s: %s", baseName, err.Error())
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s.src.wld", dirTest, baseName), data, 0644)
	if err != nil {
		t.Fatalf("failed to write wld %s: %s", baseName, err.Error())
	}

	args := []string{
		fmt.Sprintf("%s/global3_chr.src.wld:1", dirTest),
		//fmt.Sprintf("%s/mim_chr.quail", dirTest),
	}

	testCmd := &cobra.Command{}
	testCmd.Flags().String("path", "", "path to inspect")
	testCmd.Flags().String("path2", "", "path to compare")

	err = runInspectE(testCmd, args)
	if err != nil {
		t.Fatal(err)
	}
}
