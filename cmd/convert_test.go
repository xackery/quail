package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/helper"
)

func TestConvertQuail(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	testCmd := &cobra.Command{}
	dirTest := helper.DirTest()
	args := []string{
		//fmt.Sprintf("%s/dbx.eqg", eqPath),
		//fmt.Sprintf("%s/dbx.quail", dirTest),
		fmt.Sprintf("%s/dbx.quail", dirTest), //eqPath),
		fmt.Sprintf("%s/dbx.eqg", dirTest),
	}
	err := runConvertE(testCmd, args)
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
	testCmd := &cobra.Command{}

	keyword := "iceclad"
	ext := "s3d"

	fmt.Printf("quail convert %s.%s %s.quail\n", keyword, ext, keyword)
	err := runConvertE(testCmd, []string{
		fmt.Sprintf("%s/%s.%s", eqPath, keyword, ext),
		fmt.Sprintf("%s/%s.quail", dirTest, keyword),
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("quail convert %s.quail %s.%s\n", keyword, keyword, ext)
	err = runConvertE(testCmd, []string{
		fmt.Sprintf("%s/%s.quail", dirTest, keyword),
		fmt.Sprintf("%s/%s.%s", dirTest, keyword, ext),
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("quail convert %s.%s %s2.quail\n", keyword, ext, keyword)
	err = runConvertE(testCmd, []string{
		fmt.Sprintf("%s/%s.%s", dirTest, keyword, ext),
		fmt.Sprintf("%s/%s2.quail", dirTest, keyword),
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("quail convert %s2.quail %s2.%s\n", keyword, keyword, ext)
	err = runConvertE(testCmd, []string{
		fmt.Sprintf("%s/%s2.quail", dirTest, keyword),
		fmt.Sprintf("%s/%s2.%s", dirTest, keyword, ext),
	})
	if err != nil {
		t.Fatal(err)
	}

}
