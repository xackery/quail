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
	dirTest := common.DirTest()
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

func TestPtsWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	tests := []struct {
		name    string
		wantErr bool
		isDump  bool
	}{
		// .pts|1|aam.pts|aam.eqg
		//{name: "aam.eqg", isDump: true}, // PASS
		// .pts|1|ae3.pts|ae3.eqg
		// {name: "ae3.eqg"}, // PASS
		// .pts|1|ahf.pts|ahf.eqg
		// {name: "ahf.eqg"}, // PASS
		// .pts|1|ahm.pts|ahm.eqg
		// {name: "ahm.eqg"}, // PASS
		// .pts|1|aie.pts|aie.eqg
		// {name: "aie.eqg"}, // PASS
		// .pts|1|ala.pts|ala.eqg
		//{name: "ala.eqg"}, // PASS
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
				if tt.isDump {
					os.WriteFile(fmt.Sprintf("%s/%s.src.pts", dirTest, file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.src.pts.tags", dirTest, file.Name()))
					fmt.Printf("dumped to %s\n", fmt.Sprintf("%s/%s.src.pts", dirTest, file.Name()))
				}
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := common.NewByteSeekerTest()
				err = pts.Write(buf)
				if tt.isDump {
					os.WriteFile(fmt.Sprintf("%s/%s.dst.pts", dirTest, file.Name()), buf.Bytes(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.dst.pts.tags", dirTest, file.Name()))
					fmt.Printf("dumped to %s\n", fmt.Sprintf("%s/%s.dst.pts", dirTest, file.Name()))
				}
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				// TODO: add fluff data
				/* srcData := file.Data()
				dstData := buf.Bytes()

				err = common.ByteCompareTest(srcData, dstData)
				if err != nil {
					t.Fatalf("%s byteCompare: %s", tt.name, err)
				} */
			}
		})
	}
}
