package vwld

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

func TestVWldRead(t *testing.T) {
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

			vwld := &VWld{}
			err = vwld.Read(wld)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", tt.name, err.Error())
			}

		})
	}
}

func TestVWldWrite(t *testing.T) {
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
		//{baseName: "gequip4"},
		//{baseName: "global_chr"}, // TODO:  anarelion asked mesh of EYE_DMSPRITEDEF check if the eye is just massive 22 units in size, where the other units in that file are just 1-2 units in size
		//{baseName: "load2"},
		//{baseName: "load2", wldName: "lights.wld"},
		{baseName: "load2", wldName: "objects.wld"},
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
			w, err := os.Create(fmt.Sprintf("%s/%s.src.wld", dirTest, baseName))
			if err != nil {
				t.Fatalf("failed to create %s: %s", baseName, err.Error())
			}
			defer w.Close()

			_, err = w.Write(data)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			wld := &raw.Wld{}
			err = wld.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read %s: %s", baseName, err.Error())
			}

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

			vwld := &VWld{}
			err = vwld.Read(wld)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", baseName, err.Error())
			}

			err = vwld.Write(w)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			buf := bytes.NewBuffer(nil)
			err = wld.Write(buf)
			if err != nil {
				t.Fatalf("failed to write %s: %s", baseName, err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.dst.wld", dirTest, baseName), buf.Bytes(), 0644)
			if err != nil {
				t.Fatalf("failed to write wld %s: %s", baseName, err.Error())
			}

		})
	}
}
