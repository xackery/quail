package helper

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
)

// Deflate takes a byte slice and compresses it down
func Deflate(in []byte) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	pos := 0
	remain := len(in)
	blockSize := 8192
	for remain > 0 {
		sz := int(0)
		if remain > blockSize {
			sz = blockSize
			remain -= blockSize
		} else {
			sz = remain
			remain = 0
		}

		//fmt.Println(pos+blockSize, "vs", sz)
		buf := bytes.NewBuffer(nil)

		//w, err := zlib.NewWriterLevel(buf, 2)
		w := zlib.NewWriter(buf)
		//w, err := zlib.NewWriterRaw(buf, 5, 0, 15, 8)
		//if err != nil {
		//	return nil, fmt.Errorf("newWriter: %w", err)
		//}

		inflateSize, err := w.Write(in[pos : pos+sz])
		if err != nil {
			return nil, fmt.Errorf("write: %w", err)
		}
		pos += sz

		w.Close()

		/*err = binary.Write(out, binary.LittleEndian, adler32.Checksum(buf.Bytes()))
		if err != nil {
			return nil, fmt.Errorf("checksum: %w", err)
		}*/

		//fmt.Println("deflate size", buf.Len())
		err = binary.Write(out, binary.LittleEndian, uint32(buf.Len()))
		if err != nil {
			return nil, fmt.Errorf("write deflateSize: %w", err)
		}

		//fmt.Println("inflate size", inflateSize)
		err = binary.Write(out, binary.LittleEndian, uint32(inflateSize))
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
