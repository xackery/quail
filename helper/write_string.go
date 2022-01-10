package helper

import (
	"fmt"
	"io"
)

// WriteString is used to write a CString (zero terminated) to w
func WriteString(w io.Writer, in string) error {
	in += "\000"
	//fmt.Println(hex.Dump([]byte(in)))
	_, err := w.Write([]byte(in))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
