package wce

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-test/deep"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

func TestBWldRead(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	tests := []struct {
		name string
	}{
		{"crushbone"},
		//{"gequip6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s.s3d", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.name, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(fmt.Sprintf("%s.wld", tt.name))
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.name, err.Error())
			}

			wld := &raw.Wld{}
			err = wld.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			vwld := New(tt.name)
			err = vwld.ReadWldRaw(wld)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", tt.name, err.Error())
			}

		})
	}
}

func TestBWldReadWriteRead(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	tests := []struct {
		baseName string
		wldName  string
	}{
		//{baseName: "load2"},
		//{baseName: "gequip4"},
		//{baseName: "global_chr"}, // TODO:  anarelion asked mesh of EYE_DMSPRITEDEF check if the eye is just massive 22 units in size, where the other units in that file are just 1-2 units in size
		{baseName: "qeynos"},
		//{baseName: "overc", wldName: "lights.wld"},
		//{baseName: "gequip6"},
		//{baseName: "load2", wldName: "lights.wld"},
		//{baseName: "load2", wldName: "objects.wld"},
		//{baseName: "neriakc"},
		//{baseName: "westwastes"},
		//{baseName: "globalfroglok_chr"},
		//{baseName: "dulak_obj"}, // TODO: dmtrackdef2
		//{baseName: "griegsend_chr"}, // long to load but good stress test
	}
	for _, tt := range tests {
		t.Run(tt.baseName, func(t *testing.T) {

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
			} else {
				baseName = tt.wldName[:len(tt.wldName)-4]
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

			wld := &raw.Wld{}
			err = wld.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}
			fmt.Println("read", fmt.Sprintf("%s/%s.src.wld", dirTest, baseName))

			if tt.wldName == "objects.wld" {
				data, err = archive.File(fmt.Sprintf("%s.wld", tt.baseName))
				if err != nil {
					t.Fatalf("failed to open wld %s: %s", tt.baseName, err.Error())
				}
				tmpWld := &raw.Wld{}
				err = tmpWld.Read(bytes.NewReader(data))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.baseName, err.Error())
				}
			}

			vwld := New(baseName)
			err = vwld.ReadWldRaw(wld)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", baseName, err.Error())
			}

			buf := bytes.NewBuffer(nil)
			err = vwld.WriteWldRaw(buf)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.dst.wld", dirTest, baseName), buf.Bytes(), 0644)
			if err != nil {
				t.Fatalf("failed to write wld %s: %s", baseName, err.Error())
			}

			fmt.Println("wrote", fmt.Sprintf("%s/%s.dst.wld", dirTest, baseName))

			// read back in
			wld2 := &raw.Wld{}
			r := bytes.NewReader(buf.Bytes())
			err = wld2.Read(r)
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

			fmt.Println("read", fmt.Sprintf("%s/%s.dst.wld", dirTest, baseName))

		})
	}
}

func TestWCEWldReadWriteRead(t *testing.T) {

	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	tests := []struct {
		baseName string
		wldName  string
	}{
		//{baseName: "gequip6"},
		//{baseName: "crushbone"},
		//{baseName: "gequip2"}, // hierarchical sprite
		//{baseName: "overthere_chr"},
		//{baseName: "hollows"},
		//{baseName: "illithid_chr"},
		{baseName: "beetle_chr"},
		//{baseName: "globalogm_chr"},
		//{baseName: "qeynos_chr"},
		//{baseName: "global_chr"}, // TODO:  anarelion asked mesh of EYE_DMSPRITEDEF check if the eye is just massive 22 units in size, where the other units in that file are just 1-2 units in size
		//{baseName: "load2"},
		//{baseName: "qeynos"},
		//	{baseName: "qeynos", wldName: "lights.wld"},
		//{baseName: "load2", wldName: "lights.wld"},
		//{baseName: "load2", wldName: "load2.wld"},
		//{baseName: "overthere_chr", wldName: "overthere_chr.wld"},

		//{baseName: "neriakc"},
		//{baseName: "westwastes"},
		//{baseName: "globalfroglok_chr"},
		//{baseName: "dulak_obj"}, // TODO: dmtrackdef2
		//{baseName: "griegsend_chr"}, // long to load but good stress test
	}
	for _, tt := range tests {
		t.Run(tt.baseName, func(t *testing.T) {

			if os.Getenv("TEST_ARG") != "" {
				tt.baseName = os.Getenv("TEST_ARG")
			}

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
			} else {
				baseName = tt.wldName[:len(tt.wldName)-4]
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
			fmt.Println("wrote", fmt.Sprintf("%s/%s.src.wld", dirTest, baseName))
			wld := &raw.Wld{}
			err = wld.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}
			/*
				if tt.wldName == "objects.wld" {
					data, err = archive.File(fmt.Sprintf("%s.wld", tt.baseName))
					if err != nil {
						t.Fatalf("failed to open wld %s: %s", tt.baseName, err.Error())
					}
					tmpWld := &raw.Wld{}
					err = tmpWld.Read(bytes.NewReader(data))
					if err != nil {
						t.Fatalf("failed to read %s: %s", tt.baseName, err.Error())
					}
				} */

			vwld := New(baseName)
			err = vwld.ReadWldRaw(wld)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", baseName, err.Error())
			}

			fmt.Println("read", fmt.Sprintf("%s/%s.src.wld", dirTest, baseName))

			vwld.FileName = baseName + ".wld"

			err = vwld.WriteAscii(dirTest + "/" + baseName)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			fmt.Println("wrote", fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))

			// read back in

			vwld2 := New(baseName)
			err = vwld2.ReadAscii(fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

			fmt.Println("read", fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))

			// write back out

			buf2 := bytes.NewBuffer(nil)

			err = vwld2.WriteWldRaw(buf2)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.dst.wld", dirTest, baseName), buf2.Bytes(), 0644)
			if err != nil {
				t.Fatalf("failed to write wld %s: %s", baseName, err.Error())
			}

			fmt.Println("wrote", fmt.Sprintf("%s/%s.dst.wld", dirTest, baseName))

		})
	}
}

