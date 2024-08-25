package wld_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
	"github.com/xackery/quail/wld"
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

			fmt.Printf("Processed (src: %d, dst: %d) fragments for %s in %0.2f seconds\n", len(rawWldSrc.Fragments), len(rawWldDst.Fragments), tt.baseName, time.Since(totalStart).Seconds())

			srcFragByCodes := make(map[int]int)
			srcFragByTags := make(map[int][]*tagEntry)
			for i := 0; i < len(rawWldSrc.Fragments); i++ {
				srcFrag := rawWldSrc.Fragments[i]
				srcFragByCodes[srcFrag.FragCode()]++
				srcFragByTags[srcFrag.FragCode()] = append(srcFragByTags[srcFrag.FragCode()], &tagEntry{tag: rawWldSrc.TagByFrag(srcFrag), offset: i})
			}

			dstFragByCodes := make(map[int]int)
			dstFragByTags := make(map[int][]*tagEntry)

			for i := 0; i < len(rawWldDst.Fragments); i++ {
				dstFrag := rawWldDst.Fragments[i]
				dstFragByCodes[dstFrag.FragCode()]++
				dstFragByTags[dstFrag.FragCode()] = append(dstFragByTags[dstFrag.FragCode()], &tagEntry{tag: rawWldDst.TagByFrag(dstFrag), offset: i})
			}

			for i := range srcFragByTags {
				srcTags := srcFragByTags[i]
				dstTags := dstFragByTags[i]
				if i == rawfrag.FragCodeSimpleSprite {
					continue
				}
				//fmt.Printf("Comparing %d (%s) tags %d total\n", i, rawfrag.FragName(i), len(srcTags))
				// find a matching dstTag, and pop from both
				for _, srcTag := range srcTags {
					found := false
					for j, dstTag := range dstTags {
						if srcTag.tag == dstTag.tag {
							found = true
							dstTags = append(dstTags[:j], dstTags[j+1:]...)
							break
						}
					}
					if !found {
						t.Fatalf("fragment %d (%s) tag %s not found in dst", srcTag.offset, rawfrag.FragName(i), srcTag.tag)
					}
				}
				//if len(dstTags) > 0 {
				//fmt.Printf("Warning: fragment (%s) tags %v not found in src\n", rawfrag.FragName(i), dstTags)
				//}
			}

			for code, count := range srcFragByCodes {
				if count > dstFragByCodes[code] {
					t.Fatalf("fragment code %d (%s) count mismatch: src: %d, dst: %d", code, rawfrag.FragName(code), count, dstFragByCodes[code])
				}
			}
			for code, tags := range srcFragByTags {
				if len(tags) > len(dstFragByTags[code]) {
					t.Fatalf("fragment code %d (%s) tag count mismatch: src: %d, dst: %d", code, rawfrag.FragName(code), len(tags), len(dstFragByTags[code]))
				}
			}

			if len(rawWldSrc.Fragments) > len(rawWldDst.Fragments) {
				t.Fatalf("fragment count mismatch: src: %d, dst: %d", len(rawWldSrc.Fragments), len(rawWldDst.Fragments))
			}

			/*
				for i := 0; i < len(rawWldSrc.Fragments); i++ {
					srcFrag := rawWldSrc.Fragments[i]
					for j := 0; j < len(rawWldDst.Fragments); j++ {
						dstFrag := rawWldDst.Fragments[j]
						if srcFrag.FragCode() != dstFrag.FragCode() {
							continue
						}
						if rawWldSrc.TagByFrag(srcFrag) == "" {
							continue
						}
						if rawWldSrc.TagByFrag(srcFrag) != rawWldDst.TagByFrag(dstFrag) {
							continue
						}
						diff := deep.Equal(srcFrag, dstFrag)
						if diff != nil {
							t.Fatalf("fragment %d diff mismatch: src: %s (%s), dst: %s (%s), diff: %s", i, raw.FragName(srcFrag.FragCode()), rawWldSrc.TagByFrag(srcFrag), raw.FragName(dstFrag.FragCode()), rawWldDst.TagByFrag(dstFrag), diff)
						}
					}
				}*/
			/*
				diff := deep.Equal(rawWldSrc, rawWldDst)
				if diff != nil {
					t.Fatalf("rawWld diff: %s", diff)
				}
			*/
			/* for i := 0; i < len(rawWldSrc.Fragments); i++ {
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
			} */
		})
	}
}
