package raw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/tag"
)

func TestBmpRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()
	type args struct {
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .bmp|1|sidl_ba_1_tln.bmp|tln.eqg
		{name: "tln.eqg"}, // PASS
		// .bmp|2|stnd_ba_1_exo.bmp|exo.eqg eye_chr.s3d pfs import: s3d load: decode: dirName for crc 655939147 not found
		// .bmp|2|walk_ba_1_vaf.bmp|vaf.eqg valdeholm.eqg pfs import: eqg load: decode: read nameData unexpected EOF
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open pfs %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".bmp" {
					continue
				}
				bmp := &Bmp{}
				err = bmp.Read(bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}
			}
		})
	}
}

func TestBmpWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "gequip2.s3d"},
	}
	if !common.IsTestExtensive() {
		tests = []struct {
			name    string
			wantErr bool
		}{
			{name: "gequip2.s3d"}, // pass
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			matchCount := 0
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".bmp" {
					continue
				}
				bmp := &Bmp{}
				err = bmp.Read(bytes.NewReader(file.Data()))

				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = bmp.Write(buf)
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				srcData := file.Data()
				dstData := buf.Bytes()

				err = common.ByteCompareTest(srcData, dstData)
				if err != nil {
					t.Fatalf("%s byteCompare: %s", tt.name, err)
				}

				matchCount++
			}
			fmt.Printf("matchCount for %s: %d\n", tt.name, matchCount)
		})
	}
}
