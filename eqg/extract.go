package eqg

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/xackery/quail/helper"
)

func (e *EQG) Extract(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return err
		}
		fmt.Println("creating directory", path)
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return fmt.Errorf("mkdirall: %w", err)
		}
		fi, err = os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat after mkdirall: %w", err)
		}
	}
	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	extractStdout := ""
	for _, file := range e.files {
		fmt.Println("extract:", file.Name())
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", path, file.Name()), file.Data(), 0644)
		if err != nil {
			return err
		}
		extractStdout += file.Name() + ", "
	}
	if len(e.files) == 0 {
		return fmt.Errorf("no files found to extract")
	}
	extractStdout = extractStdout[0:len(extractStdout)-2] + "\n"
	fmt.Printf("extracted %d file%s to %s: %s", len(e.files), helper.Pluralize(len(e.files)), path, extractStdout)
	return nil
}
