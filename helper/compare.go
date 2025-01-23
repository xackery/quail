package helper

import (
	"fmt"
	"strings"
)

func Compare(val, val2 interface{}) error {
	switch val := val.(type) {
	case string:
		val2, ok := val2.(string)
		if !ok {
			return fmt.Errorf("type mismatch (%T != %T)", val, val2)
		}

		return compareString(val, val2)
	}

	valBuf := fmt.Sprintf("%#v", val)
	val2Buf := fmt.Sprintf("%#v", val2)

	if valBuf == val2Buf {
		fmt.Println("Files are identical")
		return nil
	}

	lines := strings.Split(valBuf, "\n")
	lines2 := strings.Split(val2Buf, "\n")
	for i := 0; i < len(lines); i++ {
		if i >= len(lines2) {
			fmt.Printf("-%s\n", lines[i])
			continue
		}
		if lines[i] != lines2[i] {
			fmt.Printf("-%s\n", lines[i])
			fmt.Printf("+%s\n", lines2[i])
		}
	}

	return nil
}

func compareString(val, val2 string) error {
	if val == val2 {
		fmt.Println("Files are identical")
		return nil
	}

	lines := strings.Split(val, "\n")
	lines2 := strings.Split(val2, "\n")
	for i := 0; i < len(lines); i++ {
		if i >= len(lines2) {
			fmt.Printf("-%s\n", lines[i])
			continue
		}
		if lines[i] != lines2[i] {
			fmt.Printf("-%s\n", lines[i])
			fmt.Printf("+%s\n", lines2[i])
		}
	}

	return nil
}
