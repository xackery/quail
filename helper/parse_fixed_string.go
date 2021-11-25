package helper

import (
	"fmt"
	"io"
)

func ParseFixedString(r io.Reader, size uint32) (string, error) {
	in := make([]byte, size)
	_, err := r.Read(in)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}
	final := ""
	for _, char := range in {
		final += string(char)
	}
	fmt.Println(len(final))
	return final, nil
}
