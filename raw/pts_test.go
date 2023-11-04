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

func TestPtsRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)
	type args struct {
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .pts|1|aam.pts|aam.eqg
		{name: "aam.eqg"},
		// .pts|1|ae3.pts|ae3.eqg
		{name: "ae3.eqg"},
		// .pts|1|ahf.pts|ahf.eqg
		{name: "ahf.eqg"},
		// .pts|1|ahm.pts|ahm.eqg
		{name: "ahm.eqg"},
		// .pts|1|aie.pts|aie.eqg
		{name: "aie.eqg"},
		// .pts|1|ala.pts|ala.eqg
		{name: "ala.eqg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".pts" {
					continue
				}
				pts := &Pts{}
				err = pts.Read(bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}

func TestEncode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)

	tests := []struct {
		name    string
		wantErr bool
	}{
		// .pts|1|aam.pts|aam.eqg
		{name: "aam.eqg"},
		// .pts|1|ae3.pts|ae3.eqg
		{name: "ae3.eqg"},
		// .pts|1|ahf.pts|ahf.eqg
		{name: "ahf.eqg"},
		// .pts|1|ahm.pts|ahm.eqg
		{name: "ahm.eqg"},
		// .pts|1|aie.pts|aie.eqg
		{name: "aie.eqg"},
		// .pts|1|ala.pts|ala.eqg
		{name: "ala.eqg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".pts" {
					continue
				}
				pts := &Pts{}
				err = pts.Read(bytes.NewReader(file.Data()))

				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = pts.Write(buf)
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				srcData := file.Data()
				dstData := buf.Bytes()

				err = common.ByteCompareTest(srcData, dstData)
				if err != nil {
					t.Fatalf("%s failed byteCompare: %s", tt.name, err)
				}
			}
		})
	}
}
