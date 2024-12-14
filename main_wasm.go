package main

import (
	"fmt"
	"regexp"
	"syscall/js"

	"github.com/xackery/quail/os"
	"github.com/xackery/quail/quail"
)

var q = quail.New()

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

	// We are passing in a Uint8array to copy into the fs for processing here.
	// Probably make this an independent function for writing to the FS
	if len(args) == 3 {
		buffer := make([]byte, args[2].Length())
		js.CopyBytesToGo(buffer, args[2])
		err := os.WriteFile(srcPath, buffer, 0755)
		if err != nil {
			fmt.Println("Could not write file: " + err.Error())
			return err
		}
		fmt.Println("Wrote srcPath" + srcPath)
	}

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
		} else {
			fmt.Println("Did open file: " + srcPath)
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

func quailConvert(this js.Value, args []js.Value) interface{} {
	err := convert(args)
	if err != nil {
		fmt.Println("Got error: " + err.Error())
		return 1
	}
	return 0
}

func quailFs(this js.Value, args []js.Value) interface{} {
	return os.ExportFileSystem()
}

func main() {
	fmt.Println("Initialized Quail")
	js.Global().Set("quailConvert", js.FuncOf(quailConvert))
	js.Global().Set("quailFs", js.FuncOf(quailFs))
	select {}
}
