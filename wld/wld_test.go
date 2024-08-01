package wld

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

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
		//{"crushbone"},
		{"gequip6"},
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

			vwld := &Wld{}
			err = vwld.ReadRaw(wld)
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

			vwld := &Wld{}
			err = vwld.ReadRaw(wld)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", baseName, err.Error())
			}

			buf := bytes.NewBuffer(nil)
			err = vwld.WriteRaw(buf)
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
		//{baseName: "gequip2"},
		//{baseName: "qeynos_chr"},
		//{baseName: "global_chr"}, // TODO:  anarelion asked mesh of EYE_DMSPRITEDEF check if the eye is just massive 22 units in size, where the other units in that file are just 1-2 units in size
		//{baseName: "load2"},
		{baseName: "qeynos"},
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

			vwld := &Wld{}
			err = vwld.ReadRaw(wld)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", baseName, err.Error())
			}

			fmt.Println("read", fmt.Sprintf("%s/%s.src.wld", dirTest, baseName))

			vwld.FileName = baseName + ".wld"

			err = vwld.WriteAscii(dirTest+"/"+baseName, true)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			fmt.Println("wrote", fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))

			// read back in

			vwld2 := &Wld{}
			err = vwld2.ReadAscii(fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

			fmt.Println("read", fmt.Sprintf("%s/%s/_root.wce", dirTest, baseName))

			// write back out

			buf2 := bytes.NewBuffer(nil)

			err = vwld2.WriteRaw(buf2)
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

			wld := &Wld{}
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

			wld := &Wld{
				FileName: filepath.Base(tt.asciiName),
			}
			err := wld.ReadAscii(fmt.Sprintf("testdata/%s", tt.asciiName))
			if err != nil {
				t.Fatalf("Failed readascii: %s", err.Error())
			}

			err = wld.WriteAscii("testdata/temp/", true)
			if err != nil {
				t.Fatalf("Failed writeascii: %s", err.Error())
			}

			ext := filepath.Ext(wld.FileName)
			if ext == ".wld" {
				wld.FileName = wld.FileName[:len(wld.FileName)-len(ext)] + ".spk"
			}

			wld2 := &Wld{}
			err = wld2.ReadAscii("testdata/temp/" + wld.FileName)
			if err != nil {
				t.Fatalf("Failed re-readascii: %s", err.Error())
			}

		})
	}
}
