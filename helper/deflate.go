package helper

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
)

func Deflate(in []byte) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	lastPos := 0
	pos := 0
	sz := len(in)
	for sz > pos {
		if pos+8192 > sz {
			pos = sz
		} else {
			pos += 8192
		}
		buf := bytes.NewBuffer(nil)
		w := zlib.NewWriter(buf)

		total, err := w.Write(in[lastPos:pos])
		if err != nil {
			return nil, fmt.Errorf("write: %w", err)
		}

		err = w.Flush()
		if err != nil {
			return nil, fmt.Errorf("flush: %w", err)
		}
		lastPos = pos
		w.Close()
		if err != nil {
			return nil, fmt.Errorf("close: %w", err)
		}

		err = binary.Write(out, binary.LittleEndian, uint32(buf.Len()))
		if err != nil {
			return nil, fmt.Errorf("write deflateSize: %w", err)
		}

		err = binary.Write(out, binary.LittleEndian, uint32(total))
		if err != nil {
			return nil, fmt.Errorf("write sz: %w", err)
		}

		err = binary.Write(out, binary.LittleEndian, buf.Bytes())
		if err != nil {
			return nil, fmt.Errorf("write block: %w", err)
		}
	}
	return out.Bytes(), nil
}
