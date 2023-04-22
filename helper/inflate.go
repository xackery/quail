package helper

import (
	"bytes"
	"fmt"
	"io"

	"compress/zlib"
)

func Inflate(in []byte, size int) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, fmt.Errorf("newReader: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, r)
	if err != nil {
		return nil, fmt.Errorf("copy: %w", err)
	}

	//fmt.Println("inflated\n", hex.Dump(out))
	return buf.Bytes(), nil
}
