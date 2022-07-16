package lay

import (
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	//	var err error

	//	fmt.Println(hex.Dump(buf.Bytes()))
}
