package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// wldOpen prepares a wld file, loading fragments. This is usually then called by Decode
func loadFragments(fragmentCount uint32, r io.ReadSeeker) (fragments [][]byte, err error) {

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(fragmentCount); fragOffset++ {

		fragSize := dec.Uint32()
		totalFragSize += fragSize

		fragCode := dec.Bytes(4)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("frag position seek %d/%d: %w", fragOffset, fragmentCount, err)
		}
		data := make([]byte, fragSize)
		_, err = r.Read(data)
		if err != nil {
			return nil, fmt.Errorf("read frag %d/%d: %w", fragOffset, fragmentCount, err)
		}

		data = append(fragCode, data...)

		fragments = append(fragments, data)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek end of frag %d/%d: %w", fragOffset, fragmentCount, err)
		}
	}

	if dec.Error() != nil {
		return nil, fmt.Errorf("decode: %w", dec.Error())
	}
	return fragments, nil
}
