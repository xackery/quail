package s3d

import (
	"fmt"
	"os"

	"github.com/xackery/quail/helper"
)

func (e *S3D) Extract(path string) (string, error) {

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("creating directory %s/\n", path)
			err = os.MkdirAll(path, 0766)
			if err != nil {
				return "", fmt.Errorf("mkdirall: %w", err)
			}
		}
		fi, err = os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("stat after mkdirall: %w", err)
		}
	}
	if !fi.IsDir() {
		return "", fmt.Errorf("%s is not a directory", path)
	}

	extractStdout := ""
	for i, file := range e.files {
		//fmt.Println(fmt.Sprintf("%s/%s", path, file.Name()), len(file.Data()))
		err = os.WriteFile(fmt.Sprintf("%s/%s", path, file.Name()), file.Data(), 0644)
		if err != nil {
			return "", fmt.Errorf("index %d: %w", i, err)
		}
		extractStdout += file.Name() + ", "
	}
	if len(e.files) == 0 {
		return "", fmt.Errorf("no files found to extract")
	}
	extractStdout = extractStdout[0 : len(extractStdout)-2]
	return fmt.Sprintf("extracted %d file%s to %s: %s", len(e.files), helper.Pluralize(len(e.files)), path, extractStdout), nil
}
