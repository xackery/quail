package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/helper"
)

func TestDoubleZipQuail(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()

	keyword := "alkabormare"
	ext := "eqg"

	os.Chdir(dirTest)

	var cmd *cobra.Command

	fmt.Printf("quail unzip %s.%s\n", keyword, ext)
	err := runUnzipE(cmd, []string{
		fmt.Sprintf("%s/%s.%s", eqPath, keyword, ext),
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("quail zip _%s.%s\n", keyword, ext)
	err = runZipE(cmd, []string{
		fmt.Sprintf("%s/_%s.%s", dirTest, keyword, ext),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = os.Rename(fmt.Sprintf("%s/%s.%s", dirTest, keyword, ext), fmt.Sprintf("%s/%s2.%s", dirTest, keyword, ext))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("quail unzip %s2.%s\n", keyword, ext)
	err = runUnzipE(cmd, []string{
		fmt.Sprintf("%s/%s2.%s", dirTest, keyword, ext),
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("quail zip _%s2.%s\n", keyword, ext)
	err = runZipE(cmd, []string{
		fmt.Sprintf("%s/_%s2.%s", dirTest, keyword, ext),
	})
	if err != nil {
		t.Fatal(err)
	}

}
