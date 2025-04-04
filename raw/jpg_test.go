package raw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestJpgRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()
	type args struct {
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "atiiki.eqg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open pfs %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".jpg" {
					continue
				}
				jpg := &Jpg{}
				err = jpg.Read(bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}
			}
		})
	}
}

func TestJpgWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()

	tests := []struct {
		name    string
		wantErr bool
	}{
		// .jpg|1|sidl_ba_1_tln.jpg|tln.eqg
		{name: "tln.eqg"}, // PASS
		// .jpg|2|stnd_ba_1_exo.jpg|exo.eqg eye_chr.s3d pfs import: s3d load: read: dirName for crc 655939147 not found
		// .jpg|2|walk_ba_1_vaf.jpg|vaf.eqg valdeholm.eqg pfs import: eqg load: read: read nameData unexpected EOF
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".jpg" {
					continue
				}
				jpg := &Jpg{}
				err = jpg.Read(bytes.NewReader(file.Data()))

				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = jpg.Write(buf)
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				srcData := file.Data()
				dstData := buf.Bytes()

				err = helper.ByteCompareTest(srcData, dstData)
				if err != nil {
					t.Fatalf("%s byteCompare: %s", tt.name, err)
				}
			}
		})
	}
}
