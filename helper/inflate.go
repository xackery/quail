package helper

import (
	"fmt"

	"github.com/xackery/go-zlib"
)

func Inflate(in []byte, size int) ([]byte, error) {
	r, err := zlib.NewReader(nil)
	if err != nil {
		return nil, fmt.Errorf("newReader: %w", err)
	}
	out := make([]byte, size)
	_, _, err = r.ReadBuffer(in, out)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	//fmt.Println("inflated\n", hex.Dump(out))
	return out, nil
}