func TestWCEWldFragMatch(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	tests := []struct {
		baseName string
		wldName  string
	}{
		//{baseName: "load2"},
		//{baseName: "load2", wldName: "objects.wld"},
		//{baseName: "beetle_chr"},
		{baseName: "qeynos_chr"},
		//{baseName: "overthere_chr"},
		//{baseName: "globalogm_chr"},

	}
	for _, tt := range tests {
		t.Run(tt.baseName, func(t *testing.T) {

			if os.Getenv("TEST_ARG") != "" {
				tt.baseName = os.Getenv("TEST_ARG")
			}

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
			fmt.Println("wrote", fmt.Sprintf("%s/%s.src.wld", dirTest, baseName))
			rawWldSrc := &raw.Wld{}
			err = rawWldSrc.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

			wldSrc := New(baseName)
			err = wldSrc.ReadWldRaw(rawWldSrc)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", baseName, err.Error())
			}

			fmt.Println("read", fmt.Sprintf("%s/%s.src.wld", dirTest, baseName))

			wldSrc.FileName = baseName + ".wld"

			err = wldSrc.WriteAscii(dirTest + "/" + baseName)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			fmt.Println("wrote", fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))

			// read back in

			wldDst := New(baseName + ".wld")
			err = wldDst.ReadAscii(fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

			fmt.Println("read", fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))

			// write back out

			dstBuf := bytes.NewBuffer(nil)

			err = wldDst.WriteWldRaw(dstBuf)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.dst.wld", dirTest, baseName), dstBuf.Bytes(), 0644)
			if err != nil {
				t.Fatalf("failed to write wld %s: %s", baseName, err.Error())
			}

			fmt.Println("wrote", fmt.Sprintf("%s/%s.dst.wld", dirTest, baseName))

			rawWldDst := &raw.Wld{}

			/* diff := deep.Equal(wldSrc, wldDst)
			if diff != nil {
				t.Fatalf("wld diff: %s", diff)
			} */

			err = rawWldDst.Read(bytes.NewReader(dstBuf.Bytes()))
			if err != nil {
				t.Fatalf("failed to read wld3 %s: %s", baseName, err.Error())
			}

			diff := deep.Equal(rawWldSrc, rawWldDst)
			if diff != nil {
				t.Fatalf("wld diff: %s", diff)
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
				err = srcFrag.Write(srcFragBuf, rawWldSrc.IsNewWorld)
				if err != nil {
					t.Fatalf("failed to write src frag %d: %s", i, err.Error())
				}

				dstFragBuf := bytes.NewBuffer(nil)
				err = dstFrag.Write(dstFragBuf, rawWldSrc.IsNewWorld)
				if err != nil {
					t.Fatalf("failed to write dst frag %d: %s", i, err.Error())
				}

				err = common.ByteCompareTest(srcFragBuf.Bytes(), dstFragBuf.Bytes())
				if err != nil {
					t.Fatalf("%s byteCompare frag %d %s: %s", raw.FragName(srcFrag.FragCode()), i, tt.baseName, err)
				}
			}

		})
	}
}

func TestAsciiRead(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		asciiName string
		wantErr   bool
	}{}
	if !common.IsTestExtensive() {
		tests = []struct {
			asciiName string
			wantErr   bool
		}{
			//{"all/all.spk", false},
			//{"fis/fis.spk", false},
			{"pre/pre.spk", false},
		}
	}
	for _, tt := range tests {
		t.Run(tt.asciiName, func(t *testing.T) {

			wld := New(tt.asciiName)
			err := wld.ReadAscii(fmt.Sprintf("testdata/%s", tt.asciiName))
			if err != nil {
				t.Fatalf("Failed readascii: %s", err.Error())
			}
		})
	}
}

func TestAsciiReadWriteRead(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		asciiName string
		wantErr   bool
	}{}
	if !common.IsTestExtensive() {
		tests = []struct {
			asciiName string
			wantErr   bool
		}{
			//{"all/all.spk", false},
			//{"fis/fis.spk", false},
			//{"pre/pre.spk", false},
			{"akheva.wld", false},
		}
	}
	for _, tt := range tests {
		t.Run(tt.asciiName, func(t *testing.T) {

			wld := New(filepath.Base(tt.asciiName))

			err := wld.ReadAscii(fmt.Sprintf("testdata/%s", tt.asciiName))
			if err != nil {
				t.Fatalf("Failed readascii: %s", err.Error())
			}

			err = wld.WriteAscii("testdata/temp/")
			if err != nil {
				t.Fatalf("Failed writeascii: %s", err.Error())
			}

			ext := filepath.Ext(wld.FileName)
			if ext == ".wld" {
				wld.FileName = wld.FileName[:len(wld.FileName)-len(ext)] + ".spk"
			}

			wld2 := New(filepath.Base(wld.FileName))
			err = wld2.ReadAscii("testdata/temp/" + wld.FileName)
			if err != nil {
				t.Fatalf("Failed re-readascii: %s", err.Error())
			}

		})
	}
}
