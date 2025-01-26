package helper

import (
	"encoding/hex"
	"fmt"
)

// ByteCompareTest compares two byte arrays, and returns an error if they are not equal. This was made for testing
func ByteCompareTest(src []byte, dst []byte) error {

	if len(dst) < len(src) {
		min := 0
		max := len(src)
		fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(src[min:max]))
		min = 0
		max = len(dst)
		fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dst[min:max]))
		return fmt.Errorf("dst is too small by %d bytes", len(src)-len(dst))
	}

	for j := 0; j < len(dst); j++ {
		if len(src) <= j {
			min := 0
			max := len(src)
			fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(src[min:max]))
			max = len(dst)
			fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dst[min:max]))
			return fmt.Errorf("src eof at offset %d (dst is too large by %d bytes)", j, len(dst)-len(src))
		}
		if len(dst) <= j {
			return fmt.Errorf("dst eof at offset %d (dst is too small by %d bytes)", j, len(src)-len(dst))
		}
		if dst[j] == src[j] {
			continue
		}
		fmt.Printf("mismatch at offset %d (src: 0x%x vs dst: 0x%x aka %d)\n", j, src[j], dst[j], dst[j])
		max := j + 16
		if max > len(src) {
			max = len(src)
		}

		min := j - 16
		if min < 0 {
			min = 0
		}
		fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(src[min:max]))
		if max > len(dst) {
			max = len(dst)
		}

		fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dst[min:max]))
		return fmt.Errorf("%d data mismatch", j)
	}
	return nil
}
