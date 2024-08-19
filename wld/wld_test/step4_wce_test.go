package wld_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wld"
)

func TestWceReadWrite(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	for _, tt := range tests {
		t.Run(tt.baseName, func(t *testing.T) {

			if os.Getenv("TEST_ARG") != "" {
				tt.baseName = os.Getenv("TEST_ARG")
			}

			totalStart := time.Now()

			start := time.Now()

			baseName := tt.baseName
			// copy original
			copyData, err := os.ReadFile(fmt.Sprintf("%s/%s.s3d", eqPath, baseName))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", baseName, err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.src.s3d", dirTest, baseName), copyData, 0644)
			if err != nil {
				t.Fatalf("failed to write s3d %s: %s", baseName, err.Error())
			}

			archive, err := pfs.NewFile(fmt.Sprintf("%s/%s.s3d", eqPath, baseName))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", baseName, err.Error())
			}
			defer archive.Close()

			if tt.wldName == "" {
				tt.wldName = fmt.Sprintf("%s.wld", tt.baseName)
			}
			// get wld
			data, err := archive.File(tt.wldName)
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", baseName, err.Error())
			}
			err = os.WriteFile(fmt.Sprintf("%s/%s.src.wld", dirTest, baseName), data, 0644)
			if err != nil {
				t.Fatalf("failed to write wld %s: %s", baseName, err.Error())
			}
			fmt.Println("Wrote", fmt.Sprintf("%s/%s.src.wld in %0.2f seconds", dirTest, baseName, time.Since(start).Seconds()))
			start = time.Now()

			rawWldSrc := &raw.Wld{}
			err = rawWldSrc.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

			wldSrc := &wld.Wld{}
			err = wldSrc.ReadRaw(rawWldSrc)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", baseName, err.Error())
			}

			fmt.Println("Read", fmt.Sprintf("%s/%s.src.wld", dirTest, baseName))

			wldSrc.FileName = baseName + ".wld"

			err = wldSrc.WriteAscii(dirTest + "/" + baseName)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			fmt.Println("Wrote", fmt.Sprintf("%s/%s/_root.wce in %0.2f seconds", dirTest, baseName, time.Since(start).Seconds()))

			start = time.Now()
			wldDst := &wld.Wld{
				FileName: baseName + ".wld",
			}
			err = wldDst.ReadAscii(fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

			fmt.Println("Read", fmt.Sprintf("%s/%s/_root.wce in %0.2f seconds", dirTest, baseName, time.Since(start).Seconds()))
			start = time.Now()

			// write back out

			dstBuf := bytes.NewBuffer(nil)

			err = wldDst.WriteRaw(dstBuf)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.dst.wld", dirTest, baseName), dstBuf.Bytes(), 0644)
			if err != nil {
				t.Fatalf("failed to write wld %s: %s", baseName, err.Error())
			}

			fmt.Println("Wrote", fmt.Sprintf("%s/%s.dst.wld in %0.2f seconds", dirTest, baseName, time.Since(start).Seconds()))

			rawWldDst := &raw.Wld{}

			/* diff := deep.Equal(wldSrc, wldDst)
			if diff != nil {
				t.Fatalf("wld diff: %s", diff)
			} */

			err = rawWldDst.Read(bytes.NewReader(dstBuf.Bytes()))
			if err != nil {
				t.Fatalf("failed to read wld3 %s: %s", baseName, err.Error())
			}

			/* diff := deep.Equal(wldSrc, wldDst)
			if diff != nil {
				t.Fatalf("wld diff: %s", diff)
			} */

			diff := deep.Equal(rawWldSrc, rawWldDst)
			if diff != nil {
				t.Fatalf("rawWld diff: %s", diff)
			}

			for i := 0; i < len(rawWldSrc.Fragments); i++ {
				srcFrag := rawWldSrc.Fragments[i]
				dstFrag := rawWldDst.Fragments[i]
				if srcFrag.FragCode() != dstFrag.FragCode() {
					t.Fatalf("fragment %d fragcode mismatch: src: %s, dst: %s", i, raw.FragName(srcFrag.FragCode()), raw.FragName(dstFrag.FragCode()))
				}
			}

			for i := 0; i < len(rawWldSrc.Fragments); i++ {
				srcFrag := rawWldSrc.Fragments[i]
				dstFrag := rawWldDst.Fragments[i]

				srcFragBuf := bytes.NewBuffer(nil)
				err = srcFrag.Write(srcFragBuf, rawWldDst.IsNewWorld)
				if err != nil {
					t.Fatalf("failed to write src frag %d: %s", i, err.Error())
				}

				dstFragBuf := bytes.NewBuffer(nil)
				err = dstFrag.Write(dstFragBuf, rawWldDst.IsNewWorld)
				if err != nil {
					t.Fatalf("failed to write dst frag %d: %s", i, err.Error())
				}

				err = common.ByteCompareTest(srcFragBuf.Bytes(), dstFragBuf.Bytes())
				if err != nil {
					t.Fatalf("%s byteCompare frag %d %s: %s", raw.FragName(srcFrag.FragCode()), i, tt.baseName, err)
				}
			}
			fmt.Printf("Processed %d fragments for %s in %0.2f seconds\n", len(rawWldSrc.Fragments), tt.baseName, time.Since(totalStart).Seconds())
		})
	}
}
