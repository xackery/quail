package wce_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
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

			ext := ".eqg"
			archive, err := pfs.NewFile(fmt.Sprintf("%s/%s.eqg", eqPath, baseName))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", baseName, err.Error())
			}
			defer archive.Close()

			wldSrc := wce.New(baseName + ".eqg")
			err = wldSrc.ReadEqgRaw(archive)
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

			fmt.Printf("Processed %s in %0.2f seconds\n", tt.baseName, time.Since(totalStart).Seconds())

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

			outArchive, err := pfs.New(baseName + ".dst.eqg")
			if err != nil {
				t.Fatalf("failed to create pfs: %s", err.Error())
			}

			err = wldDst.WriteEqgRaw(outArchive)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			w, err := os.Create(fmt.Sprintf("%s/%s.dst.eqg", dirTest, baseName))
			if err != nil {
				t.Fatalf("failed to create %s: %s", baseName, err.Error())
			}
			defer w.Close()

			err = outArchive.Write(w)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			fmt.Println("Wrote", fmt.Sprintf("%s/%s.dst%s in %0.2f seconds", dirTest, baseName, ext, time.Since(start).Seconds()))

		})
	}
}
