package helper

import (
	"bytes"
	"fmt"

	"compress/zlib"
	//"github.com/xackery/go-zlib"
)

func Inflate(in []byte) ([]byte, error) {
	buf := bytes.NewBuffer(in)
	r, err := zlib.NewReader(buf)
	if err != nil {
		return nil, fmt.Errorf("newReader: %w", err)
	}
	data := []byte{}
	_, err = r.Read(data)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	//fmt.Println("deflated:", len(data), ":", hex.Dump(data))
	return data, nil
}
