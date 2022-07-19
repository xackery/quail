package helper

import (
	"fmt"
	"io"
	"strings"
)

func ReadFixedString(r io.Reader, size uint32) (string, error) {
	in := make([]byte, size)
	_, err := r.Read(in)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}
	final := ""
	for _, char := range in {
		if char == 0x0 {
			continue
		}
		final += string(char)
	}
	final = strings.TrimSpace(final)
	return final, nil
}
