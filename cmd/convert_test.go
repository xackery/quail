package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/helper"
)

func TestConvertQuail(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
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
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	dirTest := helper.DirTest()
	testCmd := &cobra.Command{}

	keyword := "arcstone"
	ext := "eqg"

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

func TestDoubleConvertQuailDir(t *testing.T) {
	totalTime := time.Now()

	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	dirTest := helper.DirTest()
	testCmd := &cobra.Command{}

	dirPaths := []string{
		"/src/eq/rof2",
		"/src/eq/takp",
		"/src/eq/ls",
	}

	//nextPath := "/src/eq/rof2/B09.eqg"
	nextPath := ""

	isNextPathFound := false
	if nextPath == "" {
		isNextPathFound = true
	}
	for _, dirPath := range dirPaths {
		fmt.Printf("dirPath: %s\n", dirPath)

		err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				t.Fatalf("error walking dir: %v", err)
			}
			if d.IsDir() {
				return nil
			}
			if !isNextPathFound {
				nextPathTemp := filepath.Join(dirPath, d.Name())
				if nextPathTemp != nextPath {
					return nil
				}
				isNextPathFound = true
			}

			ext := filepath.Ext(d.Name())
			keyword := d.Name()[:len(d.Name())-len(ext)]
			if len(ext) > 1 {
				ext = ext[1:]
			}
			if ext != "s3d" && ext != "eqg" {
				return nil
			}
			fmt.Println(d.Name())

			fmt.Printf("quail convert %s/%s %s.quail\n", dirPath, d.Name(), keyword)
			err = runConvertE(testCmd, []string{
				fmt.Sprintf("%s/%s.%s", dirPath, keyword, ext),
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
			return nil
		})
		if err != nil {
			t.Fatalf("error walking dir: %v", err)
		}
	}
	fmt.Printf("Total time: %0.2fs\n", time.Since(totalTime).Seconds())

}
