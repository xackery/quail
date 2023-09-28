package pfs

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkEQG(b *testing.B) {
	eqgPath := os.Getenv("EQ_PATH")
	if eqgPath == "" {
		b.Skip("EQ_PATH not set")
	}

	for i := 0; i < b.N; i++ {
		pfs, err := NewFile(fmt.Sprintf("%s/xhf.eqg", eqgPath))
		if err != nil {
			b.Fatalf("Failed newfile: %s", err.Error())
		}
		pfs.Close()
	}
}
