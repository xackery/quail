package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/xackery/quail/pfs"
)

func TestExtractWldFragment(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	name := "gequip.s3d"
	out := fmt.Sprintf("%s/_test_data/%s", eqPath, name)

	pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, name))
	if err != nil {
		t.Fatalf("pfs.NewFile: %v", err)
		return
	}
	wldName := strings.TrimSuffix(name, ".s3d") + ".wld"

	data, err := pfs.File(wldName)
	if err != nil {
		t.Fatalf("pfs.File: %v", err)
		return
	}

	e, err := New(wldName, pfs)
	if err != nil {
		t.Fatalf("New: %v", err)
		return
	}

	r := bytes.NewReader(data)
	fragmentCount, err := e.readHeader(r)
	if err != nil {
		t.Fatalf("readHeader: %v", err)
		return
	}

	err = os.MkdirAll(out, 0755)
	if err != nil {
		t.Fatalf("mkdir: %v", err)
		return
	}

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(fragmentCount); fragOffset++ {
		var fragSize uint32
		var fragCode int32

		err = binary.Read(r, binary.LittleEndian, &fragSize)
		if err != nil {
			t.Fatalf("read fragment size %d/%d: %v", fragOffset, fragmentCount, err)
		}
		totalFragSize += fragSize
		//dump.Hex(fragSize, "%d(%s)fragSize=%d", i, name, fragSize)
		err = binary.Read(r, binary.LittleEndian, &fragCode)
		if err != nil {
			t.Fatalf("read fragment index %d/%d: %v", fragOffset, fragmentCount, err)
		}
		//dump.Hex(fragSize, "%dfragCode=%d", i, fragCode)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			t.Fatalf("frag position seek %d/%d: %v", fragOffset, fragmentCount, err)
		}

		buf := make([]byte, fragSize)
		_, err = r.Read(buf)
		if err != nil {
			t.Fatalf("read: %v", err)
		}

		err = os.WriteFile(fmt.Sprintf("%s/%04d_0x%02x.hex", out, fragOffset, fragCode), buf, 0644)
		if err != nil {
			t.Fatalf("write: %v", err)
		}

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			t.Fatalf("seek end of frag %d/%d: %v", fragOffset, fragmentCount, err)
		}
	}

}
