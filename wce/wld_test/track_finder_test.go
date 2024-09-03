package wce_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

func TestFragTrack(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	dirs, err := os.ReadDir(eqPath)
	if err != nil {
		t.Fatalf("Failed to read eq path: %s", err.Error())
	}

	w, err := os.Create("tracks.txt")
	if err != nil {
		t.Fatalf("Failed to create tracks.txt: %s", err.Error())
	}
	defer w.Close()

	fmt.Printf("Checking %d files\n", len(dirs))
	parseCount := 0
	for _, dir := range dirs {
		parseCount++
		if parseCount%100 == 0 {
			fmt.Printf("Checking %d/%d\n", parseCount, len(dirs))
		}

		ext := filepath.Ext(dir.Name())
		if ext != ".s3d" {
			continue
		}
		if dir.IsDir() {
			continue
		}

		s3dName := dir.Name()
		s3dPath := fmt.Sprintf("%s/%s", eqPath, s3dName)

		pfs, err := pfs.NewFile(s3dPath)
		if err != nil {
			t.Fatalf("Failed to open s3d %s: %s", s3dName, err.Error())
		}

		for _, file := range pfs.Files() {
			if filepath.Ext(file.Name()) != ".wld" {
				continue
			}
			wldName := file.Name()
			if wldName == "objects.wld" {
				continue
			}
			if wldName == "lights.wld" {
				continue
			}

			rawWld := &raw.Wld{}
			err = rawWld.Read(bytes.NewReader(file.Data()))
			if err != nil {
				t.Fatalf("Failed to read wld: %s", err.Error())
			}

			for i := 0; i < len(rawWld.Fragments); i++ {
				fragRaw := rawWld.Fragments[i]

				tagName := ""
				frag, ok := fragRaw.(*rawfrag.WldFragTrack)
				if ok {
					tagName = rawWld.Name(int32(frag.NameRef))
				} else {
					frag2, ok := fragRaw.(*rawfrag.WldFragTrackDef)
					if !ok {
						continue
					}
					tagName = rawWld.Name(int32(frag2.NameRef))
				}

				fmt.Fprintf(w, "Track found in %s/%s/%s fragID %d\n", s3dName, wldName, tagName, i)

			}
		}
	}

}
