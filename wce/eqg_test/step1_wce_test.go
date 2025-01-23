package wce_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wce"
)

type tagEntry struct {
	tag    string
	offset int
}

func (e *tagEntry) String() string {
	return fmt.Sprintf("%s (%d)", e.tag, e.offset)
}

func TestWceReadWrite(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	knownExts := []string{".mds", ".mod", ".ter"}
	start := time.Now()
	for _, tt := range tests {
		t.Run(tt.baseName, func(t *testing.T) {

			if os.Getenv("TEST_ARG") != "" {
				tt.baseName = os.Getenv("TEST_ARG")
			}

			totalStart := time.Now()

			baseName := tt.baseName
			// copy original
			copyData, err := os.ReadFile(fmt.Sprintf("%s/%s.eqg", eqPath, baseName))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", baseName, err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.src.eqg", dirTest, baseName), copyData, 0644)
			if err != nil {
				t.Fatalf("Failed to write eqg %s: %s", baseName, err.Error())
			}

			archive, err := pfs.NewFile(fmt.Sprintf("%s/%s.eqg", eqPath, baseName))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", baseName, err.Error())
			}
			defer archive.Close()

			wldSrc := wce.New(baseName + ".eqg")

			var ext string
			if tt.fileName == "" {
				files := archive.Files()
				for _, file := range files {
					ext := filepath.Ext(file.Name())
					for _, knownExt := range knownExts {
						if ext == knownExt {
							parseEQGEntry(t, file.Name(), dirTest, baseName, archive, wldSrc)
							tt.fileName = file.Name()
						}
					}
				}

			} else {
				parseEQGEntry(t, tt.fileName, dirTest, baseName, archive, wldSrc)
			}
			dstBuf := bytes.NewBuffer(nil)

			fmt.Printf("Processed %s in %0.2f seconds\n", tt.baseName, time.Since(totalStart).Seconds())
			wldSrc.FileName = baseName + ext

			err = wldSrc.WriteAscii(dirTest + "/" + baseName)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			fmt.Println("Wrote", fmt.Sprintf("%s/%s/_root.wce in %0.2f seconds", dirTest, baseName, time.Since(start).Seconds()))

			start := time.Now()
			wldDst := wce.New(baseName + ext)
			err = wldDst.ReadAscii(fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

			fmt.Println("Read", fmt.Sprintf("%s/%s/_root.wce in %0.2f seconds", dirTest, baseName, time.Since(start).Seconds()))
			start = time.Now()

			// write back out

			err = wldDst.WriteEqgRaw(dstBuf)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.dst%s", dirTest, baseName, ext), dstBuf.Bytes(), 0644)
			if err != nil {
				t.Fatalf("failed to write %s %s: %s", baseName, ext, err.Error())
			}

			fmt.Println("Wrote", fmt.Sprintf("%s/%s.dst%s in %0.2f seconds", dirTest, baseName, ext, time.Since(start).Seconds()))

		})
	}
}

func parseEQGEntry(t *testing.T, fileName string, dirTest string, baseName string, archive *pfs.Pfs, wldSrc *wce.Wce) {

	start := time.Now()
	ext := filepath.Ext(fileName)

	if ext == "" {
		t.Fatalf("failed to find file to parse in %s", baseName)
	}

	data, err := archive.File(fileName)
	if err != nil {
		t.Fatalf("failed to open %s %s: %s", ext, baseName, err.Error())
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s.src%s", dirTest, baseName, ext), data, 0644)
	if err != nil {
		t.Fatalf("failed to write %s %s: %s", ext, baseName, err.Error())
	}
	fmt.Println("Wrote", fmt.Sprintf("%s/%s.src%s in %0.2f seconds", dirTest, baseName, ext, time.Since(start).Seconds()))

	var rawWldSrc raw.Reader
	switch ext {
	case ".mds":
		rawWldSrc = &raw.Mds{
			MetaFileName: strings.TrimSuffix(fileName, ".mds"),
		}
	case ".mod":
		rawWldSrc = &raw.Mod{
			MetaFileName: strings.TrimSuffix(fileName, ".mod"),
		}
	case ".ter":
		rawWldSrc = &raw.Ter{
			MetaFileName: strings.TrimSuffix(fileName, ".ter"),
		}

	default:
		t.Fatalf("unknown ext %s", ext)
	}

	err = rawWldSrc.Read(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to read %s: %s", baseName, err.Error())
	}

	err = wldSrc.ReadEqgRaw(rawWldSrc)
	if err != nil {
		t.Fatalf("failed to convert %s: %s", baseName, err.Error())
	}

	fmt.Println("Read", fmt.Sprintf("%s/%s.src%s", dirTest, baseName, ext))

}
