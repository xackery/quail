package main

import (
	"fmt"
	"os"
	"regexp"
	"syscall/js"

	"github.com/xackery/quail/quail"
)

var q = quail.New()
var quailModule js.Value

func getFileExtension(path string) string {
	re := regexp.MustCompile(`\.[^./\\]+$`)
	match := re.FindString(path)
	return match
}

func convert(args []js.Value) error {
	if len(args) < 2 {
		fmt.Println("Expected >= 2 args for convert")
		return nil
	}
	srcPath := args[0].String()
	dstPath := args[1].String()

	srcExt := getFileExtension(srcPath)
	var err = fmt.Errorf("")
	switch srcExt {
	case ".quail":
		err = q.DirRead(srcPath)
		if err != nil {
			return fmt.Errorf("quail read dir: %w", err)
		}

	case ".json":
		err = q.JsonRead(srcPath)
		if err != nil {
			return fmt.Errorf("json read: %w", err)
		}
	default:
		_, e := os.ReadFile(srcPath)
		if e != nil {
			fmt.Println("Error reading file: " + e.Error())
		}
		err = q.PfsRead(srcPath)
		if err != nil {
			return fmt.Errorf("pfs read: %w", err)
		}
	}

	dstExt := getFileExtension(dstPath)
	switch dstExt {
	case ".quail":
		err = q.DirWrite(dstPath)
		if err != nil {
			return fmt.Errorf("dir write: %w", err)
		}
		return nil
	case ".json":
		err = q.JsonWrite(dstPath)

		if err != nil {
			return fmt.Errorf("json write: %w", err)
		}
	default:
		err = q.PfsWrite(1, 1, dstPath)
		if err != nil {
			return fmt.Errorf("pfs write: %w", err)
		}

	}
	return nil
}

func main() {
	fmt.Println("Initialized Quail")
	quailModule = (js.ValueOf(map[string]interface{}{
		"quail": js.ValueOf(true),
		"convert": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			err := convert(args)
			if err != nil {
				fmt.Println("Got error: " + err.Error())
				return 1
			}
			return 0
		}),
	}))
	select {}
}
