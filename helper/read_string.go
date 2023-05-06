package helper

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ReadString is used to read a CString (zero terminated) from r
func ReadString(r io.ReadSeeker) (string, error) {
	var buf uint8
	var err error
	value := ""
	for {
		err = binary.Read(r, binary.LittleEndian, &buf)
		if err != nil {
			return "", fmt.Errorf("read: %w", err)
		}
		if buf == 0x00 {
			return value, nil
		}
		value += string(buf)
	}
}
