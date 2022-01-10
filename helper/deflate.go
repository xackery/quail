package helper

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
)

func Deflate(in []byte) ([]byte, error) {
	out := bytes.NewBuffer(nil)

	pos := uint32(0)
	remain := uint32(len(in))
	maxBlockSize := uint32(8192)

	zw, err := flate.NewWriter(out, flate.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("newWriter: %w", err)
	}
	for remain > 0 {
		sz := uint32(0)
		if remain >= maxBlockSize {
			sz = maxBlockSize
			remain -= maxBlockSize
		} else {
			sz = remain
			remain = 0
		}

		blockLength := sz + 128
		block := make([]byte, blockLength)

		deflateSize, err := zw.Write(block)
		if err != nil {
			return nil, fmt.Errorf("write: %w", err)
		}
		pos += sz

		err = binary.Write(out, binary.LittleEndian, uint32(deflateSize))
		if err != nil {
			return nil, fmt.Errorf("write deflateSize: %w", err)
		}

		err = binary.Write(out, binary.LittleEndian, sz)
		if err != nil {
			return nil, fmt.Errorf("write sz: %w", err)
		}

		err = binary.Write(out, binary.LittleEndian, block)
		if err != nil {
			return nil, fmt.Errorf("write block: %w", err)
		}

		err = zw.Flush()
		if err != nil {
			return nil, fmt.Errorf("flush: %w", err)
		}

	}

	return out.Bytes(), nil
}
