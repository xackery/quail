package helper

import (
	"fmt"
	"io"
)

// WriteFixedString will write a string and pad 0x00 until size is satisfied
func WriteFixedString(w io.Writer, in string, size int) error {
	i, err := w.Write([]byte(in))
	if err != nil {
		return fmt.Errorf("write string: %w", err)
	}
	for ; i < size; i++ {
		_, err = w.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("write pad: %w", err)
		}
	}
	return nil
}
