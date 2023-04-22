package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs/s3d"
)

func fragmentTests(t *testing.T, isSingleRun bool, names []string, fragCode int32, fragOffset int, parser func(r io.ReadSeeker, fragOffset int) error) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {

			if fragOffset == -1 {
				count, err := parseFragments(t, isSingleRun, fmt.Sprintf("%s/_test_data/%s", eqPath, name), fragCode, parser)
				if err != nil {
					t.Fatalf("%s parseFragments: %v", name, err)
					return
				}
				log.Debugf("%s total parsed: %d", name, count)
				return
			}
			err := parseFragment(t, fmt.Sprintf("%s/_test_data/%s", eqPath, name), fragCode, fragOffset, parser)
			if err != nil {
				t.Fatalf("%s parseFragment: %v", name, err)
				return
			}
		})
	}
}

func parseFragment(t *testing.T, path string, fragCode int32, fragOffset int, parser func(r io.ReadSeeker, fragOffset int) error) error {
	r, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("open: %w", err)
		} else {
			extractFragments(t, filepath.Base(path))
			r, err = os.Open(path)
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
		}
	}
	defer r.Close()

	return parser(r, fragOffset)
}

func parseFragments(t *testing.T, isSingleRun bool, path string, fragCode int32, parser func(r io.ReadSeeker, fragOffset int) error) (int, error) {

	total := 0
	files, err := os.ReadDir(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return total, fmt.Errorf("read dir: %w", err)
		} else {
			extractFragments(t, filepath.Base(path))
			files, err = os.ReadDir(path)
			if err != nil {
				return total, fmt.Errorf("read dir: %w", err)
			}
		}
	}
	for _, fe := range files {
		if fe.IsDir() {
			continue
		}

		if !strings.HasSuffix(fe.Name(), fmt.Sprintf("0x%02x.hex", fragCode)) {
			continue
		}

		fragOffset, err := strconv.Atoi(fe.Name()[0:4])
		if err != nil {
			return total, fmt.Errorf("strconv %s: %w", fe.Name(), err)
		}

		log.Debugf(fe.Name())
		err = parseFragment(t, path+"/"+fe.Name(), fragCode, fragOffset, parser)
		if err != nil {
			return total, fmt.Errorf("parse fragment %s: %w", fe.Name(), err)
		}
		total++
		if isSingleRun {
			log.Debugf("stopping early, single run")
			return total, nil
		}
	}

	return total, nil
}

func compareReadAndWrite(t *testing.T, path string, fragCode int32, fragOffset int, e *WLD) error {
	total := 0
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}
	for _, fe := range files {
		if fe.IsDir() {
			continue
		}

		if !strings.HasSuffix(fe.Name(), fmt.Sprintf("0x%02x.hex", fragCode)) {
			continue
		}

		fragOffset, err := strconv.Atoi(fe.Name()[0:4])
		if err != nil {
			return fmt.Errorf("parse frag offset: %w", err)
		}

		log.Debugf(fe.Name())
		err = parseFragment(t, path+"/"+fe.Name(), fragCode, fragOffset, e.textureListRead)
		if err != nil {
			return fmt.Errorf("parse fragment: %w", err)
		}
		total++

		err = e.fragments[fragOffset].build(e)
		if err != nil {
			return fmt.Errorf("build: %w", err)
		}
		buf := bytes.NewBuffer(nil)

		err = e.textureListWrite(buf, fragOffset)
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}

		r, err := os.Open(path + "/" + fe.Name())
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer r.Close()

		// read r to a byte slice
		// compare byte slice to buf
		// if not equal, fail
		buf2 := bytes.NewBuffer(nil)
		_, err = buf2.ReadFrom(r)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		for i := 0; i < buf.Len(); i++ {
			if buf.Bytes()[i] != buf2.Bytes()[i] {
				return fmt.Errorf("byte %d not equal", i)
			}
		}

	}
	log.Debugf("Total parsed: %d", total)
	return nil
}

func extractFragments(t *testing.T, name string) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	log.Debugf("extracting fragments %s", name)

	out := fmt.Sprintf("%s/_test_data/%s", eqPath, name)

	pfs, err := s3d.NewFile(fmt.Sprintf("%s/%s", eqPath, name))
	if err != nil {
		t.Fatalf("s3d.NewFile: %v", err)
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
